package busylight_sync

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
)

const (
	TOML_CONFIG_FILE string = ".busylight-sync.toml"
)

type Config struct{}

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
