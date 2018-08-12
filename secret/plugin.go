// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secret

type (
	// Request defines a secret request.
	Request struct {
		Name  string `json:"name"`
		Repo  Repo   `json:"repo,omitempty"`
		Build Build  `json:"build,omitempty"`
	}

	// Response defines a secret response.
	Response struct {
		Name string `json:"name,omitempty"`
		Data string `json:"data,omitempty"`
		Pull bool   `json:"pull,omitempty"`
		Fork bool   `json:"fork,omitempty"`
	}

	// Plugin responds to a secret request.
	Plugin interface {
		Handle(req *Request) (*Response, error)
	}
)
