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
	for loaderType, tC := range Loaders {
		filepath := "/" + tC.filename
		originalconfig := Config{}

		t.Run(loaderType+"/file_not_found", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()

			configPassedToLoader := originalconfig
			err := LoadTOMLFile(fs, filepath, configPassedToLoader)

			assert.EqualError(err, fmt.Sprintf("open %s: file does not exist", filepath))
			assert.Equal(originalconfig, configPassedToLoader)
		})

		t.Run(loaderType+"/is_a_dir", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()
			err := fs.Mkdir(filepath, 0777)
			if err != nil {
				panic(err)
			}

			configPassedToLoader := originalconfig
			err = LoadTOMLFile(fs, filepath, configPassedToLoader)

			assert.EqualError(err, fmt.Sprintf("%s is a directory", filepath))
			assert.Equal(originalconfig, configPassedToLoader)
		})
	}
}

func TestNoneOfConfigFileFound(t *testing.T) {
	assert := assert.New(t)
	fs := NewFileSystem()
	configDir := "/"
	want := Config{}

	got, err := LoadConfigFileFromDir(fs, configDir)

	assert.EqualError(err, "no configuration file was found")
	assert.Equal(want, got)
}
