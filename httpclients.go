package main

import (
	"fmt"
	"net/http"
)

type HTTPBasicAuthConfig struct {
	username string
	password string
}

type HTTPBasicAuthRoundTripper struct {
	rt     http.RoundTripper
	config HTTPBasicAuthConfig
}

func (brt HTTPBasicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(brt.config.username, brt.config.password)

	return brt.rt.RoundTrip(req)
}

func NewBasicAuthClient(config HTTPBasicAuthConfig) *http.Client {
	roundTripper := HTTPBasicAuthRoundTripper{
		config: config,
		rt:     http.DefaultTransport,
	}

	return &http.Client{
		Transport: roundTripper,
	}
}

func NewHTTPClient(authType string, authConfig interface{}) (*http.Client, error) {
	switch authType {
	case "no-auth":
		return &http.Client{}, nil
	case "basic-auth":
		return NewBasicAuthClient(authConfig.(HTTPBasicAuthConfig)), nil
	default:
		return nil, fmt.Errorf("%s is not implemented", authType)
	}
}
