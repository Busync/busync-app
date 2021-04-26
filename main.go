package main

import (
	"errors"
	"net/http"
)


func getHTTPClientFromAppConfig(config AppConfig) (*http.Client, error) {
	if config.BasicAuth.isNotEmpty() {
		return NewHTTPClient("basic-auth", config.BasicAuth)
	}

	return nil, errors.New("given app config is empty")
}

