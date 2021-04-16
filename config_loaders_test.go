package busylight_sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigInGivenFormat(t *testing.T) {
	testCases := []struct {
		desc       string
		fileFormat string
	}{
		{
			desc:       "file format not implemented",
			fileFormat: "foobar",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := LoadConfigInGivenFormat(tC.fileFormat)

			assert.Error(err)
			assert.EqualError(err, fmt.Sprintf("%s is not implemented", tC.fileFormat))
		})
	}
}
