package main

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
		desc             string
		clientAuthConfig HTTPBasicAuthConfig
		serverAuthConfig HTTPBasicAuthConfig
		statusCode       int
	}{
		{
			desc:             "OK",
			clientAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			serverAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			statusCode:       200,
		},
		{
			desc:             "OK no password",
			clientAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: ""},
			serverAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: ""},
			statusCode:       200,
		},
		{
			desc:             "OK no username",
			clientAuthConfig: HTTPBasicAuthConfig{username: "", password: "spameggs"},
			serverAuthConfig: HTTPBasicAuthConfig{username: "", password: "spameggs"},
			statusCode:       200,
		},
		{
			desc:             "OK no username and password",
			clientAuthConfig: HTTPBasicAuthConfig{username: "", password: ""},
			serverAuthConfig: HTTPBasicAuthConfig{username: "", password: ""},
			statusCode:       200,
		},
		{
			desc:             "Unauthorized wrong password",
			clientAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			serverAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "barbaz"},
		},
		{
			desc:             "Unauthorized wrong username",
			clientAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			serverAuthConfig: HTTPBasicAuthConfig{username: "hamspam", password: "spameggs"},
		},
		{
			desc:             "Unauthorized wrong username and password",
			clientAuthConfig: HTTPBasicAuthConfig{username: "foobar", password: "spameggs"},
			serverAuthConfig: HTTPBasicAuthConfig{username: "hamspam", password: "barbaz"},
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
					"Basic "+basicAuth(tC.serverAuthConfig.username, tC.serverAuthConfig.password),
				).
				Reply(tC.statusCode)

			sut, err := NewHTTPClient(authType, tC.clientAuthConfig)
			assert.NoError(err)

			got, err := sut.Get(url)
			if err != nil {
				assert.Nil(got)
				assert.NotEqual(tC.serverAuthConfig, tC.clientAuthConfig)
			} else {
				assert.Equal(tC.serverAuthConfig, tC.clientAuthConfig)
				assert.Equal(tC.statusCode, got.StatusCode)
			}
		})
	}
}
