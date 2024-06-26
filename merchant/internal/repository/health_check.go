package repository

import (
	"context"

	"github.com/PickHD/LezPay/merchant/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// HealthCheckRepository is an interface that has all the function to be implemented inside health check repository
	HealthCheckRepository interface {
		Check() (bool, error)
	}

	// HealthCheckRepositoryImpl is an app health check struct that consists of all the dependencies needed for health check repository
	HealthCheckRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *pgxpool.Pool
		Redis   *redis.Client
	}
)

// NewHealthCheckRepository return new instances health check repository
func NewHealthCheckRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client) *HealthCheckRepositoryImpl {
	return &HealthCheckRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
		Redis:   rds,
	}
}

func (hr *HealthCheckRepositoryImpl) Check() (bool, error) {
	tr := hr.Tracer.Tracer("Merchant-Check Repository")
	_, span := tr.Start(hr.Context, "Start Check")
	defer span.End()

	if err := hr.DB.Ping(hr.Context); err != nil {
		hr.Logger.Error("HealthCheckRepositoryImpl.Check() Ping DB ERROR, ", err)
		return false, nil
	}

	if err := hr.Redis.Ping(hr.Context).Err(); err != nil {
		hr.Logger.Error("HealthCheckRepositoryImpl.Check() Ping Redis ERROR, ", err)
		return false, nil
	}

	return true, nil
}
