package busylight_sync

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
)

const (
	TOML_CONFIG_FILE string = ".busylight-sync.toml"
)

var Loaders = map[string]struct {
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

func AppConfigIsEmpty(appConfig AppConfig) bool {
	return appConfig == AppConfig{}
}

func NoAppInConfig(config Config) bool {
	return len(config.apps) == 0
}

func LoadConfigFileFromDir(fs afero.Afero, configDir string) (Config, error) {
	for _, loader := range Loaders {
		config := Config{}
		filepath := configDir + loader.filename

		if err := loader.load(fs, filepath, config); err == nil {
			return config, nil
		}
	}

	return Config{}, fmt.Errorf("no configuration file was found")
}
