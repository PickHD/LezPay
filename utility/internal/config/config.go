package config

import "github.com/PickHD/LezPay/utility/internal/helper"

type (
	Configuration struct {
		Server *Server
		Redis  *Redis
		Tracer *Tracer
		Mailer *Mailer
	}

	Server struct {
		AppPort int
		AppEnv  string
		AppName string
		AppID   string
	}

	Redis struct {
		Host string
		Port int
		TTL  int
	}

	Tracer struct {
		JaegerURL string
	}

	Mailer struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
		IsTLS    bool
		SSL      int
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Server: &Server{
			AppPort: helper.GetEnvInt("APP_PORT"),
			AppEnv:  helper.GetEnvString("APP_ENV"),
			AppName: helper.GetEnvString("APP_NAME"),
			AppID:   helper.GetEnvString("APP_ID"),
		},
		Redis: &Redis{
			Host: helper.GetEnvString("REDIS_HOST"),
			Port: helper.GetEnvInt("REDIS_PORT"),
			TTL:  helper.GetEnvInt("REDIS_TTL"),
		},
		Tracer: &Tracer{
			JaegerURL: helper.GetEnvString("JAEGER_URL"),
		},
		Mailer: &Mailer{
			Host:     helper.GetEnvString("SMTP_HOST"),
			Port:     helper.GetEnvInt("SMTP_PORT"),
			Username: helper.GetEnvString("SMTP_USERNAME"),
			Password: helper.GetEnvString("SMTP_PASSWORD"),
			Sender:   helper.GetEnvString("SMTP_SENDER"),
			SSL:      helper.GetEnvInt("SMTP_SSL"),
			IsTLS:    helper.GetEnvBool("SMTP_IS_TLS"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
