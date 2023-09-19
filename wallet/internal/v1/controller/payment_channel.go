package controller

import (
	"context"
	"net/http"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/helper"
	"github.com/PickHD/LezPay/wallet/internal/v1/service"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// PaymentChannelController is an interface that has all the function to be implemented inside payment channel controller
	PaymentChannelController interface {
		GetActiveList(ctx echo.Context) error
	}

	// PaymentChannelControllerImpl is an app payment channel struct that consists of all the dependencies needed for payment channel controller
	PaymentChannelControllerImpl struct {
		Context           context.Context
		Config            *config.Configuration
		Tracer            *trace.TracerProvider
		PaymentChannelSvc service.PaymentChannelService
	}
)

// NewPaymentChannelController return new instances payment channel controller
func NewPaymentChannelController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, paymentChannelSvc service.PaymentChannelService) *PaymentChannelControllerImpl {
	return &PaymentChannelControllerImpl{
		Context:           ctx,
		Config:            config,
		Tracer:            tracer,
		PaymentChannelSvc: paymentChannelSvc,
	}
}

// Check godoc
// @Summary      Get Active List Payment Channels
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /payment-channel [get]
func (pc *PaymentChannelControllerImpl) GetActiveList(ctx echo.Context) error {
	tr := pc.Tracer.Tracer("Wallet-GetActiveList Controller")
	_, span := tr.Start(pc.Context, "Start GetActiveList")
	defer span.End()

	data, err := pc.PaymentChannelSvc.GetActiveList(pc.Context)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed Get Active List Payment Channel", nil, nil, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get Active List Payment Channel", data, nil, nil)
}
