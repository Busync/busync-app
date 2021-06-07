package main

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func newFileSystem() afero.Afero {
	return afero.Afero{Fs: afero.NewMemMapFs()}
}

func TestConfigFileFormats(t *testing.T) {
	for configFileFormat, tC := range configFileFormats {
		filepath := "/" + tC.filename

		t.Run(configFileFormat+"/file_not_found", func(t *testing.T) {
			assert := assert.New(t)
			fs := newFileSystem()

			sut, err := loadConfigFromTOMLFile(fs, filepath)

			assert.EqualError(err, fmt.Sprintf("open %s: file does not exist", filepath))
			assert.Nil(sut)
		})

		t.Run(configFileFormat+"/is_a_dir", func(t *testing.T) {
			assert := assert.New(t)
			fs := newFileSystem()
			err := fs.Mkdir(filepath, 0777)
			if err != nil {
				panic(err)
			}

			sut, err := loadConfigFromTOMLFile(fs, filepath)

			assert.EqualError(err, fmt.Sprintf("%s is a directory", filepath))
			assert.Nil(sut)
		})
	}
}

func TestGetConfigFilePathAndItsLoader(t *testing.T) {
	testCases := []struct {
		desc         string
		configDir    string
		filepath     string
		fileContent  string
		wantFilepath string
		wantLoad     func(afero.Afero, string) (*configuration, error)
		wantErr      string
	}{
		{
			desc:      "no configuration file found",
			configDir: "/",
			wantErr:   "no configuration file was found",
		},
		{
			desc:         "toml configuration on rootdir",
			configDir:    "/",
			wantFilepath: "/" + TOML_CONFIG_FILE,
			wantLoad:     loadConfigFromTOMLFile,
		},
		{
			desc:         "configuration on subdir with trailing slash",
			configDir:    "/subdir/",
			wantFilepath: "/subdir/" + TOML_CONFIG_FILE,
			wantLoad:     loadConfigFromTOMLFile,
		},
		{
			desc:         "configuration on subdir without trailing slash",
			configDir:    "/subdir",
			wantFilepath: "/subdir/" + TOML_CONFIG_FILE,
			wantLoad:     loadConfigFromTOMLFile,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			fs := newFileSystem()
			wantLoadFuncName := getFuncName(tC.wantLoad)

			if tC.wantFilepath != "" {
				fs.WriteFile(tC.wantFilepath, []byte(tC.fileContent), 0755)
			}

			gotFilepath, gotLoad, err := getConfigFilePathAndItsLoader(fs, tC.configDir)
			gotLoadFuncName := getFuncName(gotLoad)

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
		config *configuration
	}{
		{
			desc:   "no app",
			config: &configuration{},
		},
		{
			desc: "one app",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc: "two apps",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"bar": {
						BasicAuth: httpBasicAuthConfig{
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

			got := noAppInConfig(tC.config)

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
		wantAppConfig appConfiguration
	}{
		{
			desc:      "is empty",
			wantAppConfig: appConfiguration{},
		},
		{
			desc: "is not empty",
			wantAppConfig: appConfiguration{
				BasicAuth: httpBasicAuthConfig{
					Username: "foobar",
					Password: "spameggs",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)

			got := appConfigIsEmpty(tC.wantAppConfig)

			if got {
				assert.Equal(appConfiguration{}, tC.wantAppConfig)
			} else {
				assert.NotEqual(appConfiguration{}, tC.wantAppConfig)
			}
		})
	}
}

func TestGetNameOfAppsWithMissingAuth(t *testing.T) {
	testCases := []struct {
		desc       string
		appConfigs map[string]appConfiguration
		want       []string
	}{
		{
			desc:       "no apps",
			appConfigs: make(map[string]appConfiguration),
			want:       []string{},
		},
		{
			desc: "one app config with auth",
			appConfigs: map[string]appConfiguration{
				"foo": {
					BasicAuth: httpBasicAuthConfig{
						Username: "foobar",
						Password: "spameggs",
					},
				},
			},
			want: []string{},
		},
		{
			desc: "one app config with missing auth",
			appConfigs: map[string]appConfiguration{
				"foo": {},
			},
			want: []string{"foo"},
		},
		{
			desc: "two app configs with auth",
			appConfigs: map[string]appConfiguration{
				"foo": {
					BasicAuth: httpBasicAuthConfig{
						Username: "foobar",
						Password: "spameggs",
					},
				},
				"bar": {
					BasicAuth: httpBasicAuthConfig{
						Username: "barbaz",
						Password: "hamspam",
					},
				},
			},
			want: []string{},
		},
		{
			desc: "two app configs with missing auth",
			appConfigs: map[string]appConfiguration{
				"foo": {},
				"bar": {},
			},
			want: []string{"bar", "foo"},
		},
		{
			desc: "two app configs with one missing auth",
			appConfigs: map[string]appConfiguration{
				"foo": {
					BasicAuth: httpBasicAuthConfig{
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

			got := getNameOfAppsWithMissingAuth(tC.appConfigs)

			assert.Equal(tC.want, got)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		config  *configuration
		wantErr string
	}{
		{
			desc:    "no app config",
			config:  &configuration{},
			wantErr: "no app in configuration file",
		},
		{
			desc: "one app config with auth",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
				},
			},
		},
		{
			desc: "one app config with missing auth",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {},
				},
			},
			wantErr: "foo has no authentication provided",
		},
		{
			desc: "two app configs with auth",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
							Username: "foobar",
							Password: "spameggs",
						},
					},
					"bar": {
						BasicAuth: httpBasicAuthConfig{
							Username: "barbaz",
							Password: "hamspam",
						},
					},
				},
			},
		},
		{
			desc: "two app configs with missing auth",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {},
					"bar": {},
				},
			},
			wantErr: "bar, foo has no authentication provided",
		},
		{
			desc: "two app configs with one missing auth",
			config: &configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
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

			err := validateConfig(tC.config)

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
		want        configuration
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

			want: configuration{
				Apps: map[string]appConfiguration{
					"foo": {
						BasicAuth: httpBasicAuthConfig{
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
			fs := newFileSystem()
			if tC.filename != "" {
				filepath := configDir + tC.filename
				fs.WriteFile(filepath, []byte(tC.fileContent), 0755)
			}

			if tC.desc == "one app config with auth" {
				fmt.Println("test")
			}
			got, err := loadConfigFileFromDir(fs, configDir)
			if err != nil {
				assert.EqualError(err, tC.wantErr)
				assert.Nil(got)
			}
		})
	}
}
