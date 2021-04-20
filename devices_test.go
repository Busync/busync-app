package busylight_sync

import (
	"testing"

	"github.com/google/gousb"
	"github.com/stretchr/testify/assert"
)

func TestUSBDeviceOpen(t *testing.T) {
	testCases := []struct {
		desc      string
		vendorID  gousb.ID
		productID gousb.ID
		wantErr   string
	}{
		{
			desc:      "not found",
			vendorID:  0x0000,
			productID: 0x0000,
			wantErr:   "device not found",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			dev, err := NewUSBDevice(tC.vendorID, tC.productID)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
				assert.Nil(dev)
			}
		})
	}
}
