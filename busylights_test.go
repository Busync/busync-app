package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotImplementedBusyLight(t *testing.T) {
	assert := assert.New(t)
	busylightName := "NotImplementedBusylight"

	sut, err := newBusyLight(busylightName)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s busylight is not implemented", busylightName))
}

func TestFakeDeviceStaticColor(t *testing.T) {
	testCases := []struct {
		desc  string
		color rgbColor
	}{
		{
			desc:  "led off",
			color: OffColor,
		},
		{
			desc:  "white",
			color: rgbColor{red: 255, green: 255, blue: 255},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			sut := fakeBusyLight{}
			sut.setStaticColor(tC.color)

			gotColor, err := sut.getStaticColor()
			assert.NoError(err)

			assert.Equal(tC.color, gotColor)
		})
	}
}

func TestFakeDeviceOff(t *testing.T) {
	testCases := []struct {
		desc             string
		previousRGBColor rgbColor
	}{
		{
			desc:             "already off",
			previousRGBColor: OffColor,
		},
		{
			desc:             "was white",
			previousRGBColor: rgbColor{red: 255, green: 255, blue: 255},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			sut := fakeBusyLight{}
			sut.setStaticColor(tC.previousRGBColor)

			sut.off()

			assert.Equal(OffColor, sut.color)
		})
	}
}
