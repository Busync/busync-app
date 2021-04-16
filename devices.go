package busylight_sync

type Device interface {
	SetStaticColor(color RGBColor)
}

type FakeDevice struct {
	color RGBColor
}

func (f *FakeDevice) SetStaticColor(color RGBColor) {
	f.color = color
}
