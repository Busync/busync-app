package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotImplementedBusyLight(t *testing.T) {
	assert := assert.New(t)
	busylightName := "NotImplementedBusylight"

	sut, err := NewBusyLight(busylightName)

	assert.Nil(sut)
	assert.EqualError(err, fmt.Sprintf("%s busylight is not implemented", busylightName))
}

func TestFakeDeviceSetStaticColor(t *testing.T) {
	testCases := []struct {
		desc  string
		color RGBColor
	}{
		{
			desc:  "led off",
			color: OffColor,
		},
		{
			desc:  "white",
			color: RGBColor{red: 255, green: 255, blue: 255},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			sut := FakeBusyLight{}
			sut.SetStaticColor(tC.color)

			assert.Equal(tC.color, sut.color)
		})
	}
}

func TestFakeDeviceOff(t *testing.T) {
	testCases := []struct {
		desc             string
		previousRGBColor RGBColor
	}{
		{
			desc:             "already off",
			previousRGBColor: OffColor,
		},
		{
			desc:             "was white",
			previousRGBColor: RGBColor{red: 255, green: 255, blue: 255},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			sut := FakeBusyLight{}
			sut.SetStaticColor(tC.previousRGBColor)

			sut.Off()

			assert.Equal(OffColor, sut.color)
		})
	}
}
