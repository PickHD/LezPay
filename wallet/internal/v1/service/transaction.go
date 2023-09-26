package service

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/PickHD/LezPay/wallet/internal/v1/repository"
	"github.com/jackc/pgx"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// TransactionService is an interface that has all the function to be implemented inside Transaction service
	TransactionService interface {
		CreateTransaction(req *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
	}

	// TransactionServiceImpl is an app Transaction struct that consists of all the dependencies needed for Transaction service
	TransactionServiceImpl struct {
		Context            context.Context
		Config             *config.Configuration
		Tracer             *trace.TracerProvider
		PaymentChannelRepo repository.PaymentChannelRepository
		WalletRepo         repository.WalletRepository
		TransactionRepo    repository.TransactionRepository
	}
)

// NewTransactionService return new instances Transaction service
func NewTransactionService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, paymentChannelRepo repository.PaymentChannelRepository, walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		Context:            ctx,
		Config:             config,
		Tracer:             tracer,
		PaymentChannelRepo: paymentChannelRepo,
		WalletRepo:         walletRepo,
		TransactionRepo:    transactionRepo,
	}
}

func (ts *TransactionServiceImpl) CreateTransaction(req *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	tr := ts.Tracer.Tracer("Wallet-CreateTransaction Service")
	_, span := tr.Start(ts.Context, "Start CreateTransaction")
	defer span.End()

	// check payment channel is exists & active or not
	// if exists & active continue
	_, err := ts.PaymentChannelRepo.GetByID(ts.Context, int64(req.PaymentChannelID))
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, model.NewError(model.NotFound, "Payment Channel Not Found / Not Active")
		}

		return nil, err
	}

	// check customer wallet is exists or not
	// if exists continue
	_, err = ts.WalletRepo.GetCustomerWalletByCustomerID(req.CustomerID)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, model.NewError(model.NotFound, "Wallet customer not found")
		}

		return nil, err
	}

	transactionID, status, err := ts.TransactionRepo.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &model.CreateTransactionResponse{
		ID:     uint64(transactionID),
		Status: status,
	}, nil
}
