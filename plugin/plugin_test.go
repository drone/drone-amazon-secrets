// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package plugin

import (
	"net/http"
	"testing"

	"github.com/drone/drone-amazon-secrets/secret"
	"github.com/google/go-cmp/cmp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/h2non/gock"
)

func TestPlugin(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "secretsmanager.GetSecretValue").
		Reply(200).
		File("testdata/secret.json")

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           "https://ec2.us-east-1.amazonaws.com",
		SigningRegion: config.Region,
	})

	manager := secretsmanager.New(config)
	req := &secret.Request{
		Name: "docker#username",
		Build: secret.Build{
			Event: "push",
		},
		Repo: secret.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin := New(manager)
	got, err := plugin.Handle(req)
	if err != nil {
		t.Error(err)
		return
	}

	want := &secret.Response{
		Name: "username",
		Data: "david",
		Pull: true,
		Fork: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
		return
	}

	if gock.IsPending() {
		t.Errorf("Unfinished requests")
		return
	}
}

func TestPlugin_FilterRepo(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "secretsmanager.GetSecretValue").
		Reply(200).
		File("testdata/secret.json")

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           "https://ec2.us-east-1.amazonaws.com",
		SigningRegion: config.Region,
	})

	manager := secretsmanager.New(config)
	req := &secret.Request{
		Name: "docker#username",
		Build: secret.Build{
			Event: "push",
		},
		Repo: secret.Repo{
			Slug: "spaceghost/hello-world",
		},
	}
	plugin := New(manager)
	_, err := plugin.Handle(req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "access denied: repository does not match"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
		return
	}

	if gock.IsPending() {
		t.Errorf("Unfinished requests")
		return
	}
}

func TestPlugin_FilterEvent(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "secretsmanager.GetSecretValue").
		Reply(200).
		File("testdata/secret.json")

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           "https://ec2.us-east-1.amazonaws.com",
		SigningRegion: config.Region,
	})

	manager := secretsmanager.New(config)
	req := &secret.Request{
		Name: "docker#username",
		Build: secret.Build{
			Event: "pull_request",
		},
		Repo: secret.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin := New(manager)
	_, err := plugin.Handle(req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "access denied: event does not match"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
		return
	}

	if gock.IsPending() {
		t.Errorf("Unfinished requests")
		return
	}
}

func TestPlugin_InvalidName(t *testing.T) {
	req := &secret.Request{
		Name: "secrets/docker/username",
	}
	_, err := New(nil).Handle(req)
	if err == nil {
		t.Errorf("Expect invalid path error")
		return
	}
	if got, want := err.Error(), "invalid or missing secret key"; got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestPlugin_NotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "secretsmanager.GetSecretValue").
		Reply(404).Done()

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           "https://ec2.us-east-1.amazonaws.com",
		SigningRegion: config.Region,
	})

	manager := secretsmanager.New(config)
	req := &secret.Request{
		Name: "docker#username",
		Build: secret.Build{
			Event: "pull_request",
		},
		Repo: secret.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin := New(manager)
	_, err := plugin.Handle(req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "secret not found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
		return
	}

	if gock.IsPending() {
		t.Errorf("Unfinished requests")
		return
	}
}
