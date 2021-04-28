package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/afero"
)

var (
	fs                    = afero.Afero{Fs: afero.NewOsFs()}
	getConfigDir          = os.UserHomeDir
	busyLightsToTryToOpen = []string{
		"luxafor-flag",
	}
)

func tryToGetGivenBusylights(busylightNames []string) ([]BusyLight, error) {
	if len(busylightNames) == 0 {
		return nil, errors.New("no busylights on given slice")
	}

	busylights := make([]BusyLight, 0)
	for _, busylightName := range busylightNames {
		openedBusylight, err := NewBusyLight(busylightName)
		if err != nil {
			log.Printf("error when trying to open the busylight %s: %s", busylightName, err)
		} else {
			busylights = append(busylights, openedBusylight)
		}
	}

	if len(busylights) == 0 {
		return nil, errors.New("no busylight found")
	}

	return busylights, nil
}

func getHTTPClientFromAppConfig(config AppConfig) (*http.Client, error) {
	if config.BasicAuth.isNotEmpty() {
		return NewHTTPClient("basic-auth", config.BasicAuth)
	}

	return nil, errors.New("given app config is empty")
}

func getAppsFromGivenConfig(config *Config) ([]app, error) {
	apps := make([]app, 0)

	if config == nil {
		return apps, errors.New("given config is nil")
	}

	for appName, appConfig := range config.Apps {

		httpClient, err := getHTTPClientFromAppConfig(appConfig)
		if err != nil {
			log.Printf("error when trying to load %s app http client: %s", appName, err)
			continue
		}

		app, err := NewApp(appName, httpClient)
		if err != nil {
			log.Printf("error when trying to load %s app: %s", appName, err)
		} else {
			apps = append(apps, app)
		}
	}

	if len(apps) == 0 {
		return apps, errors.New("no app could be loaded from given config")
	} else {
		return apps, nil
	}
}

func AnyOfGivenAppIsBusy(apps []app) bool {
	for _, app := range apps {
		isBusy, err := app.isBusy()
		if err != nil {
			log.Println(err)
		}
		if isBusy {
			return true
		}
	}

	return false
}

func ChangeBusyStateOfAllGivenBusylights(isBusy bool, busylights []BusyLight) error {
	if len(busylights) == 0 {
		return errors.New("no busylights has been given to change their states")
	}

	for _, busylight := range busylights {
		if isBusy {
			busylight.SetStaticColor(BusyColor)
		} else {
			busylight.SetStaticColor(UnoccupiedColor)
		}
	}
	return nil
}

func AdaptBusylightsBusyStateAccordingToBusyStateOfApps(busylights []BusyLight, apps []app, wasBusy bool) (bool, error) {
	if len(busylights) == 0 {
		return wasBusy, errors.New("no busylights on given slice")
	}

	if len(apps) == 0 {
		return wasBusy, errors.New("no apps on given slice")
	}

	isBusy := AnyOfGivenAppIsBusy(apps)

	if isBusy != wasBusy {
		err := ChangeBusyStateOfAllGivenBusylights(isBusy, busylights)
		if err != nil {
			log.Println(err)
		}
	}

	return isBusy, nil
}

func main() {
	configDir, err := getConfigDir()
	if err != nil {
		panic(err)
	}

	config, err := LoadConfigFileFromDir(fs, configDir)
	if err != nil {
		panic(err)
	}

	busylights, err := tryToGetGivenBusylights(busyLightsToTryToOpen)
	if err != nil {
		panic(err)
	}

	apps, err := getAppsFromGivenConfig(config)
	if err != nil {
		panic(err)
	}

	// Init busy state of busylights
	wasBusy := false
	ChangeBusyStateOfAllGivenBusylights(false, busylights)
	if err != nil {
		log.Println(err)
	}

	for {
		wasBusy, err = AdaptBusylightsBusyStateAccordingToBusyStateOfApps(busylights, apps, wasBusy)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}
