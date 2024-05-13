package config

import "github.com/PickHD/LezPay/payment/internal/helper"

type (
	Configuration struct {
		Server   *Server
		Database *Database
		Redis    *Redis
		Tracer   *Tracer
		Kafka    *Kafka
	}

	Server struct {
		AppPort int
		AppEnv  string
		AppName string
		AppID   string
	}

	Database struct {
		Port     int
		Host     string
		Username string
		Password string
		Name     string
	}

	Redis struct {
		Host string
		Port int
		TTL  int
	}

	Tracer struct {
		JaegerURL string
	}

	Kafka struct {
		FirstBrokerHost     string
		SecondBrokerHost    string
		TopicRequestPayment string
		RequiredAcks        int
		GroupID             string
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
		Database: &Database{
			Port:     helper.GetEnvInt("DB_PORT"),
			Host:     helper.GetEnvString("DB_HOST"),
			Username: helper.GetEnvString("DB_USERNAME"),
			Password: helper.GetEnvString("DB_PASSWORD"),
			Name:     helper.GetEnvString("DB_NAME"),
		},
		Redis: &Redis{
			Host: helper.GetEnvString("REDIS_HOST"),
			Port: helper.GetEnvInt("REDIS_PORT"),
			TTL:  helper.GetEnvInt("REDIS_TTL"),
		},
		Tracer: &Tracer{
			JaegerURL: helper.GetEnvString("JAEGER_URL"),
		},
		Kafka: &Kafka{
			FirstBrokerHost:     helper.GetEnvString("KAFKA_FIRST_BROKER_HOST"),
			SecondBrokerHost:    helper.GetEnvString("KAFKA_SECOND_BROKER_HOST"),
			TopicRequestPayment: helper.GetEnvString("KAFKA_TOPIC_REQUEST_PAYMENT"),
			RequiredAcks:        helper.GetEnvInt("KAFKA_REQUIRED_ACKS"),
			GroupID:             helper.GetEnvString("KAFKA_GROUP_ID"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}