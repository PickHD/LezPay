package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/LezPay/auth/internal/config"
	"github.com/PickHD/LezPay/auth/internal/helper"
	"github.com/PickHD/LezPay/auth/internal/model"
	"github.com/PickHD/LezPay/auth/internal/requester"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// AuthRepository is an interface that has all the function to be implemented inside auth repository
	AuthRepository interface {
		SendMailVerification(req interface{}, verificationType model.VerificationType) error
		GetVerificationByCode(ctx context.Context, code string, verificationType model.VerificationType) (string, error)
	}

	// AuthRepositoryImpl is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepositoryImpl struct {
		Context          context.Context
		Config           *config.Configuration
		Logger           *logrus.Logger
		Redis            *redis.Client
		Tracer           *trace.TracerProvider
		UtilityRequester requester.UtilityRequester
	}
)

// NewAuthRepository return new instances auth repository
func NewAuthRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, rds *redis.Client, tracer *trace.TracerProvider, utilityRequester requester.UtilityRequester) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		Context:          ctx,
		Config:           config,
		Logger:           logger,
		Redis:            rds,
		Tracer:           tracer,
		UtilityRequester: utilityRequester,
	}
}

func (ar *AuthRepositoryImpl) SendMailVerification(req interface{}, verificationType model.VerificationType) error {
	tr := ar.Tracer.Tracer("Auth-SendMailVerification Repository")
	_, span := tr.Start(ar.Context, "Start SendMailVerification")
	defer span.End()

	codeVerification := helper.RandomStringBytesMaskImprSrcSB(25)
	expiredCodeDuration := time.Minute * time.Duration(ar.Config.Redis.TTL)

	switch verificationType {
	case model.RegisterVerification:
		reqRegister, ok := req.(*model.RegisterRequest)
		if !ok {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification type assertion ERROR, ", model.NewError(model.Type, "failed to convert data"))

			return model.NewError(model.Type, "failed to convert data")
		}

		err := ar.Redis.SetEx(ar.Context, fmt.Sprintf(model.VerificationKey, model.RegisterVerification, codeVerification), reqRegister.Email, expiredCodeDuration).Err()
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification SetEx ERROR, ", err)

			return err
		}

		emailLink := fmt.Sprintf("http://localhost:%d/v1/register/verify?code=%s&user_type=%s", ar.Config.Server.AppPort, codeVerification, reqRegister.UserType)

		_, err = ar.UtilityRequester.SendMail(ar.Context, &requester.SendMailRequest{
			Subject:          "Please Complete the Verification of your Request Registration",
			Cc:               reqRegister.Email,
			CcTitle:          "Registration Confirmations",
			To:               []string{reqRegister.Email},
			Username:         reqRegister.Fullname,
			Link:             emailLink,
			UserType:         string(reqRegister.UserType),
			VerificationType: string(model.RegisterVerification),
		})
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification UtilityRequester.SendMail ERROR, ", err)

			return err
		}

		return nil
	case model.ForgotPasswordVerification:
		reqForgotPassRequest, ok := req.(*model.ForgotPasswordRequest)
		if !ok {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification type assertion ERROR, ", model.NewError(model.Type, "failed to convert data"))

			return model.NewError(model.Type, "failed to convert data")
		}

		err := ar.Redis.SetEx(ar.Context, fmt.Sprintf(model.VerificationKey, model.ForgotPasswordVerification, codeVerification), reqForgotPassRequest.Email, expiredCodeDuration).Err()
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification SetEx ERROR, ", err)

			return err
		}

		emailLink := fmt.Sprintf("http://localhost:%d/v1/forgot-password/verify?code=%s&user_type=%s", ar.Config.Server.AppPort, codeVerification, reqForgotPassRequest.UserType)

		_, err = ar.UtilityRequester.SendMail(ar.Context, &requester.SendMailRequest{
			Subject:          "Please Complete the Verification of your Request Registration",
			Cc:               reqForgotPassRequest.Email,
			CcTitle:          "Please Complete the Verification of your Request Forgot Password",
			To:               []string{reqForgotPassRequest.Email},
			Username:         reqForgotPassRequest.Email,
			Link:             emailLink,
			UserType:         string(reqForgotPassRequest.UserType),
			VerificationType: string(model.RegisterVerification),
		})
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification UtilityRequester.SendMail ERROR, ", err)

			return err
		}

		return nil
	}

	return model.NewError(model.Validation, "invalid verification_type")
}

func (ar *AuthRepositoryImpl) GetVerificationByCode(ctx context.Context, code string, verificationType model.VerificationType) (string, error) {
	tr := ar.Tracer.Tracer("Auth-GetVerificationByCode repository")
	ctx, span := tr.Start(ctx, "Start GetVerificationByCode")
	defer span.End()

	result := ar.Redis.Get(ctx, fmt.Sprintf(model.VerificationKey, verificationType, code))
	if result.Err() != nil {
		ar.Logger.Error("AuthRepositoryImpl.GetVerificationByCode Get ERROR, ", result.Err())

		return "", result.Err()
	}

	return result.Val(), nil
}
