package repository

import (
	"context"

	"github.com/PickHD/LezPay/wallet/internal/v1/config"
	"github.com/PickHD/LezPay/wallet/internal/v1/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// PaymentChannelRepository is an interface that has all the function to be implemented inside payment channel repository
	PaymentChannelRepository interface {
		GetActiveList(ctx context.Context) (*[]model.PaymentChannel, error)
	}

	// PaymentChannelRepositoryImpl is an app payment channel struct that consists of all the dependencies needed for payment channel repository
	PaymentChannelRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *pgxpool.Pool
		Redis   *redis.Client
	}
)

// NewPaymentChannelRepository return new instances payment channel repository
func NewPaymentChannelRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client) *PaymentChannelRepositoryImpl {
	return &PaymentChannelRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
		Redis:   rds,
	}
}

func (pr *PaymentChannelRepositoryImpl) GetActiveList(ctx context.Context) (*[]model.PaymentChannel, error) {
	tr := pr.Tracer.Tracer("Wallet-GetActiveList Repository")
	_, span := tr.Start(ctx, "Start GetActiveList")
	defer span.End()

	sql := `
		SELECT
			id,
			name,
			code,
			image_url,
			COALESCE(payment_instruction,'') AS payment_instruction,
			status
		FROM
			payment_channel
	`

	paymentChannels := []model.PaymentChannel{}

	rows, err := pr.DB.Query(ctx, sql)
	if err != nil {
		pr.Logger.Error("PaymentChannelRepositoryImpl.GetActiveList Query ERROR ", err)

		return nil, err
	}

	for rows.Next() {
		pChannel := model.PaymentChannel{}

		err = rows.Scan(
			&pChannel.ID,
			&pChannel.Name,
			&pChannel.Code,
			&pChannel.ImageURL,
			&pChannel.PaymentInstruction,
			&pChannel.Status,
		)
		if err != nil {
			pr.Logger.Error("PaymentChannelRepositoryImpl.GetActiveList rows.Scan ERROR ", err)

			return nil, err
		}

		paymentChannels = append(paymentChannels, pChannel)
	}
	defer rows.Close()

	return &paymentChannels, nil
}
