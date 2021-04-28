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

func AllBusylightsAreInGivenBusyState(busylights []BusyLight, expectIsBusy bool) bool {
	colors := GetColorOfAllGivenBusylights(busylights)
	busyStates := GetBusyStatesFromColors(colors)
	for _, currentIsBusy := range busyStates {
		if currentIsBusy != expectIsBusy {
			return false
		}
	}

	return true
}
