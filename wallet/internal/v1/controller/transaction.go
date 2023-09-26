package controller

import (
	"context"
	"strings"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/PickHD/LezPay/wallet/internal/v1/service"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	walletpb "github.com/PickHD/LezPay/wallet/pkg/proto/v1/wallet"
)

type (
	// TransactionController is an interface that has all the function to be implemented inside transaction controller
	TransactionController interface {
		CreateTransaction(ctx context.Context, req *walletpb.CreateTransactionRequest) (*walletpb.CreateTransactionResponse, error)
	}

	// TransactionControllerImpl is an app Transaction struct that consists of all the dependencies needed for transaction controller
	TransactionControllerImpl struct {
		Context        context.Context
		Config         *config.Configuration
		Tracer         *trace.TracerProvider
		TransactionSvc service.TransactionService
		walletpb.UnimplementedTransactionServiceServer
	}
)

// NewTransactionController return new instances transaction controller
func NewTransactionController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, transactionSvc service.TransactionService) *TransactionControllerImpl {
	return &TransactionControllerImpl{
		Context:        ctx,
		Config:         config,
		Tracer:         tracer,
		TransactionSvc: transactionSvc,
	}
}

func (tc *TransactionControllerImpl) CreateTransaction(ctx context.Context, req *walletpb.CreateTransactionRequest) (*walletpb.CreateTransactionResponse, error) {
	tr := tc.Tracer.Tracer("Wallet-CreateTransaction Controller")
	_, span := tr.Start(ctx, "Start CreateTransaction")
	defer span.End()

	newTrx := model.CreateTransactionRequest{
		CustomerID:       req.GetCustomerId(),
		PaymentChannelID: req.GetPaymentChannelId(),
		Amount:           req.GetAmount(),
		TypeTransaction:  model.TransactionType(req.GetTypeTransaction()),
	}

	data, err := tc.TransactionSvc.CreateTransaction(&newTrx)
	if err != nil {
		if strings.Contains(err.Error(), string(model.NotFound)) {
			return nil, status.Error(codes.NotFound, "Data Not Found")
		}

		return nil, status.Errorf(codes.Internal, "Failed Create Transaction %s", err.Error())
	}

	return &walletpb.CreateTransactionResponse{
		TransactionId: data.ID,
		Status:        data.Status,
	}, nil
}
