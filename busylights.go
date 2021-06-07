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
	OffColor        rgbColor = rgbColor{red: 0, green: 0, blue: 0}
	BusyColor       rgbColor = rgbColor{red: 255, green: 0, blue: 0}
	UnoccupiedColor rgbColor = rgbColor{red: 0, green: 255, blue: 0}
)

type busyLight interface {
	getStaticColor() (rgbColor, error)
	setStaticColor(rgbColor) error
	off()
}

func newBusyLight(name string) (busyLight, error) {
	switch name {
	case "luxafor-flag":
		return newLuxaforFlag()
	case "fake-busylight":
		return &fakeBusyLight{}, nil
	default:
		return nil, fmt.Errorf("%s busylight is not implemented", name)
	}
}

type fakeBusyLight struct {
	color rgbColor
}

func (f fakeBusyLight) getStaticColor() (rgbColor, error) {
	return f.color, nil
}

func (f *fakeBusyLight) setStaticColor(color rgbColor) error {
	f.color = color

	return nil
}

func (f *fakeBusyLight) off() {
	f.setStaticColor(OffColor)
}

type luxaforFlag struct {
	device *usbDevice
}

func newLuxaforFlag() (busyLight, error) {
	dev, err := newUSBDevice(LUXAFOR_FLAG_VENDOR_ID, LUXAFOR_FLAG_PRODUCT_ID)
	if err != nil {
		return nil, err
	}

	return &luxaforFlag{
		device: dev,
	}, nil
}

func (l *luxaforFlag) setStaticColor(color rgbColor) error {
	data := []byte{0x01, LUXAFOR_FLAG_ALL_LED, color.red, color.green, color.blue, 0, 0x0, 0x0}

	err := l.device.writeCommand(data)
	return err
}

func (luxaforFlag) getStaticColor() (rgbColor, error) {
	return rgbColor{}, errors.New("GetStaticColor is not implemented on Luxafor Flag")
}

func (l *luxaforFlag) off() {
	l.setStaticColor(OffColor)
}
