package busylight_sync

import (
	"fmt"
)

type app interface {
	isBusy() (bool, error)
}

func NewApp(appName string) (app, error) {
	switch appName {
	default:
		return nil, fmt.Errorf("%s is not implemented", appName)
	}
}

