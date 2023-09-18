package controller

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/PickHD/LezPay/wallet/internal/v1/service"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	walletpb "github.com/PickHD/LezPay/wallet/pkg/proto/v1/wallet"
)

type (
	// WalletController is an interface that has all the function to be implemented inside Wallet controller
	WalletController interface {
		CreateWallet(ctx context.Context, req *walletpb.WalletRequest) (*walletpb.WalletResponse, error)
		GetCustomerWallet(ctx context.Context, req *walletpb.GetCustomerWalletRequest) (*walletpb.GetCustomerWalletResponse, error)
	}

	// WalletControllerImpl is an app Wallet struct that consists of all the dependencies needed for Wallet controller
	WalletControllerImpl struct {
		Context   context.Context
		Config    *config.Configuration
		Tracer    *trace.TracerProvider
		WalletSvc service.WalletService
		walletpb.UnimplementedWalletServiceServer
	}
)

// NewWalletController return new instances wallet controller
func NewWalletController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, walletSvc service.WalletService) *WalletControllerImpl {
	return &WalletControllerImpl{
		Context:   ctx,
		Config:    config,
		Tracer:    tracer,
		WalletSvc: walletSvc,
	}
}

func (wc *WalletControllerImpl) CreateWallet(ctx context.Context, req *walletpb.WalletRequest) (*walletpb.WalletResponse, error) {
	tr := wc.Tracer.Tracer("Wallet-CreateWallet Controller")
	_, span := tr.Start(ctx, "Start CreateWallet")
	defer span.End()

	newWallet := model.CreateWalletRequest{
		CustomerID: req.GetCustomerId(),
	}

	data, err := wc.WalletSvc.CreateWallet(&newWallet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Create Wallet %s", err.Error())
	}

	return &walletpb.WalletResponse{
		Id: data.ID,
	}, nil
}

func (wc *WalletControllerImpl) GetCustomerWallet(ctx context.Context, req *walletpb.GetCustomerWalletRequest) (*walletpb.GetCustomerWalletResponse, error) {
	tr := wc.Tracer.Tracer("Wallet-GetCustomerWallet Controller")
	_, span := tr.Start(ctx, "Start GetCustomerWallet")
	defer span.End()

	data, err := wc.WalletSvc.GetCustomerWalletByCustomerID(req.GetCustomerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Get Customer Wallet %s", err.Error())
	}

	return &walletpb.GetCustomerWalletResponse{
		Id:      data.ID,
		Balance: data.Balance,
	}, nil
}
