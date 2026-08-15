package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Debug    Debug    `mapstructure:"debug" json:"debug"`
	BotToken BotToken `mapstructure:"bot-token" json:"bot_token"`
}

type BotToken struct {
	Token string `mapstructure:"token" json:"bot_token"`
}

type Debug struct {
	Flag bool `mapstructure:"flag" json:"flag"`
}

func LoadConfig() Config {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("env.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("ERR: Failed to read config file: ", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("ERR: Failed to unmarshal config: ", err)
	}

	fmt.Println(cfg.String())
	return cfg
}

func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Config{BotToken:%+v\nDebug:%+v}", c.BotToken, c.Debug)
	}

	return "Config: " + string(data)
}
