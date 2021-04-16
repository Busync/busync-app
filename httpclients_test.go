package busylight_sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

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
