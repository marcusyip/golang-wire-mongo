package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Mongodb  MongodbConfig  `json:"mongodb"`
	Facebook FacebookConfig `json:"facebook"`
}

type MongodbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type FacebookConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func ProvideConfig() *Config {
	viper.SetEnvPrefix("fh")
	viper.BindEnv("mongo.host")
	viper.BindEnv("mongo.port")

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("/etc/sample/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	conf := &Config{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Could not unmarshal config: %s", err))
	}
	return conf
}
