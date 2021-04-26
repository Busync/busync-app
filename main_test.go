package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPClientFromAppConfig(t *testing.T) {
	testCases := []struct {
		desc      string
		appConfig AppConfig
		wantErr   string
	}{
		{
			desc:    "app config is empty",
			wantErr: "given app config is empty",
		},
		{
			desc: "got http client from basic auth",
			appConfig: AppConfig{
				BasicAuth: HTTPBasicAuthConfig{
					Username: "foobar",
					Password: "spameggs",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			httpClient, err := getHTTPClientFromAppConfig(tC.appConfig)

			if err != nil {
				assert.Nil(httpClient)
				assert.EqualError(err, tC.wantErr)
			} else {
				assert.NoError(err)
				assert.IsType(&http.Client{}, httpClient)
			}
		})
	}
}
