package busylight_sync

import (
	"fmt"
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

func mockHTTPRequestWithBasicAuth(mockedRequest *gock.Request, authConfig HTTPBasicAuthConfig) *gock.Request {
	return mockedRequest.MatchHeader(
		"Authorization",
		"Basic "+basicAuth(authConfig.username, authConfig.password),
	)
}

func TestFakeApp(t *testing.T) {
	testCases := []struct {
		desc           string
		httpClientType string
		authConfig     interface{}
		json           FakeAppJSONResponse
		statusCode     int
	}{
		{
			desc:           "is busy without auth",
			httpClientType: "no-auth",
			authConfig:     nil,
			json:           FakeAppJSONResponse{IsBusy: true},
			statusCode:     200,
		},
		{
			desc:           "is not busy without auth",
			httpClientType: "no-auth",
			authConfig:     nil,
			json:           FakeAppJSONResponse{IsBusy: false},
			statusCode:     200,
		},
		{
			desc:           "is busy with basic auth",
			httpClientType: "basic-auth",
			authConfig:     HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			json:           FakeAppJSONResponse{IsBusy: true},
			statusCode:     200,
		},
		{
			desc:           "is not busy with basic auth",
			httpClientType: "basic-auth",
			authConfig:     HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			json:           FakeAppJSONResponse{IsBusy: false},
			statusCode:     200,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			jsonMap := structs.Map(tC.json)

			defer gock.Off()
			mockedRequest := gock.New(FAKEAPP_API_URL).
				Get(FAKEAPP_PATH)

			switch tC.httpClientType {
			case "basic-auth":
				mockedRequest = mockHTTPRequestWithBasicAuth(mockedRequest, tC.authConfig.(HTTPBasicAuthConfig))
			}

			mockedRequest.Reply(tC.statusCode).
				JSON(jsonMap)

			httpClient, err := NewHTTPClient(tC.httpClientType, tC.authConfig)
			assert.NoError(err)

			sut, err := NewApp("fake", httpClient)
			assert.NoError(err)

			res, err := sut.isBusy()
			assert.NoError(err, "isBusy() method should not raise an error")
			assert.Equal(res, tC.json.IsBusy)
		})
	}
}
