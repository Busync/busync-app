package main

func GetAppNames(apps []app) []string {
	var appNames = make([]string, 0)
	for _, currentApp := range apps {
		appName := GetStructName(currentApp)
		appNames = append(appNames, appName)
	}

	return appNames
}
