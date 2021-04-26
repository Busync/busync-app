package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fatih/structs"
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
func mockHTTPRequestWithBasicAuth(mockedRequest *gock.Request, authConfig HTTPBasicAuthConfig) *gock.Request {
	return mockedRequest.MatchHeader(
		"Authorization",
		"Basic "+basicAuth(authConfig.Username, authConfig.Password),
	)
}

func GetAPIMockedRequestForGivenApp(appName string) *gock.Request {
	var url, path string
	switch appName {
	case "fake":
		url = FAKEAPP_API_URL
		path = FAKEAPP_API_PATH
	case "toggl":
		url = TOGGL_API_URL
		path = TOGGL_API_PATH
	default:
		panic(fmt.Sprintf("%s app mock api not implemented", appName))
	}

	return gock.New(url).
		Get(path)
}

func MockRequestWithGivenAuth(req *gock.Request, httpClientType string, authConfig interface{}) *gock.Request {
	switch httpClientType {
	case "no-auth":
		return req
	case "basic-auth":
		return mockHTTPRequestWithBasicAuth(req, authConfig.(HTTPBasicAuthConfig))
	default:
		panic(fmt.Sprintf("%s mock http client not implemented", httpClientType))
	}
}

func GetFakeAppMockedJSONMapResponse(isBusy bool) map[string]interface{} {
	mockedJSONResponse := FakeAppJSONResponse{IsBusy: isBusy}
	return structs.Map(mockedJSONResponse)
}

func GetTogglMockedJSONMapResponse(isBusy bool) map[string]interface{} {
	var id int
	if isBusy {
		id = 12345
	} else {
		id = 0
	}

	mockedJSONResponse := TogglJSONResponse{Data: struct {
		Id int `json:"id"`
	}{Id: id}}
	return structs.Map(mockedJSONResponse)
}

func MockJSONMapResponseForGivenApp(appName string, isBusy bool) map[string]interface{} {
	switch appName {
	case "fake":
		return GetFakeAppMockedJSONMapResponse(isBusy)
	case "toggl":
		return GetTogglMockedJSONMapResponse(isBusy)
	default:
		panic(fmt.Sprintf("%s app mock json map response not implemented", appName))
	}
}

func TestAppsIsBusy(t *testing.T) {
	apps := []string{
		"fake",
		"toggl",
	}

	testCases := []struct {
		desc           string
		httpClientType string
		authConfig     interface{}
		isBusy         bool
		statusCode     int
	}{
		{
			desc:           "is busy without auth",
			httpClientType: "no-auth",
			authConfig:     nil,
			isBusy:         true,
			statusCode:     200,
		},
		{
			desc:           "is not busy without auth",
			httpClientType: "no-auth",
			authConfig:     nil,
			isBusy:         false,
			statusCode:     200,
		},
		{
			desc:           "is busy with basic auth",
			httpClientType: "basic-auth",
			authConfig:     HTTPBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			isBusy:         true,
			statusCode:     200,
		},
		{
			desc:           "is not busy with basic auth",
			httpClientType: "basic-auth",
			authConfig:     HTTPBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			isBusy:         false,
			statusCode:     200,
		},
	}
	for _, appName := range apps {
		for _, tC := range testCases {
			t.Run(appName+"/"+tC.desc, func(t *testing.T) {
				assert := assert.New(t)
				jsonMap := MockJSONMapResponseForGivenApp(appName, tC.isBusy)

				defer gock.Off()
				mockedRequest := GetAPIMockedRequestForGivenApp(appName)
				mockedRequest = MockRequestWithGivenAuth(
					mockedRequest,
					tC.httpClientType,
					tC.authConfig,
				)
				mockedRequest.Reply(tC.statusCode).
					JSON(jsonMap)

				httpClient, err := NewHTTPClient(tC.httpClientType, tC.authConfig)
				assert.NoError(err)

				sut, err := NewApp(appName, httpClient)
				assert.NoError(err)

				res, err := sut.isBusy()
				assert.NoError(err, "isBusy() method should not raise an error")
				assert.Equal(res, tC.isBusy)
			})
		}
	}
}
