package busylight_sync

var (
	OffColor RGBColor = RGBColor{red: 0, green: 0, blue: 0}
)

type Device interface {
	SetStaticColor(color RGBColor)
	Off()
}

type FakeDevice struct {
	color RGBColor
}

func (f *FakeDevice) SetStaticColor(color RGBColor) {
	f.color = color
}

func (f *FakeDevice) Off() {
	f.color = OffColor
}
