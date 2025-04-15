package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr           string
	ProductionMode       bool   `mapstructure:"PROD"`
	PostgresUrl          string `mapstructure:"POSTGRES_URL"`
	SecretKey            string `mapstructure:"SECRET_KEY"`
	AdminName            string `mapstructure:"ADMIN_NAME"`
	MinIOEndpoint        string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKeyID     string `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MinIOSecretAccessKey string `mapstructure:"MINIO_SECRET_ACCESS_KEY"`
	MinIOBucketName      string `mapstructure:"MINIO_BUCKET"`
	MinIOPublicHost      string `mapstructure:"MINIO_PUBLIC_HOST"`
	YaClientId           string `mapstructure:"YANDEX_API_CLIENT_ID"`
	YaSecret             string `mapstructure:"YANDEX_API_SECRET"`
	YaRedirectUrl        string `mapstructure:"YANDEX_API_REDIRECT_URL"`
}

func (c *Config) InitProd() {
}

func (c *Config) InitDev() {
	var (
		name     string = viper.GetString("POSTGRES_USER")
		password string = viper.GetString("POSTGRES_PASSWORD")
		host     string = "localhost"
		port     string = viper.GetString("POSTGRES_PORT")
		db       string = viper.GetString("POSTGRES_DB")
	)
	c.PostgresUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", name, password, host, port, db)
}

func New() *Config {
	var config Config = Config{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err.Error())
		return nil
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Config error: %s", err.Error())
		return nil
	}
	config.ServerAddr = fmt.Sprintf("%s:%s", viper.Get("SERVER_HOST"), viper.Get("SERVER_PORT"))

	if config.ProductionMode {
		config.InitProd()
	} else {
		config.InitDev()
	}
	return &config
}
