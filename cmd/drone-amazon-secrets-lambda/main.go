// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/drone/drone-amazon-secrets/plugin"
	"github.com/drone/drone-amazon-secrets/secret"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	addr = os.Getenv("SERVER_ADDR")
	key  = os.Getenv("SECRET_KEY")
)

func main() {
	if key == "" {
		log.Fatalln("fatal: missing secret key")
	}
	if addr == "" {
		addr = ":3000"
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalln(err)
	}

	handler := secret.Handler(
		key,
		plugin.New(
			secretsmanager.New(cfg),
		),
	)

	log.Printf("server listening on address %s", addr)

	http.Handle("/", handler)
	log.Fatal(gateway.ListenAndServe(addr, nil))
}
