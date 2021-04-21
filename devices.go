package busylight_sync

import (
	"errors"

	"github.com/google/gousb"
)

type USBDevice struct {
	context     *gousb.Context
	outEndpoint *gousb.OutEndpoint
	closer      func()
}

func NewUSBDevice(vendorID, productID gousb.ID) (*USBDevice, error) {
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

	return &USBDevice{
		context:     ctx,
		outEndpoint: outEndpoint,
		closer:      closer,
	}, nil
}

func (u USBDevice) Close() error {
	u.closer()
	err := u.context.Close()

	return err
}

func (u *USBDevice) WriteCommand(command []byte) error {
	_, err := u.outEndpoint.Write(command)

	return err
}
