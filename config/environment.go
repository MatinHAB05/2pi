package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Environment struct {
	BotToken         BotToken         `mapstructure:",squash"`
	DebugModeOptions DebugModeOptions `mapstructure:",squash" json:"debug_mode_options"`
	DataBase         DataBase         `mapstructure:",squash" json:"database"`
	Redis            Redis            `mapstructure:",squash" json:"redis"`
	Logger           LoggerConfig     `mapstructure:",squash" json:"logger"`
	Casbin           Casbin           `mapstructure:",squash" json:"casbin"`
}

type BotToken struct {
	Token string `mapstructure:"BOT_TOKEN" json:"bot_token"`
}

type DebugModeOptions struct {
	Flag bool `mapstructure:"DEBUG_FLAG" json:"flag"`
}

type DataBase struct {
	Host     string `mapstructure:"PDB_HOST"     json:"host"`
	Port     int    `mapstructure:"PDB_PORT"     json:"port"`
	User     string `mapstructure:"PDB_USER"     json:"user"`
	Password string `mapstructure:"PDB_PASSWORD" json:"password"`
	Name     string `mapstructure:"PDB_DB_NAME"  json:"name"`
	SSLMode  string `mapstructure:"PDB_SSL_MODE" json:"ssl_mode"`
}

type Redis struct {
	Host      string `mapstructure:"RDB_HOST"      json:"host"`
	Port      int    `mapstructure:"RDB_PORT"      json:"port"`
	Password  string `mapstructure:"RDB_PASSWORD"  json:"password"`
	RDBNumber int    `mapstructure:"RDB_DB_NUMBER" json:"db"`
}

type LoggerConfig struct {
	FilePath string `mapstructure:"LOGGER_FILE_PATH" json:"file_path"`
	Encoding string `mapstructure:"LOGGER_ENCODING"  json:"encoding"`
	Level    string `mapstructure:"LOGGER_LEVEL"     json:"level"`
	Logger   string `mapstructure:"LOGGER_NAME"      json:"logger"`
}

type Casbin struct {
	ModelConfigFilePath string `mapstructure:"CASBIN_MODEL_PATH" json:"model_path"`
}

func NewEnvironment() *Environment {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("../..")
	v.SetConfigType("env")

	v.SetConfigName(".env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var env Environment
	if err := v.Unmarshal(&env); err != nil {
		panic(fmt.Errorf("failed to unmarshal environment config: %w", err))
	}

	return &env
}

func (c *Environment) String() string {
	if c == nil {
		return "<nil>"
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Config{BotToken:%+v, Debug:%+v, DataBase:%+v, Redis:%+v}",
			c.BotToken, c.DebugModeOptions, c.DataBase, c.Redis)
	}

	return "Environment: " + string(data)
}
