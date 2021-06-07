package main

func getBusylightNames(busylights []busyLight) []string {
	var busylightNames = make([]string, 0)
	for _, currentBusylight := range busylights {
		busylightName := getStructName(currentBusylight)
		busylightNames = append(busylightNames, busylightName)
	}

	return busylightNames
}

func getAppNames(apps []busyApps) []string {
	var appNames = make([]string, 0)
	for _, currentApp := range apps {
		appName := getStructName(currentApp)
		appNames = append(appNames, appName)
	}

	return appNames
}
