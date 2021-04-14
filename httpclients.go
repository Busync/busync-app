package busylight_sync

import (
	"fmt"
	"net/http"
)

func NewHTTPClient(authType string) (*http.Client, error) {
	switch authType {
	default:
		return nil, fmt.Errorf("%s is not implemented", authType)
	}
}
