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
	httpClient, err := NewHTTPClient("no-auth", nil)
	assert.NoError(err)

	sut, err := NewApp(appName, httpClient)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", appName))
}

func TestAppsGetBusyStateFromJSONMapResponse(t *testing.T) {
	httpClient := &http.Client{}

	testCases := map[app][]struct {
		desc         string
		jsonResponse interface{}
		want         bool
	}{
		FakeApp{httpClient}: {
			{
				desc:         "is busy",
				jsonResponse: FakeAppJSONResponse{IsBusy: true},
				want:         true,
			},
			{
				desc:         "is not busy",
				jsonResponse: FakeAppJSONResponse{IsBusy: false},
				want:         false,
			},
		},
		Toggl{httpClient}: {
			{
				desc: "is busy",
				jsonResponse: TogglJSONResponse{Data: struct {
					Id int `json:"id"`
				}{1}},
				want: true,
			},
			{
				desc: "is not busy",
				jsonResponse: TogglJSONResponse{Data: struct {
					Id int `json:"id"`
				}{0}},
				want: false,
			},
		},
	}
	for sut, tCs := range testCases {
		appName := GetStructName(sut)
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
				authConfig:     HTTPBasicAuthConfig{Username: "foobar", Password: "spameggs"},
				isBusy:         true,
				statusCode:     200,
			},
		},
		{
			desc: "is not busy with basic auth",
			mockConfig: appMockConfig{
				httpClientType: "basic-auth",
				authConfig:     HTTPBasicAuthConfig{Username: "foobar", Password: "spameggs"},
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

				sut, err := GetMockedApp(appName, tC.mockConfig)
				assert.NoError(err)

				res, err := sut.isBusy()
				assert.NoError(err, "isBusy() method should not raise an error")
				assert.Equal(res, tC.mockConfig.isBusy)
			})
		}
	}
}
