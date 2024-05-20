package config

import "github.com/PickHD/LezPay/auth/internal/helper"

type (
	Configuration struct {
		Server   *Server
		Common   *Common
		Database *Database
		Redis    *Redis
		Secret   *Secret
		Tracer   *Tracer
		Service  *Service
	}

	Common struct {
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

	Redis struct {
		Host string
		Port int
		TTL  int
	}

	Secret struct {
		JWTSecret string
	}

	Tracer struct {
		JaegerURL string
	}

	Service struct {
		GRPCCustomerHost string
		GRPCMerchantHost string
		GRPCWalletHost   string
		UtilityURL       string
	}
)

func loadConfiguration() *Configuration {
	return &Configuration{
		Common: &Common{
			JWTExpire: helper.GetEnvInt("JWT_EXPIRE"),
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
		Service: &Service{
			GRPCCustomerHost: helper.GetEnvString("GRPC_CUSTOMER_HOST"),
			GRPCMerchantHost: helper.GetEnvString("GRPC_MERCHANT_HOST"),
			GRPCWalletHost:   helper.GetEnvString("GRPC_WALLET_HOST"),
			UtilityURL:       helper.GetEnvString("UTILITY_URL"),
		},
	}
}

func NewConfig() *Configuration {
	return loadConfiguration()
}
