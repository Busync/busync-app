package main

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

		t.Run(configFileFormat+"/file_not_found", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()

			config, err := LoadConfigFromTOMLFile(fs, filepath)

			assert.EqualError(err, fmt.Sprintf("open %s: file does not exist", filepath))
			assert.Nil(config)
		})

		t.Run(configFileFormat+"/is_a_dir", func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()
			err := fs.Mkdir(filepath, 0777)
			if err != nil {
				panic(err)
			}

			config, err := LoadConfigFromTOMLFile(fs, filepath)

			assert.EqualError(err, fmt.Sprintf("%s is a directory", filepath))
			assert.Nil(config)
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
		wantLoad     func(afero.Afero, string) (*Config, error)
		wantErr      string
	}{
		{
			desc:    "no configuration file found",
			wantErr: "no configuration file was found",
		},
		{
			desc:         "toml configuration found",
			wantFilepath: configDir + TOML_CONFIG_FILE,
			wantLoad:     LoadConfigFromTOMLFile,
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
		config *Config
	}{
		{
			desc:   "no app",
			config: &Config{},
		},
		{
			desc: "one app",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc: "two apps",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"bar": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "barbaz",
							Password: "hamspam",
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
				assert.Equal(0, len(tC.config.Apps))
			} else {
				assert.Less(0, len(tC.config.Apps))
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
				BasicAuth: HTTPBasicAuthConfig{
					Username: "foobar",
					Password: "spameggs",
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

func TestGetNameOfAppsWithMissingAuth(t *testing.T) {
	testCases := []struct {
		desc       string
		appConfigs map[string]AppConfig
		want       []string
	}{
		{
			desc:       "no apps",
			appConfigs: make(map[string]AppConfig),
			want:       []string{},
		},
		{
			desc: "one app config with auth",
			appConfigs: map[string]AppConfig{
				"foo": {
					BasicAuth: HTTPBasicAuthConfig{
						Username: "foobar",
						Password: "spameggs",
					},
				},
			},
			want: []string{},
		},
		{
			desc: "one app config with missing auth",
			appConfigs: map[string]AppConfig{
				"foo": {},
			},
			want: []string{"foo"},
		},
		{
			desc: "two app configs with auth",
			appConfigs: map[string]AppConfig{
				"foo": {
					BasicAuth: HTTPBasicAuthConfig{
						Username: "foobar",
						Password: "spameggs",
					},
				},
				"bar": {
					BasicAuth: HTTPBasicAuthConfig{
						Username: "barbaz",
						Password: "hamspam",
					},
				},
			},
			want: []string{},
		},
		{
			desc: "two app configs with missing auth",
			appConfigs: map[string]AppConfig{
				"foo": {},
				"bar": {},
			},
			want: []string{"bar", "foo"},
		},
		{
			desc: "two app configs with one missing auth",
			appConfigs: map[string]AppConfig{
				"foo": {
					BasicAuth: HTTPBasicAuthConfig{
						Username: "foobar",
						Password: "spameggs",
					},
				},
				"bar": {},
			},
			want: []string{"bar"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := GetNameOfAppsWithMissingAuth(tC.appConfigs)

			assert.Equal(tC.want, got)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		config  *Config
		wantErr string
	}{
		{
			desc:    "no app config",
			config:  &Config{},
			wantErr: "no app in configuration file",
		},
		{
			desc: "one app config with auth",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc: "one app config with missing auth",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {},
				},
			},
			wantErr: "foo has no authentication provided",
		},
		{
			desc: "two app configs with auth",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"bar": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "barbaz",
							Password: "hamspam",
						},
					},
				},
			},
		},
		{
			desc: "two app configs with missing auth",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {},
					"bar": {},
				},
			},
			wantErr: "bar, foo has no authentication provided",
		},
		{
			desc: "two app configs with one missing auth",
			config: &Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"bar": {},
				},
			},
			wantErr: "bar has no authentication provided",
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
	configDir := "/"

	testCases := []struct {
		desc        string
		filename    string
		fileContent string
		want        Config
		wantErr     string
	}{
		{
			desc:     "one app config with auth",
			filename: TOML_CONFIG_FILE,
			fileContent: `
[apps]
  [apps.foo]
    [apps.foo.basic-auth]
      password = "spameggs"
      username = "foobar"
`,

			want: Config{
				Apps: map[string]AppConfig{
					"foo": {
						BasicAuth: HTTPBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc:    "no config file",
			wantErr: "no configuration file was found",
		},
		{
			desc:     "no app in config",
			filename: TOML_CONFIG_FILE,
			wantErr:  "no app in configuration file",
		},
		{
			desc:     "one app config with missing auth",
			filename: TOML_CONFIG_FILE,
			fileContent: `
[apps]
  [apps.foo]
    [apps.foo.basic-auth]
`,
			wantErr: "foo has no authentication provided",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			fs := NewFileSystem()
			if tC.filename != "" {
				filepath := configDir + tC.filename
				fs.WriteFile(filepath, []byte(tC.fileContent), 0755)
			}

			if tC.desc == "one app config with auth" {
				fmt.Println("test")
			}
			got, err := LoadConfigFileFromDir(fs, configDir)
			if err != nil {
				assert.EqualError(err, tC.wantErr)
				assert.Nil(got)
			}
		})
	}
}
