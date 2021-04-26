package main

import (
	"errors"
	"log"
	"net/http"
)


func getHTTPClientFromAppConfig(config AppConfig) (*http.Client, error) {
	if config.BasicAuth.isNotEmpty() {
		return NewHTTPClient("basic-auth", config.BasicAuth)
	}

	return nil, errors.New("given app config is empty")
}

func getAppsFromGivenConfig(config *Config) ([]app, error) {
	apps := make([]app, 0)

	if config == nil {
		return apps, errors.New("given config is nil")
	}
	
	for appName, appConfig := range config.Apps {

		httpClient, err := getHTTPClientFromAppConfig(appConfig)
		if err != nil {
			log.Printf("error when trying to load %s app http client: %s", appName, err)
			continue
		}

		app, err := NewApp(appName, httpClient)
		if err != nil {
			log.Printf("error when trying to load %s app: %s", appName, err)
		} else {
			apps = append(apps, app)
		}
	}

	if len(apps) == 0 {
		return apps, errors.New("no app could be loaded from given config")
	} else {
		return apps, nil
	}
}

func AnyOfGivenAppIsBusy(apps []app) bool {
	for _, app := range apps {
		isBusy, err := app.isBusy()
		if err != nil {
			log.Println(err)
		}
		if isBusy {
			return true
		}
	}

	return false
}
