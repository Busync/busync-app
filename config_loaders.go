package busylight_sync

import (
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
)

const (
	TOML_CONFIG_FILE string = ".busylight-sync.toml"
)

type Config struct{}

func LoadTOMLFile(fs afero.Afero, filepath string, v interface{}) error {
	data, err := fs.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &v)
	return err
}
