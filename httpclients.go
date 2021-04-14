package busylight_sync

import (
	"fmt"
	"net/http"
)

func NewHTTPClient(authType string) (*http.Client, error) {
	switch authType {
	case "no-auth":
		return &http.Client{}, nil
	default:
		return nil, fmt.Errorf("%s is not implemented", authType)
	}
}
