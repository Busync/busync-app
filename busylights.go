package busylight_sync

import (
	"fmt"

	"github.com/google/gousb"
)

const (
	LUXAFOR_FLAG_VENDOR_ID  gousb.ID = 0x04d8
	LUXAFOR_FLAG_PRODUCT_ID gousb.ID = 0xf372
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
