package config

import "github.com/PickHD/LezPay/customer/internal/v1/helper"

type (
	Configuration struct {
		Common   *Common
		Server   *Server
		Database *Database
		Secret   *Secret
		Redis    *Redis
		Tracer   *Tracer
		Kafka    *Kafka
		Service  *Service
	}

	Common struct {
		GrpcPort  int
		JWTExpire int
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

	Secret struct {
		JWTSecret string
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
		FirstBrokerHost        string
		SecondBrokerHost       string
		TopicTopupTransaction  string
		TopicPayoutTransaction string
		RequiredAcks           int
		GroupID                string
	}

	Service struct {
		GRPCWalletHost string
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Common: &Common{
			JWTExpire: helper.GetEnvInt("JWT_EXPIRE"),
			GrpcPort:  helper.GetEnvInt("GRPC_PORT"),
		},
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
		Secret: &Secret{
			JWTSecret: helper.GetEnvString("JWT_SECRET"),
		},
		Tracer: &Tracer{
			JaegerURL: helper.GetEnvString("JAEGER_URL"),
		},
		Kafka: &Kafka{
			FirstBrokerHost:        helper.GetEnvString("KAFKA_FIRST_BROKER_HOST"),
			SecondBrokerHost:       helper.GetEnvString("KAFKA_SECOND_BROKER_HOST"),
			TopicTopupTransaction:  helper.GetEnvString("KAFKA_TOPIC_TOPUP_TRANSACTION"),
			TopicPayoutTransaction: helper.GetEnvString("KAFKA_TOPIC_PAYOUT_TRANSACTION"),
			RequiredAcks:           helper.GetEnvInt("KAFKA_REQUIRED_ACKS"),
			GroupID:                helper.GetEnvString("KAFKA_GROUP_ID"),
		},
		Service: &Service{
			GRPCWalletHost: helper.GetEnvString("GRPC_WALLET_HOST"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
