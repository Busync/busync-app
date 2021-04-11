package busylight_sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotImplementedApp(t *testing.T) {
	assert := assert.New(t)
	appName := "NotImplementedApp"
	sut, err := NewApp(appName)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s is not implemented", appName))
}

