package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
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

func TestAnyOfGivenAppIsBusy(t *testing.T) {
	testCases := []struct {
		desc           string
		appsMockConfig map[string]appMockConfig
		want           bool
	}{
		{
			desc: "one app not busy",
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			want: false,
		},
		{
			desc: "one busy app",
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
			},
			want: true,
		},
		{
			desc: "two apps with one busy",
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
				"toggl": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
			},
			want: true,
		},
		{
			desc: "two apps not busy",
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
				"toggl": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			want: false,
		},
		{
			desc: "two busy apps",
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
				"toggl": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
			},
			want: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			defer gock.Off()

			apps := GetSliceOfMockedApps(tC.appsMockConfig)

			got := AnyOfGivenAppIsBusy(apps)

			assert.Equal(tC.want, got)
		})
	}
}

func TestChangeBusyStateOfAllGivenBusylights(t *testing.T) {
	testCases := []struct {
		desc       string
		isBusy     bool
		busylights []BusyLight
	}{
		{
			desc:       "no busylights on given slice",
			busylights: []BusyLight{},
		},
		{
			desc: "one busylight from busy to unoccupied",
			busylights: []BusyLight{
				&FakeBusyLight{color: BusyColor},
			},
			isBusy: false,
		},
		{
			desc: "one busylight stay busy",
			busylights: []BusyLight{
				&FakeBusyLight{color: BusyColor},
			},
			isBusy: true,
		}, {
			desc: "one busylight stay unoccupied",
			busylights: []BusyLight{
				&FakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: false,
		},

		{
			desc: "two busylights from busy to unnocupied with one staying unnocupied",
			busylights: []BusyLight{
				&FakeBusyLight{color: BusyColor},
				&FakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: false,
		},
		{
			desc: "two busylights from unnocupied to busy with one staying busy",
			busylights: []BusyLight{
				&FakeBusyLight{color: BusyColor},
				&FakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: true,
		},
		{
			desc: "one busylight from unoccupied to busy",
			busylights: []BusyLight{
				&FakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := ChangeBusyStateOfAllGivenBusylights(tC.isBusy, tC.busylights)

			if err != nil {
				assert.EqualError(err, "no busylights has been given to change their states")
			} else {
				assert.NoError(err)
				assert.True(AllBusylightsAreInGivenBusyState(tC.busylights, tC.isBusy))
			}
		})
	}
}
