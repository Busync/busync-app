package busylight_sync

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

// Function from:  https://github.com/golang/go/blob/master/src/net/http/client.go#L417
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestHTTPClientNotImplemented(t *testing.T) {
	assert := assert.New(t)
	authType := "foobar"

	got, err := NewHTTPClient(authType, nil)

	assert.Nil(got)
	assert.Error(err)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", authType))
}

func TestHTTPClientNoAuth(t *testing.T) {
	assert := assert.New(t)
	authType := "no-auth"
	url := "https://foo.bar"

	defer gock.Off()
	gock.New(url).
		Reply(200)

	sut, err := NewHTTPClient(authType, nil)
	assert.NoError(err)

	resp, err := sut.Get(url)
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
}

func TestHTTPClientBasicAuth(t *testing.T) {
	testCases := []struct {
		desc       string
		authConfig HTTPBasicAuthConfig
		statusCode int
	}{
		{
			desc:       "OK",
			authConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			statusCode: 200,
		},
		{
			desc:       "OK no password",
			authConfig: HTTPBasicAuthConfig{username: "foobar", password: ""},
			statusCode: 200,
		},
		{
			desc:       "OK no username",
			authConfig: HTTPBasicAuthConfig{username: "", password: "spameggs"},
			statusCode: 200,
		},
		{
			desc:       "OK no username and password",
			authConfig: HTTPBasicAuthConfig{username: "", password: ""},
			statusCode: 200,
		},
	}
	for _, tC := range testCases {
		authType := "basic-auth"
		url := "https://foo.bar"

		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			defer gock.Off()
			gock.New(url).
				MatchHeader(
					"Authorization",
					"Basic "+basicAuth(tC.authConfig.username, tC.authConfig.password),
				).
				Reply(tC.statusCode)

			sut, err := NewHTTPClient(authType, tC.authConfig)
			assert.NoError(err)

			got, err := sut.Get(url)
		})
	}
}
