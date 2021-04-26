package main

import (
	"github.com/google/gousb"
)

const (
	LUXAFOR_FLAG_VENDOR_ID  gousb.ID = 0x04d8
	LUXAFOR_FLAG_PRODUCT_ID gousb.ID = 0xf372
	LUXAFOR_FLAG_ALL_LED    byte     = 0xff
)

var (
	OffColor RGBColor = RGBColor{red: 0, green: 0, blue: 0}
)

type BusyLight interface {
	SetStaticColor(color RGBColor)
	Off()
}

type FakeBusyLight struct {
	color RGBColor
}

func (f *FakeBusyLight) SetStaticColor(color RGBColor) {
	f.color = color
}

func (f *FakeBusyLight) Off() {
	f.color = OffColor
}

type LuxaforFlag struct {
	device *USBDevice
}

func NewLuxaforFlag() (*LuxaforFlag, error) {
	dev, err := NewUSBDevice(LUXAFOR_FLAG_VENDOR_ID, LUXAFOR_FLAG_PRODUCT_ID)
	if err != nil {
		return nil, err
	}

	return &LuxaforFlag{
		device: dev,
	}, nil
}

func (l LuxaforFlag) SetStaticColor(color RGBColor) error {
	data := []byte{0x01, LUXAFOR_FLAG_ALL_LED, color.red, color.green, color.blue, 0, 0x0, 0x0}

	err := l.device.WriteCommand(data)
	return err
}

func (l *LuxaforFlag) Off() {
	l.SetStaticColor(OffColor)
}
