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

func TestAppConfigIsEmpty(t *testing.T) {
	testCases := []struct {
		desc      string
		appConfig AppConfig
	}{
		{
			desc:      "is empty",
			appConfig: AppConfig{},
		},
		{
			desc: "is not empty",
			appConfig: AppConfig{
				basicAuth: HTTPBasicAuthConfig{
					username: "foobar",
					password: "spameggs",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := AppConfigIsEmpty(tC.appConfig)

			if got {
				assert.Equal(AppConfig{}, tC.appConfig)
			} else {
				assert.NotEqual(AppConfig{}, tC.appConfig)
			}
		})
	}
}

func TestNoAppInConfig(t *testing.T) {
	testCases := []struct {
		desc   string
		config Config
	}{
		{
			desc:   "no app",
			config: Config{},
		},
		{
			desc: "one app",
			config: Config{
				apps: map[string]AppConfig{
					"foo": AppConfig{
						basicAuth: HTTPBasicAuthConfig{
							username: "foobar",
							password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc: "two apps",
			config: Config{
				apps: map[string]AppConfig{
					"foo": AppConfig{
						basicAuth: HTTPBasicAuthConfig{
							username: "foobar",
							password: "spameggs",
						},
					},
					"bar": AppConfig{
						basicAuth: HTTPBasicAuthConfig{
							username: "barbaz",
							password: "hamspam",
						},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := NoAppInConfig(tC.config)

			if got {
				assert.Equal(0, len(tC.config.apps))
			} else {
				assert.Less(0, len(tC.config.apps))
			}
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
