package repository

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/config"
	"github.com/PickHD/LezPay/wallet/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// WalletRepository is an interface that has all the function to be implemented inside wallet repository
	WalletRepository interface {
		CreateWallet(req *model.CreateWalletRequest) (int64, error)
		GetCustomerWalletByCustomerID(customerID uint64) (*model.Wallet, error)
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

	var id int64

	sql := `
		INSERT INTO 
			wallet (id,customer_id,balance) 
		VALUES 
			($1,$2,$3) RETURNING id
		`

	row := tx.QueryRow(wr.Context, sql, id, req.CustomerID, 0)

	err = row.Scan(&id)
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(wr.Context)
		if errRollback != nil {
			wr.Logger.Error("WalletRepositoryImpl.CreateWallet tx.Rollback ERROR ", errRollback)

			return 0, errRollback
		}

		wr.Logger.Error("WalletRepositoryImpl.CreateWallet row.Scan ERROR ", err)

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

func (wr *WalletRepositoryImpl) GetCustomerWalletByCustomerID(customerID uint64) (*model.Wallet, error) {
	tr := wr.Tracer.Tracer("Wallet-GetCustomerWalletByCustomerID Repository")
	_, span := tr.Start(wr.Context, "Start GetCustomerWalletByCustomerID")
	defer span.End()

	sql := `
		SELECT
			id,
			customer_id,
			balance
		FROM 
			wallet
		WHERE
			customer_id = $1
	`

	row := wr.DB.QueryRow(wr.Context, sql, customerID)

	data := model.Wallet{}
	err := row.Scan(&data.ID, &data.CustomerID, &data.Balance)
	if err != nil {
		wr.Logger.Error("WalletRepositoryImpl.GetCustomerWalletByCustomerID row.Scan ERROR ", err)

		return nil, err
	}

	return &data, nil
}
