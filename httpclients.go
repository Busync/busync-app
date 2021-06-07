package main

import (
	"fmt"
	"net/http"
)

type httpAuthConfig interface {
	isNotEmpty() bool
}

type httpBasicAuthConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func (b httpBasicAuthConfig) isNotEmpty() bool {
	return b != httpBasicAuthConfig{}
}

type httpBasicAuthRoundTripper struct {
	rt     http.RoundTripper
	config httpBasicAuthConfig
}

func (brt httpBasicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(brt.config.Username, brt.config.Password)

	return brt.rt.RoundTrip(req)
}

func newBasicAuthClient(config httpBasicAuthConfig) *http.Client {
	roundTripper := httpBasicAuthRoundTripper{
		config: config,
		rt:     http.DefaultTransport,
	}

	return &http.Client{
		Transport: roundTripper,
	}
}

func newHTTPClient(authType string, authConfig interface{}) (*http.Client, error) {
	switch authType {
	case "no-auth":
		return &http.Client{}, nil
	case "basic-auth":
		return newBasicAuthClient(authConfig.(httpBasicAuthConfig)), nil
	default:
		return nil, fmt.Errorf("%s is not implemented", authType)
	}
}
