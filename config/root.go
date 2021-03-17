package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	logger *logrus.Logger
	viper  *viper.Viper
}

func (q Config) Refresh() {
	q.viper.SetDefault("ENV", "production")
	q.viper.SetDefault("NO_SSL_PORT", "80")
	q.viper.SetDefault("SSL_PORT", "443")
	q.viper.SetDefault("SSL_CERT", "")
	q.viper.SetDefault("SSL_PEM", "")
	q.viper.SetDefault("LOG_PATH", "./logs")
	q.viper.SetDefault("USER", "admin")
	q.viper.SetDefault("PASS", "admin")
	q.viper.SetDefault("PUT_ALLOWED_IPS", []string{}) //[]string{"127.0.0.1/24", "10.0.0.1/24"}
	q.viper.SetDefault("GET_ALLOWED_IPS", []string{})
	q.viper.SetDefault("IDENTITY_KEY", "user_id")

	q.viper.SetDefault("IMAGE_TASK_QUEUE", "orc.image-task")
	q.viper.SetDefault("IMAGE_SUCCESS_QUEUE", "orc.success")
	q.viper.SetDefault("IMAGE_ERROR_QUEUE", "orc.error")

	q.viper.SetDefault("MONGODB_HOST", "localhost")
	q.viper.SetDefault("MONGODB_PORT", "27017")
	q.viper.SetDefault("MONGODB_USERNAME", "")
	q.viper.SetDefault("MONGODB_PASSWORD", "")
	q.viper.SetDefault("MONGODB_DB", "ocr")

	q.viper.SetDefault("RABBITMQ_HOST", "localhost")
	q.viper.SetDefault("RABBITMQ_PORT", "5672")
	q.viper.SetDefault("RABBITMQ_USERNAME", "guest")
	q.viper.SetDefault("RABBITMQ_PASSWORD", "guest")
	q.viper.SetDefault("RABBITMQ_VHOST", "/")

	q.viper.AutomaticEnv() // Read from os env
	q.ReadFromEnvFile()
}

var viperENV = viper.New()
var hasEnv = true

func (q *Config) ReadFromEnvFile() {
	if hasEnv {
		viperENV.AddConfigPath(".")
		viperENV.AddConfigPath("../") //fix test app .env wrong path
		viperENV.SetConfigName(".env")
		viperENV.SetConfigType("env")
		viperENV.AutomaticEnv()
		err := viperENV.ReadInConfig()
		if err != nil {
			hasEnv = false
			if q.logger != nil {
				q.logger.WithError(err).Warn("load env failed")
			}
		}
		for _, k := range viperENV.AllKeys() {
			switch k {
			case "PUT_ALLOWED_IPS", "GET_ALLOWED_IPS":
				q.viper.Set(k, strings.Split(viperENV.Get(k).(string), ","))
			default:
				q.viper.Set(k, viperENV.Get(k))
			}
		}
	}
}

func NewConfig(logger *logrus.Logger) (*viper.Viper, error) {
	var CONFIG = Config{}
	CONFIG.viper = viper.New()
	if logger != nil {
		CONFIG.logger = logger
	}
	CONFIG.Refresh()
	return CONFIG.viper, nil
}
