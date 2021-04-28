package main

func GetColorOfAllGivenBusylights(busylights []BusyLight) []RGBColor {
	colors := make([]RGBColor, 0)
	for _, busylight := range busylights {
		color, err := busylight.GetStaticColor()
		if err != nil {
			panic(err)
		}
		colors = append(colors, color)
	}

	return colors
}

