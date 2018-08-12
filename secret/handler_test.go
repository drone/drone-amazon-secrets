// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secret

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drone/drone-amazon-secrets/secret/aesgcm"

	"github.com/99designs/httpsignatures-go"
)

func TestHandler(t *testing.T) {
	key := "xVKAGlWQiY3sOp8JVc0nbuNId3PNCgWh"
	want := &Response{
		Name: "docker_password",
		Data: "correct-horse-battery-staple",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(want)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", buf)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	err := httpsignatures.DefaultSha256Signer.AuthRequest("hmac-key", key, req)
	if err != nil {
		t.Error(err)
		return
	}

	plugin := &mockPlugin{
		res: want,
		err: nil,
	}

	handler := Handler(key, plugin)
	handler.ServeHTTP(res, req)

	if got, want := res.Code, 200; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}

	resp := &Response{}
	json.Unmarshal(res.Body.Bytes(), resp)
	if got, want := resp.Name, want.Name; got != want {
		t.Errorf("Want secret name %s, got %s", want, got)
	}
}

func TestHandler_Encrypted(t *testing.T) {
	key := "xVKAGlWQiY3sOp8JVc0nbuNId3PNCgWh"
	want := &Response{
		Name: "docker_password",
		Data: "correct-horse-battery-staple",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(want)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", buf)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Add("Accept-Encoding", "aesgcm")

	err := httpsignatures.DefaultSha256Signer.AuthRequest("hmac-key", key, req)
	if err != nil {
		t.Error(err)
		return
	}

	plugin := &mockPlugin{
		res: want,
		err: nil,
	}

	handler := Handler(key, plugin)
	handler.ServeHTTP(res, req)

	if got, want := res.Code, 200; got != want {
		t.Errorf("Want status code %d, got %d", want, got)
	}
	if got, want := res.Header().Get("Content-Encoding"), "aesgcm"; got != want {
		t.Errorf("Want Content-Encoding %s, got %s", want, got)
	}
	if got, want := res.Header().Get("Content-Type"), "application/octet-stream"; got != want {
		t.Errorf("Want Content-Type %s, got %s", want, got)
	}

	keyb, err := aesgcm.Key(key)
	if err != nil {
		t.Error(err)
		return
	}
	body, err := aesgcm.Decrypt(res.Body.Bytes(), keyb)
	if err != nil {
		t.Error(err)
		return
	}

	resp := &Response{}
	json.Unmarshal(body, resp)
	if got, want := resp.Name, want.Name; got != want {
		t.Errorf("Want secret name %s, got %s", want, got)
	}
}

func TestHandler_MissingSignature(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	handler := Handler("xVKAGlWQiY3sOp8JVc0nbuNId3PNCgWh", nil)
	handler.ServeHTTP(res, req)

	got, want := res.Body.String(), "Invalid or Missing Signature\n"
	if got != want {
		t.Errorf("Want response body %q, got %q", want, got)
	}
}

func TestHandler_InvalidSignature(t *testing.T) {
	sig := `keyId="hmac-key",algorithm="hmac-sha256",signature="QrS16+RlRsFjXn5IVW8tWz+3ZRAypjpNgzehEuvJksk=",headers="(request-target) accept accept-encoding content-type date digest"`
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Signature", sig)

	handler := Handler("xVKAGlWQiY3sOp8JVc0nbuNId3PNCgWh", nil)
	handler.ServeHTTP(res, req)

	got, want := res.Body.String(), "Invalid Signature\n"
	if got != want {
		t.Errorf("Want response body %q, got %q", want, got)
	}
}

type mockPlugin struct {
	res *Response
	err error
}

func (m *mockPlugin) Handle(req *Request) (*Response, error) {
	return m.res, m.err
}
