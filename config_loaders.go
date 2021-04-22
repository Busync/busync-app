package busylight_sync

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
	basicAuth HTTPBasicAuthConfig
}

type Config struct {
	apps map[string]AppConfig
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

func NoAppInConfig(config Config) bool {
	return len(config.apps) == 0
}

func AppConfigIsEmpty(appConfig AppConfig) bool {
	return appConfig == AppConfig{}
}

func GetNamesOfEmptyAppConfigs(appConfigs map[string]AppConfig) []string {
	emptyAppConfigs := make([]string, 0)
	for appName, appConfig := range appConfigs {
		if AppConfigIsEmpty(appConfig) {
			emptyAppConfigs = append(emptyAppConfigs, appName)
		}
	}

	sort.Strings(emptyAppConfigs)
	return emptyAppConfigs
}

func ValidateConfig(config Config) error {
	if NoAppInConfig(config) {
		return errors.New("no app in configuration file")
	}

	emptyAppNames := GetNamesOfEmptyAppConfigs(config.apps)
	if len(emptyAppNames) == 1 {
		return fmt.Errorf("%s configuration is empty", emptyAppNames[0])
	} else if len(emptyAppNames) > 1 {
		joinedEmptyAppNames := strings.Join(emptyAppNames, ", ")
		return fmt.Errorf("%s configurations are empty", joinedEmptyAppNames)
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
