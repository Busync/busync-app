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

var ConfigFileFormats = map[string]struct {
	filename string
	load     func(afero.Afero, string, interface{}) error
}{
	"toml": {
		filename: TOML_CONFIG_FILE,
		load:     LoadTOMLFile,
	},
}

type AppConfig struct {
	BasicAuth HTTPBasicAuthConfig `toml:"basic-auth"`
}

type Apps map[string]AppConfig

type Config struct {
	Apps `toml:"apps"`
}

func LoadTOMLFile(fs afero.Afero, filepath string, v interface{}) error {
	if isDir, err := fs.IsDir(filepath); isDir {
		return fmt.Errorf("%s is a directory", filepath)
	} else if err != nil {
		return err
	}

	data, err := fs.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &v)
	return err
}

	return len(config.Apps) == 0
}

func AppConfigIsEmpty(appConfig AppConfig) bool {
	return appConfig == AppConfig{}
}

func GetNameOfAppsWithMissingAuth(appConfigs map[string]AppConfig) []string {
	appNamesWithMissingAuth := make([]string, 0)
	for appName, appConfig := range appConfigs {
		if AppConfigIsEmpty(appConfig) {
			appNamesWithMissingAuth = append(appNamesWithMissingAuth, appName)
		}
	}

	sort.Strings(appNamesWithMissingAuth)
	return appNamesWithMissingAuth
}

func ValidateConfig(config Config) error {
	if NoAppInConfig(config) {
		return errors.New("no app in configuration file")
	}

	appNamesWithMissingAuth := GetNameOfAppsWithMissingAuth(config.apps)
	if len(appNamesWithMissingAuth) > 0 {
		joinedAppNamesWithMissingAuth := strings.Join(appNamesWithMissingAuth, ", ")
		return fmt.Errorf("%s has no authentication provided", joinedAppNamesWithMissingAuth)
	}

	return nil
}

func GetConfigFilePathAndItsLoader(fs afero.Afero, configDir string) (string, func(afero.Afero, string, interface{}) error, error) {
	for _, configFileFormat := range ConfigFileFormats {
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

func LoadConfigFileFromDir(fs afero.Afero, configDir string) (*Config, error) {
	filepath, load, err := GetConfigFilePathAndItsLoader(fs, configDir)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = load(fs, filepath, config)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(config)
	return &config, err
}
