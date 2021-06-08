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

func getAPIMockedRequestForGivenApp(appName string) *gock.Request {
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

func mockHTTPRequestWithBasicAuth(mockedRequest *gock.Request, authConfig httpBasicAuthConfig) *gock.Request {
	return mockedRequest.MatchHeader(
		"Authorization",
		"Basic "+basicAuth(authConfig.Username, authConfig.Password),
	)
}

func mockRequestWithGivenAuth(req *gock.Request, httpClientType string, authConfig interface{}) *gock.Request {
	switch httpClientType {
	case "no-auth":
		return req
	case "basic-auth":
		return mockHTTPRequestWithBasicAuth(req, authConfig.(httpBasicAuthConfig))
	default:
		panic(fmt.Sprintf("%s mock http client not implemented", httpClientType))
	}
}

func getFakeAppMockedJSONMapResponse(isBusy bool) map[string]interface{} {
	mockedJSONResponse := fakeAppJSONResponse{IsBusy: isBusy}
	return structs.Map(mockedJSONResponse)
}

func getTogglMockedJSONMapResponse(isBusy bool) map[string]interface{} {
	var id int
	if isBusy {
		id = 12345
	} else {
		id = 0
	}

	mockedJSONResponse := togglJSONResponse{Data: struct {
		Id int `json:"id"`
	}{Id: id}}
	return structs.Map(mockedJSONResponse)
}

func mockJSONMapResponseForGivenApp(appName string, isBusy bool) map[string]interface{} {
	switch appName {
	case "fake":
		return getFakeAppMockedJSONMapResponse(isBusy)
	case "toggl":
		return getTogglMockedJSONMapResponse(isBusy)
	default:
		panic(fmt.Sprintf("%s app mock json map response not implemented", appName))
	}
}

func getMockedApp(appName string, mockConfig appMockConfig) (busyApps, error) {
	jsonMap := mockJSONMapResponseForGivenApp(appName, mockConfig.isBusy)

	mockedRequest := getAPIMockedRequestForGivenApp(appName)
	mockedRequest = mockRequestWithGivenAuth(
		mockedRequest,
		mockConfig.httpClientType,
		mockConfig.authConfig,
	)
	mockedRequest.Reply(mockConfig.statusCode).
		JSON(jsonMap)

	httpClient, err := newHTTPClient(mockConfig.httpClientType, mockConfig.authConfig)
	if err != nil {
		panic(err)
	}

	return newBusyApp(appName, httpClient)
}

func getSliceOfMockedApps(appsNameAndConfig map[string]appMockConfig) []busyApps {
	apps := make([]busyApps, 0)

	for appName, mockConfig := range appsNameAndConfig {
		app, err := getMockedApp(appName, mockConfig)
		if err != nil {
			panic(err)
		}

		apps = append(apps, app)
	}

	return apps
}
