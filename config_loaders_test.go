package busylight_sync

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func NewFileSystem() afero.Afero {
	return afero.Afero{Fs: afero.NewMemMapFs()}
}

func TestLoaders(t *testing.T) {
	var testCases = []struct {
		desc     string
		filepath string
		loader   func(afero.Afero, string, interface{}) error
	}{
		{
			desc:     "toml file",
			filepath: "/" + TOML_CONFIG_FILE,
			loader:   LoadTOMLFile,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc+"/file_not_found", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()
			originalconfig := Config{}

			configPassedToLoader := originalconfig
			err := LoadTOMLFile(fs, tC.filepath, configPassedToLoader)

			assert.EqualError(err, fmt.Sprintf("open %s: file does not exist", tC.filepath))
			assert.Equal(originalconfig, configPassedToLoader)
		})
	}
}
