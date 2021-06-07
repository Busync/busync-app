package main

import (
	"errors"

	"github.com/google/gousb"
)

type usbDevice struct {
	context     *gousb.Context
	outEndpoint *gousb.OutEndpoint
	closer      func()
}

func newUSBDevice(vendorID, productID gousb.ID) (*usbDevice, error) {
	ctx := gousb.NewContext()
	dev, err := ctx.OpenDeviceWithVIDPID(vendorID, productID)
	if err != nil {
		return nil, err
	}

	if dev == nil {
		return nil, errors.New("device not found")
	}

	if err := dev.SetAutoDetach(true); err != nil {
		return nil, err
	}

	iface, closer, err := dev.DefaultInterface()
	if err != nil {
		return nil, err
	}

	outEndpoint, err := iface.OutEndpoint(0x01)
	if err != nil {
		return nil, err
	}

	return &usbDevice{
		context:     ctx,
		outEndpoint: outEndpoint,
		closer:      closer,
	}, nil
}

func (u usbDevice) close() error {
	u.closer()
	err := u.context.Close()

	return err
}

func (u *usbDevice) writeCommand(command []byte) error {
	_, err := u.outEndpoint.Write(command)

	return err
}
