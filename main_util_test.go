package main

func GetBusylightNames(busylights []BusyLight) []string {
	var busylightNames = make([]string, 0)
	for _, currentBusylight := range busylights {
		busylightName := GetStructName(currentBusylight)
		busylightNames = append(busylightNames, busylightName)
	}

	return busylightNames
}

func GetAppNames(apps []app) []string {
	var appNames = make([]string, 0)
	for _, currentApp := range apps {
		appName := GetStructName(currentApp)
		appNames = append(appNames, appName)
	}

	return appNames
}
