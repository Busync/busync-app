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

func GetBusyStatesFromColors(colors []RGBColor) (busyStates []bool) {
	busyStates = make([]bool, 0)
	for _, color := range colors {
		switch color {
		case BusyColor:
			busyStates = append(busyStates, true)
		default:
			busyStates = append(busyStates, false)
		}
	}

	return
}
