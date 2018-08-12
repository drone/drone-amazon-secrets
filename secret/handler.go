// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secret

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/drone/drone-amazon-secrets/secret/aesgcm"

	"github.com/99designs/httpsignatures-go"
)

// Handler returns a http.Handler that accepts JSON-encoded
// HTTP requests for a secret, invokes the underlying secret
// plugin, and writes the JSON-encoded secret to the HTTP response.
//
// The handler verifies the authenticity of the HTTP request
// using the http-signature, and returns a 400 Bad Request if
// the signature is missing or invalid.
//
// The handler can optionally encrypt the response body using
// aesgcm if the HTTP request includes the Accept-Encoding header
// set to aesgcm.
func Handler(secret string, plugin Plugin) http.Handler {
	return &handler{
		secret: secret,
		plugin: plugin,
	}
}

type handler struct {
	secret string
	plugin Plugin
}

func (p *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature, err := httpsignatures.FromRequest(r)
	if err != nil {
		http.Error(w, "Invalid or Missing Signature", 400)
		return
	}
	if !signature.IsValid(p.secret, r) {
		http.Error(w, "Invalid Signature", 400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	req := &Request{}
	err = json.Unmarshal(body, req)
	if err != nil {
		http.Error(w, "Invalid Input", 400)
		return
	}

	secret, err := p.plugin.Handle(req)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	out, _ := json.Marshal(secret)

	// If the client can optionally accept an encrypted
	// response, we encrypt the payload body using secretbox.
	if r.Header.Get("Accept-Encoding") == "aesgcm" {
		key, err := aesgcm.Key(p.secret)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		out, err = aesgcm.Encrypt(out, key)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Encoding", "aesgcm")
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(out)
}
