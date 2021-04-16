package busylight_sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeDeviceSetStaticColor(t *testing.T) {
	testCases := []struct {
		desc  string
		color RGBColor
	}{
		{
			desc:  "led off",
			color: RGBColor{red: 0, green: 0, blue: 0},
		},
		{
			desc:  "white",
			color: RGBColor{red: 255, green: 255, blue: 255},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			sut := FakeDevice{}
			sut.SetStaticColor(tC.color)

			assert.Equal(tC.color, sut.color)
		})
	}
}
