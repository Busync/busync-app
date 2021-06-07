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

func TestHTTPAuthConfigIsEmptyMethod(t *testing.T) {
	testCases := []struct {
		desc       string
		authConfig httpAuthConfig
		want       bool
	}{
		{
			desc:       "basic auth is empty",
			authConfig: httpBasicAuthConfig{},
			want:       false,
		},
		{
			desc: "basic auth is not empty",
			authConfig: httpBasicAuthConfig{
				Username: "foobar",
				Password: "spameggs",
			},
			want: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tC.want, tC.authConfig.isNotEmpty())
		})
	}
}

func TestHTTPClientNotImplemented(t *testing.T) {
	assert := assert.New(t)
	authType := "foobar"

	got, err := newHTTPClient(authType, nil)

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

	sut, err := newHTTPClient(authType, nil)
	assert.NoError(err)

	resp, err := sut.Get(url)
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
}

func TestHTTPClientBasicAuth(t *testing.T) {
	testCases := []struct {
		desc             string
		clientAuthConfig httpBasicAuthConfig
		serverAuthConfig httpBasicAuthConfig
		statusCode       int
	}{
		{
			desc:             "OK",
			clientAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			serverAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			statusCode:       200,
		},
		{
			desc:             "OK no password",
			clientAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: ""},
			serverAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: ""},
			statusCode:       200,
		},
		{
			desc:             "OK no username",
			clientAuthConfig: httpBasicAuthConfig{Username: "", Password: "spameggs"},
			serverAuthConfig: httpBasicAuthConfig{Username: "", Password: "spameggs"},
			statusCode:       200,
		},
		{
			desc:             "OK no username and password",
			clientAuthConfig: httpBasicAuthConfig{Username: "", Password: ""},
			serverAuthConfig: httpBasicAuthConfig{Username: "", Password: ""},
			statusCode:       200,
		},
		{
			desc:             "Unauthorized wrong password",
			clientAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			serverAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "barbaz"},
		},
		{
			desc:             "Unauthorized wrong username",
			clientAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			serverAuthConfig: httpBasicAuthConfig{Username: "hamspam", Password: "spameggs"},
		},
		{
			desc:             "Unauthorized wrong username and password",
			clientAuthConfig: httpBasicAuthConfig{Username: "foobar", Password: "spameggs"},
			serverAuthConfig: httpBasicAuthConfig{Username: "hamspam", Password: "barbaz"},
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
					"Basic "+basicAuth(tC.serverAuthConfig.Username, tC.serverAuthConfig.Password),
				).
				Reply(tC.statusCode)

			sut, err := newHTTPClient(authType, tC.clientAuthConfig)
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
