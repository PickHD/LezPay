package service

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/PickHD/LezPay/wallet/internal/v1/repository"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// PaymentChannelService is an interface that has all the function to be implemented inside payment channel service
	PaymentChannelService interface {
		GetActiveList(ctx context.Context) (*[]model.PaymentChannel, error)
	}

	// PaymentChannelServiceImpl is an app payment channel struct that consists of all the dependencies needed for payment channel service
	PaymentChannelServiceImpl struct {
		Context            context.Context
		Config             *config.Configuration
		Tracer             *trace.TracerProvider
		PaymentChannelRepo repository.PaymentChannelRepository
	}
)

// NewPaymentChannelService return new instances payment channel service
func NewPaymentChannelService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, paymentChannelRepo repository.PaymentChannelRepository) *PaymentChannelServiceImpl {
	return &PaymentChannelServiceImpl{
		Context:            ctx,
		Config:             config,
		Tracer:             tracer,
		PaymentChannelRepo: paymentChannelRepo,
	}
}

func (ps *PaymentChannelServiceImpl) GetActiveList(ctx context.Context) (*[]model.PaymentChannel, error) {
	tr := ps.Tracer.Tracer("Wallet-GetActiveList Service")
	_, span := tr.Start(ps.Context, "Start GetActiveList")
	defer span.End()

	return ps.PaymentChannelRepo.GetActiveList(ctx)
}
