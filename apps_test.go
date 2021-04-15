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

func TestFakeApp(t *testing.T) {
	var testCases = []struct {
		desc           string
		statusCode     int
		json           FakeAppJSONResponse
		httpClientType string
		authConfig     interface{}
	}{
		{desc: "is busy", statusCode: 200, json: FakeAppJSONResponse{IsBusy: true}, httpClientType: "no-auth", authConfig: nil},
		{desc: "is not busy", statusCode: 200, json: FakeAppJSONResponse{IsBusy: false}, httpClientType: "no-auth", authConfig: nil},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			jsonMap := structs.Map(tC.json)
			httpClient, err := NewHTTPClient(tC.httpClientType, tC.authConfig)
			assert.NoError(err)

			gock.New(FAKEAPP_API_URL).
				Get(FAKEAPP_PATH).
				Reply(tC.statusCode).
				JSON(jsonMap)
			defer gock.Off()

			sut, err := NewApp("fake", httpClient)
			assert.NoError(err)

			res, err := sut.isBusy()

			assert.NoError(err, "isBusy() method should not raise an error")
			assert.Equal(res, tC.json.IsBusy)
		})
	}
}
