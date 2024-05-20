package service

import (
	"context"

	"github.com/PickHD/LezPay/utility/internal/config"
	"github.com/PickHD/LezPay/utility/internal/model"
	"github.com/PickHD/LezPay/utility/internal/repository"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// SendMailService is an interface that has all the function to be implemented inside send mail Service
	SendMailService interface {
		Send(ctx context.Context, req *model.SendMailRequest) error
	}

	// SendMailServiceImpl is an app send mail struct that consists of all the dependencies needed for send mail service
	SendMailServiceImpl struct {
		Context      context.Context
		Config       *config.Configuration
		Logger       *logrus.Logger
		Tracer       *trace.TracerProvider
		SendMailRepo repository.SendMailRepository
	}
)

// NewSendMailService return new instances send mail Service
func NewSendMailService(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, sendMailRepo repository.SendMailRepository) *SendMailServiceImpl {
	return &SendMailServiceImpl{
		Context:      ctx,
		Config:       config,
		Logger:       logger,
		Tracer:       tracer,
		SendMailRepo: sendMailRepo,
	}
}

func (ss *SendMailServiceImpl) Send(ctx context.Context, req *model.SendMailRequest) error {
	tr := ss.Tracer.Tracer("Utility-Send Service")
	_, span := tr.Start(ss.Context, "Start Send")
	defer span.End()

	err := ss.validateSendMail(req)
	if err != nil {
		return err
	}

	return ss.SendMailRepo.Send(ctx, req)
}

func (ss *SendMailServiceImpl) validateSendMail(req *model.SendMailRequest) error {
	if len(req.Subject) < 1 && len(req.To) < 1 && len(req.CcTitle) < 1 && len(req.Cc) < 1 && len(req.Username) < 1 && len(req.Link) < 1 && len(req.UserType) < 1 && len(req.VerificationType) < 1 {
		return model.NewError(model.Validation, "field cannot be empty")
	}

	for _, t := range req.To {
		if !model.IsValidEmail.MatchString(t) {
			return model.NewError(model.Validation, "invalid email")
		}
	}

	if req.UserType != model.Customer && req.UserType != model.Merchant {
		return model.NewError(model.Validation, "invalid user type")
	}

	if req.UserType != model.Customer && req.UserType != model.Merchant {
		return model.NewError(model.Validation, "invalid user type")
	}

	if req.VerificationType != model.RegisterVerification && req.VerificationType != model.ForgotPasswordVerification {
		return model.NewError(model.Validation, "invalid verification type")
	}

	return nil
}
