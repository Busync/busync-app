package busylight_sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClientNotImplemented(t *testing.T) {
	assert := assert.New(t)
	authType := "foobar"

	got, err := NewHTTPClient(authType)

	assert.Nil(got)
	assert.Error(err)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", authType))
}
