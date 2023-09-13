package service

import (
	"context"
	"time"

	"github.com/PickHD/LezPay/auth/internal/v1/config"
	"github.com/PickHD/LezPay/auth/internal/v1/helper"
	"github.com/PickHD/LezPay/auth/internal/v1/model"
	"github.com/PickHD/LezPay/auth/internal/v1/repository"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"

	customerpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/customer"
	merchantpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/merchant"
	walletpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/wallet"
)

type (
	// AuthService is an interface that has all the function to be implemented inside auth service
	AuthService interface {
		RegisterCustomerOrMerchant(req *model.RegisterRequest) error
		VerifyRegisterCode(ctx context.Context, code string, userType model.UserType, verifyType model.VerificationType) (*model.VerifyCodeResponse, error)
		LoginCustomerOrMerchant(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	}

	// AuthServiceImpl is an app auth struct that consists of all the dependencies needed for auth service
	AuthServiceImpl struct {
		Context         context.Context
		Config          *config.Configuration
		Logger          *logrus.Logger
		Tracer          *trace.TracerProvider
		AuthRepo        repository.AuthRepository
		CustomerClients customerpb.CustomerServiceClient
		MerchantClients merchantpb.MerchantServiceClient
		WalletClients   walletpb.WalletServiceClient
	}
)

// NewAuthService return new instances auth service
func NewAuthService(
	ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, authRepo repository.AuthRepository, customerClients customerpb.CustomerServiceClient, merchantClients merchantpb.MerchantServiceClient, walletClients walletpb.WalletServiceClient,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		Context:         ctx,
		Config:          config,
		Logger:          logger,
		Tracer:          tracer,
		AuthRepo:        authRepo,
		CustomerClients: customerClients,
		MerchantClients: merchantClients,
		WalletClients:   walletClients,
	}
}

func (as *AuthServiceImpl) RegisterCustomerOrMerchant(req *model.RegisterRequest) error {
	tr := as.Tracer.Tracer("Auth-RegisterCustomerOrMerchant Service")
	_, span := tr.Start(as.Context, "Start RegisterCustomerOrMerchant")
	defer span.End()

	err := as.validateRegisterRequest(req)
	if err != nil {
		return err
	}

	hashPass, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	hashPin, err := helper.HashPassword(req.Pin)
	if err != nil {
		return err
	}

	switch req.UserType {
	case model.Customer:
		newCustomer := &customerpb.CustomerRequest{
			FullName:    req.Fullname,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Password:    hashPass,
			Pin:         hashPin,
		}

		_, err = as.CustomerClients.CreateCustomer(as.Context, newCustomer)
		if err != nil {
			return err
		}
	case model.Merchant:
		newMerchant := &merchantpb.MerchantRequest{
			FullName:    req.Fullname,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Password:    hashPass,
		}

		_, err = as.MerchantClients.CreateMerchant(as.Context, newMerchant)
		if err != nil {
			return err
		}
	}

	err = as.AuthRepo.SendMailRegisterVerification(req)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthServiceImpl) VerifyRegisterCode(ctx context.Context, code string, userType model.UserType, verifyType model.VerificationType) (*model.VerifyCodeResponse, error) {
	tr := as.Tracer.Tracer("Auth-VerifyRegisterCode service")
	ctx, span := tr.Start(ctx, "Start VerifyRegisterCode")
	defer span.End()

	email, err := as.AuthRepo.GetVerificationByCode(ctx, code, verifyType)
	if err != nil {
		if err == redis.Nil {
			return nil, model.NewError(model.NotFound, "code not found / expired")
		}

		return nil, err
	}

	switch verifyType {
	case model.RegisterVerification:
		switch userType {
		case model.Customer:
			// update customer verify status to true
			_, err := as.CustomerClients.UpdateVerifiedCustomer(as.Context, &customerpb.UpdateVerifiedCustomerRequest{Email: email})
			if err != nil {
				return nil, err
			}

			// get customer id by email
			data, err := as.CustomerClients.GetCustomerIDByEmail(as.Context, &customerpb.GetCustomerIDByEmailRequest{Email: email})
			if err != nil {
				return nil, err
			}

			// create wallet with 0 balance
			_, err = as.WalletClients.CreateWallet(as.Context, &walletpb.WalletRequest{CustomerId: data.GetCustomerId()})
			if err != nil {
				return nil, err
			}

		case model.Merchant:
			// update merchant verify status to true
			_, err := as.MerchantClients.UpdateVerifiedMerchant(as.Context, &merchantpb.UpdateVerifiedMerchantRequest{Email: email})
			if err != nil {
				return nil, err
			}
		}
	case model.ForgotPasswordVerification:
	}

	return &model.VerifyCodeResponse{
		IsVerified: true,
	}, nil
}

func (as *AuthServiceImpl) LoginCustomerOrMerchant(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	tr := as.Tracer.Tracer("Auth-LoginCustomerOrMerchant service")
	_, span := tr.Start(ctx, "Start LoginCustomerOrMerchant")
	defer span.End()

	err := as.validateLoginRequest(req)
	if err != nil {
		return nil, err
	}

	dataCustomer, err := as.CustomerClients.GetCustomerDetailsByEmail(ctx, &customerpb.GetCustomerDetailsByEmailRequest{Email: req.Email})
	if err != nil {
		// if not found, check to merchant client
		if err.Error() == "Email Not Found" {
			dataMerchant, err := as.MerchantClients.GetMerchantDetailsByEmail(ctx, &merchantpb.GetMerchantDetailsByEmailRequest{Email: req.Email})
			if err != nil {
				return nil, err
			}

			if !helper.CheckPasswordHash(dataMerchant.GetPassword(), req.Password) {
				return nil, model.NewError(model.Validation, "password invalid")
			}

			// generate access token jwt
			accessToken, expiredAt, err := as.generateJWT(dataMerchant.GetId(), req.Email)
			if err != nil {
				return nil, err
			}

			return &model.LoginResponse{
				AccessToken: accessToken,
				Type:        "Bearer",
				ExpiredAt:   time.Now().Add(expiredAt),
			}, nil
		}

		return nil, err
	}

	if !helper.CheckPasswordHash(dataCustomer.GetPassword(), req.Password) {
		return nil, model.NewError(model.Validation, "password invalid")
	}

	// generate access token jwt
	accessToken, expiredAt, err := as.generateJWT(dataCustomer.GetId(), req.Email)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
		Type:        "Bearer",
		ExpiredAt:   time.Now().Add(expiredAt),
	}, nil
}

func (as *AuthServiceImpl) validateRegisterRequest(req *model.RegisterRequest) error {
	if req.UserType != model.Customer && req.UserType != model.Merchant {
		return model.NewError(model.Validation, "user_type must customer or merchant")
	}

	if !model.IsValidEmail.MatchString(req.Email) {
		return model.NewError(model.Validation, "invalid email")
	}

	if len(req.Fullname) < 1 {
		return model.NewError(model.Validation, "full_name must not empty")
	}

	if !model.IsValidPhoneNumber.MatchString(req.PhoneNumber) {
		return model.NewError(model.Validation, "invalid phone number, format number using indonesian code (+62)")
	}

	if !helper.IsValid(req.Password) {
		return model.NewError(model.Validation, "password must have at least 1 number, 1 special character, 1 capital")
	}

	if req.UserType == model.Merchant && len(req.Pin) > 0 {
		return model.NewError(model.Validation, "merchant not using pin")
	}

	if req.UserType == model.Customer && !model.IsValidPin.MatchString(req.Pin) {
		return model.NewError(model.Validation, "pin must be 6 digit numbers")
	}

	return nil
}

func (as *AuthServiceImpl) validateLoginRequest(req *model.LoginRequest) error {
	if !model.IsValidEmail.MatchString(req.Email) {
		return model.NewError(model.Validation, "invalid email")
	}

	return nil
}

func (as *AuthServiceImpl) generateJWT(userID uint64, email string) (string, time.Duration, error) {
	var (
		payloadUserID  = "user_id"
		payloadEmail   = "email"
		payloadExpires = "exp"
		JWTExpire      = time.Duration(as.Config.Common.JWTExpire) * time.Hour
	)

	claims := jwt.MapClaims{}
	claims[payloadUserID] = userID
	claims[payloadEmail] = email
	claims[payloadExpires] = time.Now().Add(JWTExpire).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(as.Config.Secret.JWTSecret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, JWTExpire, nil
}
