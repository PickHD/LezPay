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
	// TransactionRepository is an interface that has all the function to be implemented inside Transaction repository
	TransactionRepository interface {
		CreateTransaction(req *model.CreateTransactionRequest) (int64, string, error)
	}

	// TransactionRepositoryImpl is an app Transaction struct that consists of all the dependencies needed for transaction repository
	TransactionRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *pgxpool.Pool
		Redis   *redis.Client
	}
)

// NewTransactionRepository return new instances transaction repository
func NewTransactionRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
		Redis:   rds,
	}
}

func (tr *TransactionRepositoryImpl) CreateTransaction(req *model.CreateTransactionRequest) (int64, string, error) {
	trc := tr.Tracer.Tracer("Wallet-CreateTransaction Repository")
	_, span := trc.Start(tr.Context, "Start CreateTransaction")
	defer span.End()

	// begin tx
	tx, err := tr.DB.Begin(tr.Context)
	if err != nil {
		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction DB.Begin ERROR", err)

		return 0, "", err
	}

	sql := `
		INSERT INTO 
			transaction (id,customer_id,payment_channel_id,amount,type, status)
		VALUES
			($1, $2, $3, $4, $5, $6)
	`
	sqlHistory := `
		INSERT INTO
			transaction_history (id,transaction_id,status)
		VALUES
			($1, $2, $3)
	`

	id, err := helper.GenerateSnowflakeID()
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(tr.Context)
		if errRollback != nil {
			tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Rollback ERROR", errRollback)

			return 0, "", errRollback
		}

		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction GenerateSnowflakeID ERROR", err)

		return 0, "", err
	}

	idHistory, err := helper.GenerateSnowflakeID()
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(tr.Context)
		if errRollback != nil {
			tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Rollback ERROR", errRollback)

			return 0, "", errRollback
		}

		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction GenerateSnowflakeID ERROR", err)

		return 0, "", err
	}

	_, err = tx.Exec(tr.Context, sql, id, req.CustomerID, req.PaymentChannelID, req.Amount, string(req.TypeTransaction), string(model.TransactionPending))
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(tr.Context)
		if errRollback != nil {
			tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Rollback ERROR", errRollback)

			return 0, "", errRollback
		}

		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Exec ERROR", err)

		return 0, "", err
	}

	_, err = tx.Exec(tr.Context, sqlHistory, idHistory, id, string(model.TransactionPending))
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(tr.Context)
		if errRollback != nil {
			tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Rollback ERROR", errRollback)

			return 0, "", errRollback
		}

		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Exec ERROR", err)

		return 0, "", err
	}

	// do commit tx
	err = tx.Commit(tr.Context)
	if err != nil {
		tr.Logger.Error("WalletRepositoryImpl.CreateTransaction tx.Commit ERROR", err)

		return 0, "", err
	}

	return id, string(model.TransactionPending), nil
}
