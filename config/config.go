package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	AppPort      string `json:"app_port"`
	AppEnv       string `json:"app_env"`
	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type MysqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Config struct {
	App   App
	Mysql MysqlDB 
}

func Init() {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},
		Mysql: MysqlDB{
			Host:      viper.GetString("DB_HOST"), 
			Port:      viper.GetString("DB_PORT"),
			User:      viper.GetString("DB_USER"),
			Password:  viper.GetString("DB_PASSWORD"),
			DBName:    viper.GetString("DB_NAME"),
			DBMaxOpen: viper.GetInt("DB_MAX_OPEN_CONNECTION"),
			DBMaxIdle: viper.GetInt("DB_MAX_IDLE_CONNECTION"),
		},
	}
}