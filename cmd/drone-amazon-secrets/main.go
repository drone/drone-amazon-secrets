// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"os"

	"github.com/drone/drone-amazon-secrets/plugin"
	"github.com/drone/drone-amazon-secrets/server"
	"github.com/drone/drone-go/plugin/secret"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

var (
	debug = os.Getenv("DEBUG") == "true"
	addr  = os.Getenv("SERVER_ADDR")
	key   = os.Getenv("SECRET_KEY")
)

func main() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if key == "" {
		logrus.Fatalln("missing secret key")
	}
	if addr == "" {
		addr = ":3000"
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		logrus.Fatalln(err)
	}

	if cfg.Region == "" {
		metaClient := ec2metadata.New(cfg)
		if region, err := metaClient.Region(); err == nil {
			cfg.Region = region
			logrus.Infof("using region %s from ec2 metadata", cfg.Region)
		} else {
			logrus.Fatalf("failed to determine region: %s", err)
		}
	}

	handler := secret.Handler(
		key,
		plugin.New(secretsmanager.New(cfg)),
		logrus.StandardLogger(),
	)
	healthzHandler := server.HandleHealthz()

	logrus.Infof("server listening on address %s", addr)

	http.Handle("/", handler)
	http.Handle("/healthz",healthzHandler)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}
