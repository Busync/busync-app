package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func convertAppSliceToInterfaceSlice(apps []app) []interface{} {
	var interfaceSlice []interface{} = make([]interface{}, len(apps))
	for i, app := range apps {
		interfaceSlice[i] = app
	}

	return interfaceSlice
}

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

func TestGetAppsFromGivenConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		config   *Config
		wantApps []string
		wantErr  string
	}{
		{
			desc:    "config is nil",
			wantErr: "given config is nil",
		},
		{
			desc:    "no app",
			config:  &Config{},
			wantErr: "no app could be loaded from given config",
		},
		{
			desc: "one app with no auth",
			config: &Config{
				Apps{
					"toggl": {},
				},
			},
			wantErr: "no app could be loaded from given config",
		},
		{
			desc: "one non implemented app",
			config: &Config{
				Apps{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
			wantErr: "no app could be loaded from given config",
		},
		{
			desc: "one app",
			config: &Config{
				Apps{
					"fake": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
			wantApps: []string{"FakeApp"},
		},
		{
			desc: "two apps",
			config: &Config{
				Apps{
					"fake": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"toggl": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "barbaz",
							Password: "hamspam",
						},
					},
				},
			},
			wantApps: []string{"FakeApp", "Toggl"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			apps, err := getAppsFromGivenConfig(tC.config)
			appsInterfaceSlice := convertAppSliceToInterfaceSlice(apps)
			appsNames := GetSliceOfStructNames(appsInterfaceSlice)

			if err != nil {

				assert.Empty(apps)
				assert.EqualError(err, tC.wantErr)
			} else {
				assert.NoError(err)
				assert.ElementsMatch(tC.wantApps, appsNames)
			}
		})
	}
}
