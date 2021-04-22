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

func TestConfigFileFormats(t *testing.T) {
	for configFileFormat, tC := range ConfigFileFormats {
		filepath := "/" + tC.filename
		originalconfig := Config{}

		t.Run(configFileFormat+"/file_not_found", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()

			configPassedToLoader := originalconfig
			err := LoadTOMLFile(fs, filepath, configPassedToLoader)

			assert.EqualError(err, fmt.Sprintf("open %s: file does not exist", filepath))
			assert.Equal(originalconfig, configPassedToLoader)
		})

		t.Run(configFileFormat+"/is_a_dir", func(t *testing.T) {
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

func TestGetConfigFilePathAndItsLoader(t *testing.T) {
	configDir := "/"

	testCases := []struct {
		desc         string
		filepath     string
		fileContent  string
		wantFilepath string
		wantLoad     func(afero.Afero, string, interface{}) error
		wantErr      string
	}{
		{
			desc:    "no configuration file found",
			wantErr: "no configuration file was found",
		},
		{
			desc:         "toml configuration found",
			wantFilepath: configDir + TOML_CONFIG_FILE,
			wantLoad:     LoadTOMLFile,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()
			wantLoadFuncName := GetFuncName(tC.wantLoad)

			if tC.wantFilepath != "" {
				fs.WriteFile(tC.wantFilepath, []byte(tC.fileContent), 0755)
			}

			gotFilepath, gotLoad, err := GetConfigFilePathAndItsLoader(fs, configDir)
			gotLoadFuncName := GetFuncName(gotLoad)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
			}

			assert.Equal(tC.wantFilepath, gotFilepath)
			assert.Equal(wantLoadFuncName, gotLoadFuncName)
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

func TestGetNamesOfEmptyAppConfigs(t *testing.T) {
	testCases := []struct {
		desc       string
		appConfigs map[string]AppConfig
		want       []string
	}{
		{
			desc:       "no app config",
			appConfigs: make(map[string]AppConfig),
			want:       []string{},
		},
		{
			desc: "one non empty app config",
			appConfigs: map[string]AppConfig{
				"foo": AppConfig{
					basicAuth: HTTPBasicAuthConfig{
						username: "foobar",
						password: "spameggs",
					},
				},
			},
			want: []string{},
		},
		{
			desc: "one empty app config",
			appConfigs: map[string]AppConfig{
				"foo": AppConfig{},
			},
			want: []string{"foo"},
		},
		{
			desc: "two non empty app configs",
			appConfigs: map[string]AppConfig{
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
			want: []string{},
		},
		{
			desc: "two empty app configs",
			appConfigs: map[string]AppConfig{
				"foo": AppConfig{},
				"bar": AppConfig{},
			},
			want: []string{"bar", "foo"},
		},
		{
			desc: "two app configs with one empty",
			appConfigs: map[string]AppConfig{
				"foo": AppConfig{
					basicAuth: HTTPBasicAuthConfig{
						username: "foobar",
						password: "spameggs",
					},
				},
				"bar": AppConfig{},
			},
			want: []string{"bar"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := GetNamesOfEmptyAppConfigs(tC.appConfigs)

			assert.Equal(tC.want, got)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		config  Config
		wantErr string
	}{
		{
			desc:    "no app config",
			config:  Config{},
			wantErr: "no app in configuration file",
		},
		{
			desc: "one non empty app config",
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
			desc: "one empty app configs",
			config: Config{
				apps: map[string]AppConfig{
					"foo": AppConfig{},
				},
			},
			wantErr: "foo configuration is empty",
		},
		{
			desc: "two non empty app configs",
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
		{
			desc: "two empty app configs",
			config: Config{
				apps: map[string]AppConfig{
					"foo": AppConfig{},
					"bar": AppConfig{},
				},
			},
			wantErr: "bar, foo configurations are empty",
		},
		{
			desc: "two app configs with one empty",
			config: Config{
				apps: map[string]AppConfig{
					"foo": AppConfig{
						basicAuth: HTTPBasicAuthConfig{
							username: "foobar",
							password: "spameggs",
						},
					},
					"bar": AppConfig{},
				},
			},
			wantErr: "bar configuration is empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			err := ValidateConfig(tC.config)

			if err != nil {
				assert.EqualError(err, tC.wantErr)
			}
		})
	}
}

func TestLoadConfigFileFromDir(t *testing.T) {
	assert := assert.New(t)
	fs := NewFileSystem()
	configDir := "/"

	got, err := LoadConfigFileFromDir(fs, configDir)

	assert.EqualError(err, "no configuration file was found")
	assert.Nil(got)
}
