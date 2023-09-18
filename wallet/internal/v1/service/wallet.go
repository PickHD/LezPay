package service

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/PickHD/LezPay/wallet/internal/v1/repository"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// WalletService is an interface that has all the function to be implemented inside wallet service
	WalletService interface {
		CreateWallet(req *model.CreateWalletRequest) (*model.CreateWalletResponse, error)
		GetCustomerWalletByCustomerID(customerID uint64) (*model.Wallet, error)
	}

	// WalletServiceImpl is an app wallet struct that consists of all the dependencies needed for Wallet service
	WalletServiceImpl struct {
		Context    context.Context
		Config     *config.Configuration
		Tracer     *trace.TracerProvider
		WalletRepo repository.WalletRepository
	}
)

// NewWalletService return new instances wallet service
func NewWalletService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, walletRepo repository.WalletRepository) *WalletServiceImpl {
	return &WalletServiceImpl{
		Context:    ctx,
		Config:     config,
		Tracer:     tracer,
		WalletRepo: walletRepo,
	}
}

func (ws *WalletServiceImpl) CreateWallet(req *model.CreateWalletRequest) (*model.CreateWalletResponse, error) {
	tr := ws.Tracer.Tracer("Wallet-CreateWallet Service")
	_, span := tr.Start(ws.Context, "Start CreateWallet")
	defer span.End()

	Id, err := ws.WalletRepo.CreateWallet(req)
	if err != nil {
		return nil, err
	}

	return &model.CreateWalletResponse{
		ID: uint64(Id),
	}, nil
}

func (ws *WalletServiceImpl) GetCustomerWalletByCustomerID(customerID uint64) (*model.Wallet, error) {
	tr := ws.Tracer.Tracer("Wallet-GetCustomerWalletByCustomerID Service")
	_, span := tr.Start(ws.Context, "Start GetCustomerWalletByCustomerID")
	defer span.End()

	data, err := ws.WalletRepo.GetCustomerWalletByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
