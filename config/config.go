package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Chains []ChainConfig
}

type ChainConfig struct {
	RegistryName string `mapstructure:"registry-name"`
	Denom        string
	Exponent     uint
	RestAPI      string `mapstructure:"rest-api"`

	Validator string
	Wallets   []Wallet
}

type Wallet struct {
	Name    string
	Address string
}

func LoadConfig(file string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
