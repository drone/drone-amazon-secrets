// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secret

// TODO(bradrydzewski) these types should eventually be moved
// to the github.com/drone/drone-go repository.

type (
	// Repo represents a Drone repository.
	Repo struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		Private   bool   `json:"private"`
	}

	// Build represents a Drone build.
	Build struct {
		Event  string `json:"event"`
		Ref    string `json:"ref"`
		Source string `json:"source"`
		Target string `json:"target"`
		Fork   string `json:"fork"`
	}
)
