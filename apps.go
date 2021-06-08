package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	FAKEAPP_API_URL  = "http://fake.app/api"
	FAKEAPP_API_PATH = "/is-busy"
	TOGGL_API_URL    = "https://api.track.toggl.com/api/v8"
	TOGGL_API_PATH   = "/time_entries/current"
)

type busyApps interface {
	getBusyStateFromJSONResponse(interface{}) bool
	isBusy() (bool, error)
}

func newBusyApp(appName string, client *http.Client) (busyApps, error) {
	switch appName {
	case "fake":
		return &fakeBusyApp{client}, nil
	case "toggl":
		return &togglBusyApp{client}, nil
	default:
		return nil, fmt.Errorf("%s is not implemented", appName)
	}
}

type fakeAppJSONResponse struct {
	IsBusy bool `json:"isBusy"`
}

type fakeBusyApp struct {
	client *http.Client
}

func (fakeBusyApp) getBusyStateFromJSONResponse(jsonResponse interface{}) bool {
	return jsonResponse.(fakeAppJSONResponse).IsBusy
}

func (f fakeBusyApp) isBusy() (bool, error) {
	resp, err := f.client.Get(FAKEAPP_API_URL + FAKEAPP_API_PATH)
	if err != nil {
		return false, err
	}

	var jsonResponse fakeAppJSONResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return false, err
	}

	return f.getBusyStateFromJSONResponse(jsonResponse), nil
}

type togglJSONResponse struct {
	Data struct {
		Id int `json:"id"`
	} `json:"data"`
}

type togglBusyApp struct {
	client *http.Client
}

func (togglBusyApp) getBusyStateFromJSONResponse(jsonResponse interface{}) bool {
	return jsonResponse.(togglJSONResponse).Data.Id != 0
}

func (t togglBusyApp) isBusy() (bool, error) {
	resp, err := t.client.Get(TOGGL_API_URL + TOGGL_API_PATH)
	if err != nil {
		return false, err
	}

	var jsonResponse togglJSONResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return false, err
	}

	return t.getBusyStateFromJSONResponse(jsonResponse), nil
}
