package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Environment struct {
	Admin               AdminConfig         `mapstructure:",squash"`
	BotToken            BotToken            `mapstructure:",squash"`
	ModeOptions         ModeOptions         `mapstructure:",squash" json:"mode_options"`
	DataBase            DataBase            `mapstructure:",squash" json:"database"`
	Redis               Redis               `mapstructure:",squash" json:"redis"`
	Logger              LoggerConfig        `mapstructure:",squash" json:"logger"`
	Casbin              Casbin              `mapstructure:",squash" json:"casbin"`
	Email               EmailConfig         `mapstructure:",squash" json:"email"`
	ElasticSearchConfig ElasticSearchConfig `mapstructure:",squash" json:"elastic"`
}
type AdminConfig struct {
	UserID    int64  `mapstructure:"ADMIN_USER_ID" json:"admin_user_id"`
	Username  string `mapstructure:"ADMIN_USER_USERNAME" json:"admin_username"`
	FirstName string `mapstructure:"ADMIN_USER_FIRSTNAME" json:"admin_first_name"`
	LastName  string `mapstructure:"ADMIN_USER_LASTNAME" json:"admin_last_name"`
	Language  string `mapstructure:"ADMIN_USER_LANG" json:"admin_language"`
}
type BotToken struct {
	Token string `mapstructure:"BOT_TOKEN" json:"bot_token"`
}

type ModeOptions struct {
	Debug                    bool `mapstructure:"DEBUG_FLAG" json:"debug_flag"`
	ProductionFlag           bool `mapstructure:"PRODUCTION_FLAG" json:"production_flag"`
	DebugRBACMiddelware      bool `mapstructure:"DEBUG_RBAC" json:"debug_rbac_middle"`
	SaveJsonScrapperArticles bool `mapstructure:"SAVE_JSON_SCRAPPER_ARTICLES" json:"save_json_scrapper_articles"`
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

type EmailConfig struct {
	TwoPiEmail            string `mapstructure:"TWO_PI_EMAIL"              json:"2pi_email"`
	TwoPiEmailAppPassword string `mapstructure:"TWO_PI_EMAIL_APP_PASSWORD" json:"2pi_email_app_password"`
	SMTPHost              string `mapstructure:"SMTP_HOST"                        json:"smtp_host"`
	SMTPPort              string `mapstructure:"SMTP_HOST_PORT"                   json:"smtp_host_port"`
}

type ElasticSearchConfig struct {
	Host     string `mapstructure:"ES_HOST"      json:"host"`
	Port     int    `mapstructure:"ES_PORT"      json:"port"`
	Password string `mapstructure:"ES_PASSWORD"  json:"password"`
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
			c.BotToken, c.ModeOptions, c.DataBase, c.Redis)
	}

	return "Environment: " + string(data)
}
