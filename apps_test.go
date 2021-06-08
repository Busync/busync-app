package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestNotImplementedApp(t *testing.T) {
	assert := assert.New(t)
	appName := "NotImplementedApp"
	httpClient, err := newHTTPClient("no-auth", nil)
	assert.NoError(err)

	sut, err := newBusyApp(appName, httpClient)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", appName))
}

func TestAppsGetBusyStateFromJSONMapResponse(t *testing.T) {
	httpClient := &http.Client{}

	testCases := map[busyApps][]struct {
		desc         string
		jsonResponse interface{}
		want         bool
	}{
		fakeBusyApp{httpClient}: {
			{
				desc:         "is busy",
				jsonResponse: fakeAppJSONResponse{IsBusy: true},
				want:         true,
			},
			{
				desc:         "is not busy",
				jsonResponse: fakeAppJSONResponse{IsBusy: false},
				want:         false,
			},
		},
		togglBusyApp{httpClient}: {
			{
				desc: "is busy",
				jsonResponse: togglJSONResponse{Data: struct {
					Id int `json:"id"`
				}{1}},
				want: true,
			},
			{
				desc: "is not busy",
				jsonResponse: togglJSONResponse{Data: struct {
					Id int `json:"id"`
				}{0}},
				want: false,
			},
		},
	}
	for sut, tCs := range testCases {
		appName := getStructName(sut)
		for _, tC := range tCs {
			t.Run(fmt.Sprintf("%s/%s", appName, tC.desc), func(t *testing.T) {
				assert := assert.New(t)

				got := sut.getBusyStateFromJSONResponse(tC.jsonResponse)
				assert.Equal(tC.want, got)
			})
		}
	}
}

func TestAppsIsBusy(t *testing.T) {
	apps := []string{
		"fake",
		"toggl",
	}

	testCases := []struct {
		desc       string
		mockConfig appMockConfig
	}{
		{
			desc: "is busy without auth",
			mockConfig: appMockConfig{
				httpClientType: "no-auth",
				authConfig:     nil,
				isBusy:         true,
				statusCode:     200,
			},
		},
		{
			desc: "is not busy without auth",
			mockConfig: appMockConfig{

				httpClientType: "no-auth",
				authConfig:     nil,
				isBusy:         false,
				statusCode:     200,
			},
		},
		{
			desc: "is busy with basic auth",
			mockConfig: appMockConfig{
				httpClientType: "basic-auth",
				authConfig:     httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
				isBusy:         true,
				statusCode:     200,
			},
		},
		{
			desc: "is not busy with basic auth",
			mockConfig: appMockConfig{
				httpClientType: "basic-auth",
				authConfig:     httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
				isBusy:         false,
				statusCode:     200,
			},
		},
	}
	for _, appName := range apps {
		for _, tC := range testCases {
			t.Run(appName+"/"+tC.desc, func(t *testing.T) {
				assert := assert.New(t)
				defer gock.Off()

				sut, err := getMockedApp(appName, tC.mockConfig)
				assert.NoError(err)

				res, err := sut.isBusy()
				assert.NoError(err, "isBusy() method should not raise an error")
				assert.Equal(res, tC.mockConfig.isBusy)
			})
		}
	}
}
