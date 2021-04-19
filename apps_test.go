package busylight_sync

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

func TestAppsGetBusyStateFromJSONResponse(t *testing.T) {
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
		"Basic "+basicAuth(authConfig.username, authConfig.password),
	)
}

func GetFakeAppMockedRequest(httpClientType string, authConfig interface{}) *gock.Request {
	mockedRequest := gock.New(FAKEAPP_API_URL).
		Get(FAKEAPP_PATH)

	switch httpClientType {
	case "no-auth":
		return mockedRequest
	case "basic-auth":
		return mockHTTPRequestWithBasicAuth(mockedRequest, authConfig.(HTTPBasicAuthConfig))
	default:
		panic(fmt.Sprintf("%s http client not implemented", httpClientType))
	}
}

func GetFakeAppMockedJSONMapResponse(isBusy bool) map[string]interface{} {
	mockedJSONResponse := FakeAppJSONResponse{IsBusy: isBusy}
	return structs.Map(mockedJSONResponse)
}

func TestAppsIsBusy(t *testing.T) {
	apps := map[string]struct {
		getMockedRequest         func(string, interface{}) *gock.Request
		getMockedJSONMapResponse func(bool) map[string]interface{}
	}{
		"fake": {
			getMockedRequest:         GetFakeAppMockedRequest,
			getMockedJSONMapResponse: GetFakeAppMockedJSONMapResponse,
		},
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
			authConfig:     HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			isBusy:         true,
			statusCode:     200,
		},
		{
			desc:           "is not busy with basic auth",
			httpClientType: "basic-auth",
			authConfig:     HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			isBusy:         false,
			statusCode:     200,
		},
	}
	for appName, utils := range apps {
		for _, tC := range testCases {
			t.Run(appName+"/"+tC.desc, func(t *testing.T) {
				assert := assert.New(t)
				jsonMap := utils.getMockedJSONMapResponse(tC.isBusy)

				defer gock.Off()
				mockedRequest := utils.getMockedRequest(tC.httpClientType, tC.authConfig)
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
