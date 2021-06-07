package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
)

const (
	TOML_CONFIG_FILE string = ".busylight-sync.toml"
)

var configFileFormats = map[string]struct {
	filename string
	load     func(afero.Afero, string) (*configuration, error)
}{
	"toml": {
		filename: TOML_CONFIG_FILE,
		load:     loadConfigFromTOMLFile,
	},
}

type appConfiguration struct {
	BasicAuth httpBasicAuthConfig `toml:"basic-auth"`
}

type Apps map[string]appConfiguration

type configuration struct {
	Apps `toml:"apps"`
}

func loadConfigFromTOMLFile(fs afero.Afero, filepath string) (*configuration, error) {
	if isDir, err := fs.IsDir(filepath); isDir {
		return nil, fmt.Errorf("%s is a directory", filepath)
	} else if err != nil {
		return nil, err
	}

	data, err := fs.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var loadedConfig configuration
	err = toml.Unmarshal(data, &loadedConfig)
	return &loadedConfig, err
}

func noAppInConfig(config *configuration) bool {
	return len(config.Apps) == 0
}

func appConfigIsEmpty(appConfig appConfiguration) bool {
	return appConfig == appConfiguration{}
}

func getNameOfAppsWithMissingAuth(appConfigs map[string]appConfiguration) []string {
	appNamesWithMissingAuth := make([]string, 0)
	for appName, appConfig := range appConfigs {
		if appConfigIsEmpty(appConfig) {
			appNamesWithMissingAuth = append(appNamesWithMissingAuth, appName)
		}
	}

	sort.Strings(appNamesWithMissingAuth)
	return appNamesWithMissingAuth
}

func validateConfig(config *configuration) error {
	if noAppInConfig(config) {
		return errors.New("no app in configuration file")
	}

	appNamesWithMissingAuth := getNameOfAppsWithMissingAuth(config.Apps)
	if len(appNamesWithMissingAuth) > 0 {
		joinedAppNamesWithMissingAuth := strings.Join(appNamesWithMissingAuth, ", ")
		return fmt.Errorf("%s has no authentication provided", joinedAppNamesWithMissingAuth)
	}

	return nil
}

func getConfigFilePathAndItsLoader(fs afero.Afero, configDir string) (string, func(afero.Afero, string) (*configuration, error), error) {
	configDir = addTrailingSlashIfNotExistsOnGivenPath(configDir)

	for _, configFileFormat := range configFileFormats {
		filepath := configDir + configFileFormat.filename

		fileExists, err := fs.Exists(filepath)
		if err != nil {
			return "", nil, err
		} else if fileExists {
			return filepath, configFileFormat.load, nil
		}
	}

	return "", nil, errors.New("no configuration file was found")
}

func loadConfigFileFromDir(fs afero.Afero, configDir string) (*configuration, error) {
	filepath, load, err := getConfigFilePathAndItsLoader(fs, configDir)
	if err != nil {
		return nil, err
	}

	loadedConfig, err := load(fs, filepath)
	if err != nil {
		return nil, err
	}

	err = validateConfig(loadedConfig)
	if err != nil {
		return nil, err
	}

	return loadedConfig, nil
}
