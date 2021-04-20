package busylight_sync

import (
	"errors"
	"fmt"

	"github.com/google/gousb"
)

type USBDevice struct{}

func NewUSBDevice(vendorID, productID gousb.ID) (*USBDevice, error) {
	ctx := gousb.NewContext()
	dev, err := ctx.OpenDeviceWithVIDPID(vendorID, productID)
	if err != nil {
		return nil, err
	}

	if dev == nil {
		return nil, errors.New("device not found")
	}

	return nil, fmt.Errorf("test")
}
