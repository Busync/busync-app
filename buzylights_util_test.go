package main

func getColorOfAllGivenBusylights(busylights []busyLight) []rgbColor {
	colors := make([]rgbColor, 0)
	for _, busylight := range busylights {
		color, err := busylight.getStaticColor()
		if err != nil {
			panic(err)
		}
		colors = append(colors, color)
	}

	return colors
}

func getBusyStatesFromColors(colors []rgbColor) (busyStates []bool) {
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

func allBusylightsAreInGivenBusyState(busylights []busyLight, expectIsBusy bool) bool {
	colors := getColorOfAllGivenBusylights(busylights)
	busyStates := getBusyStatesFromColors(colors)
	for _, currentIsBusy := range busyStates {
		if currentIsBusy != expectIsBusy {
			return false
		}
	}

	return true
}
