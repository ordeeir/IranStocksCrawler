package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigList map[string]string

type Config struct {
	//Databases    Databases    `json:"databases" yaml:"databases"`
	Cache    CacheConfig  `json:"cache" yaml:"cache"`
	Router   RouterConfig `json:"router" yaml:"router"`
	Settings ConfigList   `mapstructure:"general_settings" json:"general_settings" yaml:"general_settings"`
}

// NewConfig is a function to Load Configuration
func NewConfig(path string) (*Config, *viper.Viper, error) {

	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err.Error()))
	}
	v.WatchConfig()

	c := new(Config)
	err = v.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err.Error()))
	}

	return c, v, err
}
