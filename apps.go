package busylight_sync

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	FAKEAPP_API_URL = "http://fake.app/api"
	FAKEAPP_PATH    = "/is-busy"
)

type app interface {
	isBusy() (bool, error)
}

func NewApp(appName string, client *http.Client) (app, error) {
	switch appName {
	case "fake":
		return &FakeApp{client}, nil
	default:
		return nil, fmt.Errorf("%s is not implemented", appName)
	}
}

type FakeAppJSONResponse struct {
	IsBusy bool `json:"isBusy"`
}

type FakeApp struct {
	client *http.Client
}

func (f *FakeApp) isBusy() (bool, error) {
	resp, err := f.client.Get(FAKEAPP_API_URL + FAKEAPP_PATH)
	if err != nil {
		return false, err
	}

	var respJSON FakeAppJSONResponse
	err = json.NewDecoder(resp.Body).Decode(&respJSON)
	if err != nil {
		return false, err
	}

	return respJSON.IsBusy, nil
}
