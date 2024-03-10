package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	PostalCodeHosts []PostalCodeHost `mapstructure:"PostalCodeHosts"`
}

type PostalCodeHost struct {
	Name string `mapstructure:"Name"`
	Host string `mapstructure:"Host"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
