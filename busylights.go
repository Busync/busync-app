package main

import (
	"errors"
	"fmt"

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
	GetStaticColor() (RGBColor, error)
	SetStaticColor(RGBColor) error
	Off()
}

func NewBusyLight(name string) (BusyLight, error) {
	switch name {
	case "luxafor-flag":
		return NewLuxaforFlag()
	case "fake-busylight":
		return &FakeBusyLight{}, nil
	default:
		return nil, fmt.Errorf("%s busylight is not implemented", name)
	}
}

type FakeBusyLight struct {
	color RGBColor
}

func (f FakeBusyLight) GetStaticColor() (RGBColor, error) {
	return f.color, nil
}

func (f *FakeBusyLight) SetStaticColor(color RGBColor) error {
	f.color = color

	return nil
}

func (f *FakeBusyLight) Off() {
	f.SetStaticColor(OffColor)
}

type LuxaforFlag struct {
	device *USBDevice
}

func NewLuxaforFlag() (BusyLight, error) {
	dev, err := NewUSBDevice(LUXAFOR_FLAG_VENDOR_ID, LUXAFOR_FLAG_PRODUCT_ID)
	if err != nil {
		return nil, err
	}

	return &LuxaforFlag{
		device: dev,
	}, nil
}

func (l *LuxaforFlag) SetStaticColor(color RGBColor) error {
	data := []byte{0x01, LUXAFOR_FLAG_ALL_LED, color.red, color.green, color.blue, 0, 0x0, 0x0}

	err := l.device.WriteCommand(data)
	return err
}

func (LuxaforFlag) GetStaticColor() (RGBColor, error) {
	return RGBColor{}, errors.New("GetStaticColor is not implemented on Luxafor Flag")
}

func (l *LuxaforFlag) Off() {
	l.SetStaticColor(OffColor)
}
