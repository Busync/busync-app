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
	sut, err := NewApp(appName)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", appName))
}

func TestFakeApp(t *testing.T) {
	var testCases = []struct {
		desc       string
		statusCode int
		json       FakeAppJSONResponse
	}{
		{desc: "is busy", statusCode: 200, json: FakeAppJSONResponse{IsBusy: true}},
		{desc: "is not busy", statusCode: 200, json: FakeAppJSONResponse{IsBusy: false}},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			jsonMap := structs.Map(tC.json)

			gock.New(FAKEAPP_API_URL).
				Get(FAKEAPP_PATH).
				Reply(tC.statusCode).
				JSON(jsonMap)
			defer gock.Off()

			sut, err := NewApp("fake")
			assert.NoError(err)

			res, err := sut.isBusy()

			assert.NoError(err, "isBusy() method should not raise an error")
			assert.Equal(res, tC.json.IsBusy)
		})
	}
}
