package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type Config struct {
	Logs LogsConfig `yaml:"logs"`
}

type LogsConfig struct {
	Encoding    string `yaml:"encoding"`
	Level       string `yaml:"level"`
	Development bool   `yaml:"development"`
}

//nolint:all
func (s Spec) readConfig(cliCtx *cli.Context) (cfg Config, err error) {
	viperCfg := viper.NewWithOptions()
	viperCfg.SetConfigName("config")
	viperCfg.SetConfigType("yaml")
	viperCfg.AddConfigPath(fmt.Sprintf("/etc/%s/", strings.ToLower(s.AppName)))
	viperCfg.AddConfigPath(fmt.Sprintf("$HOME/.%s/", strings.ToLower(s.AppName)))
	viperCfg.AddConfigPath(".")
	viperCfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperCfg.SetEnvPrefix(
		strings.ToUpper(strings.NewReplacer(".", "_", "-", "_").Replace(s.AppName)),
	)
	viperCfg.AutomaticEnv()

	if configPath := cliCtx.String(configFilePathFlag); configPath != "" {
		if _, err = os.Stat(configPath); err == nil {
			viperCfg.SetConfigFile(configPath)
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		viperErr := new(viper.ConfigFileNotFoundError)
		if !errors.As(err, viperErr) || !s.IgnoreMissingConfigFile {
			return
		} else {
		}
	}

	if err = viperCfg.Unmarshal(&cfg); err != nil {
		return
	}

	return
}
