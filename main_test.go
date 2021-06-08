package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestTryToGetGivenBusylights(t *testing.T) {
	testCases := []struct {
		desc                  string
		busyLightsToTryToOpen []string
		wantBusylights        []string
		wantErr               string
	}{
		{
			desc:    "empty slice",
			wantErr: "no busylights on given slice",
		},
		{
			desc:                  "one busylights",
			busyLightsToTryToOpen: []string{"fake-busylight"},
			wantBusylights:        []string{"fakeBusyLight"},
		},
		{
			desc:                  "two busylights",
			busyLightsToTryToOpen: []string{"fake-busylight", "fake-busylight"},
			wantBusylights:        []string{"fakeBusyLight", "fakeBusyLight"},
		},
		{
			desc:                  "one busylight and zero found",
			busyLightsToTryToOpen: []string{"foobar"},
			wantErr:               "no busylight found",
		},
		{
			desc:                  "two busylights and zero found",
			busyLightsToTryToOpen: []string{"foobar", "spameggs"},
			wantErr:               "no busylight found",
		},
		{
			desc:                  "two busylights and one found",
			busyLightsToTryToOpen: []string{"fake-busylight", "spameggs"},
			wantBusylights:        []string{"fakeBusyLight"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			busylights, err := tryToGetGivenBusylights(tC.busyLightsToTryToOpen)
			busylightsNames := getBusylightNames(busylights)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
				assert.Nil(busylights)
			} else {
				assert.NoError(err)
				assert.Equal(tC.wantBusylights, busylightsNames)
			}
		})
	}
}

func TestGetHTTPClientFromAppConfig(t *testing.T) {
	testCases := []struct {
		desc      string
		appConfig appConfiguration
		wantErr   string
	}{
		{
			desc:    "app config is empty",
			wantErr: "given app config is empty",
		},
		{
			desc: "got http client from basic auth",
			appConfig: appConfiguration{
				BasicAuth: httpBasicAuthConfig{
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
		config   *configuration
		wantApps []string
		wantErr  string
	}{
		{
			desc:    "config is nil",
			wantErr: "given config is nil",
		},
		{
			desc:    "no app",
			config:  &configuration{},
			wantErr: "no app could be loaded from given config",
		},
		{
			desc: "one app with no auth",
			config: &configuration{
				Apps{
					"toggl": {},
				},
			},
			wantErr: "no app could be loaded from given config",
		},
		{
			desc: "one non implemented app",
			config: &configuration{
				Apps{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
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
			config: &configuration{
				Apps{
					"fake": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
			wantApps: []string{"fakeBusyApp"},
		},
		{
			desc: "two apps",
			config: &configuration{
				Apps{
					"fake": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"toggl": {
						BasicAuth: httpBasicAuthConfig{
							Username: "barbaz",
							Password: "hamspam",
						},
					},
				},
			},
			wantApps: []string{"fakeBusyApp", "togglBusyApp"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			apps, err := getAppsFromGivenConfig(tC.config)
			appsNames := getAppNames(apps)

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

			apps := getSliceOfMockedApps(tC.appsMockConfig)

			got := AnyOfGivenAppIsBusy(apps)

			assert.Equal(tC.want, got)
		})
	}
}

func TestChangeBusyStateOfAllGivenBusylights(t *testing.T) {
	testCases := []struct {
		desc       string
		isBusy     bool
		busylights []busyLight
		wantErr    string
	}{
		{
			desc:       "no busylights on given slice",
			busylights: []busyLight{},
			wantErr:    "no busylights has been given to change their states",
		},
		{
			desc: "one busylight from busy to unoccupied",
			busylights: []busyLight{
				&fakeBusyLight{color: BusyColor},
			},
			isBusy: false,
		},
		{
			desc: "one busylight stay busy",
			busylights: []busyLight{
				&fakeBusyLight{color: BusyColor},
			},
			isBusy: true,
		}, {
			desc: "one busylight stay unoccupied",
			busylights: []busyLight{
				&fakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: false,
		},

		{
			desc: "two busylights from busy to unnocupied with one staying unnocupied",
			busylights: []busyLight{
				&fakeBusyLight{color: BusyColor},
				&fakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: false,
		},
		{
			desc: "two busylights from unnocupied to busy with one staying busy",
			busylights: []busyLight{
				&fakeBusyLight{color: BusyColor},
				&fakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: true,
		},
		{
			desc: "one busylight from unoccupied to busy",
			busylights: []busyLight{
				&fakeBusyLight{color: UnoccupiedColor},
			},
			isBusy: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := ChangeBusyStateOfAllGivenBusylights(tC.isBusy, tC.busylights)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
			} else {
				assert.NoError(err)
				assert.True(allBusylightsAreInGivenBusyState(tC.busylights, tC.isBusy))
			}
		})
	}
}

func TestAdaptBusylightsBusyStateAccordingToBusyStateOfApps(t *testing.T) {
	testCases := []struct {
		desc           string
		busylights     []busyLight
		appsMockConfig map[string]appMockConfig
		wasBusy        bool
		wantIsBusy     bool
		wantErr        string
	}{
		{
			desc:       "no busylights",
			busylights: []busyLight{},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
				},
			},
			wasBusy:    true,
			wantIsBusy: true,
			wantErr:    "no busylights on given slice",
		},
		{
			desc: "no apps",
			busylights: []busyLight{
				&fakeBusyLight{},
			},
			appsMockConfig: map[string]appMockConfig{},
			wasBusy:        false,
			wantIsBusy:     false,
			wantErr:        "no apps on given slice",
		},
		{
			desc: "one busy busylight & one unoccupied app",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			wasBusy:    true,
			wantIsBusy: false,
		},
		{
			desc: "one unoccupied busylight & one busy app",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
			},
			wasBusy:    false,
			wantIsBusy: true,
		},
		{
			desc: "both the busylight and the app are busy",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
			},
			wasBusy:    true,
			wantIsBusy: true,
		},
		{
			desc: "both the busylight and the app are unoccupied",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			wasBusy:    false,
			wantIsBusy: false,
		},
		{
			desc: "one busy busylight and two apps busy and unoccupied",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
				"toggl": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			wasBusy:    true,
			wantIsBusy: true,
		},
		{
			desc: "one unoccupied busylight and two apps busy and unoccupied",
			busylights: []busyLight{
				&fakeBusyLight{
					color: BusyColor,
				},
			},
			appsMockConfig: map[string]appMockConfig{
				"fake": {
					httpClientType: "no-auth",
					isBusy:         true,
					statusCode:     200,
				},
				"toggl": {
					httpClientType: "no-auth",
					isBusy:         false,
					statusCode:     200,
				},
			},
			wasBusy:    false,
			wantIsBusy: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			defer gock.Off()

			apps := getSliceOfMockedApps(tC.appsMockConfig)

			isBusy, err := AdaptBusylightsBusyStateAccordingToBusyStateOfApps(tC.busylights, apps, tC.wasBusy)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
			} else {
				assert.NoError(err)
			}
			assert.Equal(tC.wantIsBusy, isBusy)

		})
	}
}
