package requester

import (
	"context"
	"net/http"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// MidtransRequester is an interface that has all the function to be implemented inside midtrans requester
	MidtransRequester interface {
	}

	// MidtransRequesterImpl is an app health check struct that consists of all the dependencies needed for midtrans requester
	MidtransRequesterImpl struct {
		Context    context.Context
		Config     *config.Configuration
		Logger     *logrus.Logger
		Tracer     *trace.TracerProvider
		HttpClient *http.Client
	}
)

// NewRequesterMidtrans return new instances midtrans requester
func (m *MidtransRequesterImpl) NewRequesterMidtrans(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, client *http.Client) *MidtransRequesterImpl {
	return &MidtransRequesterImpl{
		Context:    ctx,
		Config:     config,
		Logger:     logger,
		Tracer:     tracer,
		HttpClient: client,
	}
}
