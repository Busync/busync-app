package main

import (
	"fmt"

	"github.com/fatih/structs"
	"gopkg.in/h2non/gock.v1"
)

type appMockConfig struct {
	httpClientType string
	authConfig     interface{}
	isBusy         bool
	statusCode     int
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

func mockHTTPRequestWithBasicAuth(mockedRequest *gock.Request, authConfig HTTPBasicAuthConfig) *gock.Request {
	return mockedRequest.MatchHeader(
		"Authorization",
		"Basic "+basicAuth(authConfig.Username, authConfig.Password),
	)
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

func GetMockedApp(appName string, mockConfig appMockConfig) (app, error) {
	jsonMap := MockJSONMapResponseForGivenApp(appName, mockConfig.isBusy)

	mockedRequest := GetAPIMockedRequestForGivenApp(appName)
	mockedRequest = MockRequestWithGivenAuth(
		mockedRequest,
		mockConfig.httpClientType,
		mockConfig.authConfig,
	)
	mockedRequest.Reply(mockConfig.statusCode).
		JSON(jsonMap)

	httpClient, err := NewHTTPClient(mockConfig.httpClientType, mockConfig.authConfig)
	if err != nil {
		panic(err)
	}

	return NewApp(appName, httpClient)
}

func GetSliceOfMockedApps(appsNameAndConfig map[string]appMockConfig) []app {
	apps := make([]app, 0)

	for appName, mockConfig := range appsNameAndConfig {
		app, err := GetMockedApp(appName, mockConfig)
		if err != nil {
			panic(err)
		}

		apps = append(apps, app)
	}

	return apps
}
