package repository

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/helper"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// WalletRepository is an interface that has all the function to be implemented inside wallet repository
	WalletRepository interface {
		CreateWallet(req *model.CreateWalletRequest) (int64, error)
	}

	// WalletRepositoryImpl is an app Wallet struct that consists of all the dependencies needed for wallet repository
	WalletRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *pgxpool.Pool
		Redis   *redis.Client
	}
)

// NewWalletRepository return new instances wallet repository
func NewWalletRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
		Redis:   rds,
	}
}

func (wr *WalletRepositoryImpl) CreateWallet(req *model.CreateWalletRequest) (int64, error) {
	tr := wr.Tracer.Tracer("Wallet-CreateWallet Repository")
	_, span := tr.Start(wr.Context, "Start CreateWallet")
	defer span.End()

	// begin tx
	tx, err := wr.DB.Begin(wr.Context)
	if err != nil {
		wr.Logger.Error("WalletRepositoryImpl.CreateWallet DB.Begin ERROR", err)

		return 0, err
	}

	sql := `
		INSERT INTO 
			wallet (id,customer_id,balance) 
		VALUES 
			($1,$2,$3)
		`

	id, err := helper.GenerateSnowflakeID()
	if err != nil {
		// do rollback tx
		err := tx.Rollback(wr.Context)
		if err != nil {
			wr.Logger.Error("WalletRepositoryImpl.CreateWallet tx.Rollback ERROR", err)

			return 0, err
		}

		wr.Logger.Error("WalletRepositoryImpl.CreateWallet GenerateSnowflakeID ERROR", err)

		return 0, err
	}

	_, err = tx.Exec(wr.Context, sql, id, req.CustomerID, 0)
	if err != nil {
		// do rollback tx
		err := tx.Rollback(wr.Context)
		if err != nil {
			wr.Logger.Error("WalletRepositoryImpl.CreateWallet tx.Rollback ERROR", err)

			return 0, err
		}

		wr.Logger.Error("WalletRepositoryImpl.CreateWallet tx.Exec ERROR", err)

		return 0, err
	}

	// do commit tx
	err = tx.Commit(wr.Context)
	if err != nil {
		wr.Logger.Error("WalletRepositoryImpl.CreateWallet tx.Commit ERROR", err)

		return 0, err
	}

	return id, nil
}
