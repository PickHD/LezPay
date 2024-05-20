package service

import (
	"context"
	"time"

	"github.com/PickHD/LezPay/auth/internal/config"
	"github.com/PickHD/LezPay/auth/internal/helper"
	"github.com/PickHD/LezPay/auth/internal/model"
	"github.com/PickHD/LezPay/auth/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/customer"
	merchantpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/merchant"
	walletpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/wallet"
)

type (
	// AuthService is an interface that has all the function to be implemented inside auth service
	AuthService interface {
		RegisterCustomerOrMerchant(req *model.RegisterRequest) error
		VerifyRegisterOrForgotPasswordCode(ctx context.Context, code string, userType model.UserType, verifyType model.VerificationType) (*model.VerifyCodeResponse, error)
		LoginCustomerOrMerchant(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
		ForgotPasswordCustomerOrMerchant(ctx context.Context, req *model.ForgotPasswordRequest) error
		ResetPasswordCustomerOrMerchant(ctx context.Context, req *model.ResetPasswordRequest, code string, userType model.UserType) error
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
		// check request email
		_, err := as.CustomerClients.GetCustomerDetailsByEmail(as.Context, &customerpb.GetCustomerDetailsByEmailRequest{Email: req.Email})
		if err != nil {
			switch status.Code(err) {
			case codes.NotFound:
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

				err = as.AuthRepo.SendMailVerification(req, model.RegisterVerification)
				if err != nil {
					return err
				}

				return nil
			}

			return err
		}

		return model.NewError(model.Validation, "Email already exists")
	case model.Merchant:
		// check request email
		_, err := as.MerchantClients.GetMerchantDetailsByEmail(as.Context, &merchantpb.GetMerchantDetailsByEmailRequest{Email: req.Email})
		if err != nil {
			switch status.Code(err) {
			case codes.NotFound:
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

				err = as.AuthRepo.SendMailVerification(req, model.RegisterVerification)
				if err != nil {
					return err
				}

				return nil
			}

			return err
		}

		return model.NewError(model.Validation, "Email already exists")
	}

	return model.NewError(model.Validation, "Invalid user type")
}

func (as *AuthServiceImpl) VerifyRegisterOrForgotPasswordCode(ctx context.Context, code string, userType model.UserType, verifyType model.VerificationType) (*model.VerifyCodeResponse, error) {
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
		switch status.Code(err) {
		case codes.NotFound:
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

func (as *AuthServiceImpl) ForgotPasswordCustomerOrMerchant(ctx context.Context, req *model.ForgotPasswordRequest) error {
	tr := as.Tracer.Tracer("Auth-ForgotPasswordCustomerOrMerchant service")
	_, span := tr.Start(ctx, "Start ForgotPasswordCustomerOrMerchant")
	defer span.End()

	err := as.validateForgotPasswordRequest(req)
	if err != nil {
		return err
	}

	_, err = as.CustomerClients.GetCustomerDetailsByEmail(ctx, &customerpb.GetCustomerDetailsByEmailRequest{Email: req.Email})
	if err != nil {
		// if not found, check to merchant client
		switch status.Code(err) {
		case codes.NotFound:
			_, err = as.MerchantClients.GetMerchantDetailsByEmail(ctx, &merchantpb.GetMerchantDetailsByEmailRequest{Email: req.Email})
			if err != nil {
				return err
			}

			err = as.AuthRepo.SendMailVerification(req, model.ForgotPasswordVerification)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	err = as.AuthRepo.SendMailVerification(req, model.ForgotPasswordVerification)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthServiceImpl) ResetPasswordCustomerOrMerchant(ctx context.Context, req *model.ResetPasswordRequest, code string, userType model.UserType) error {
	tr := as.Tracer.Tracer("Auth-ResetPasswordUser service")
	ctx, span := tr.Start(ctx, "Start ResetPasswordUser")
	defer span.End()

	err := as.validateResetPasswordRequest(req)
	if err != nil {
		return err
	}

	getEmail, err := as.AuthRepo.GetVerificationByCode(ctx, code, model.ForgotPasswordVerification)
	if err != nil {
		if err == redis.Nil {
			return model.NewError(model.NotFound, "code not found / expired")
		}

		return err
	}

	hashedNewPass, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	switch userType {
	case model.Customer:
		reqUpdateNewPassword := &customerpb.UpdateCustomerPasswordByEmailRequest{
			Email:    getEmail,
			Password: hashedNewPass,
		}
		_, err := as.CustomerClients.UpdateCustomerPasswordByEmail(as.Context, reqUpdateNewPassword)
		if err != nil {
			return err
		}
	case model.Merchant:
		reqUpdateNewPassword := &merchantpb.UpdateMerchantPasswordByEmailRequest{
			Email:    getEmail,
			Password: hashedNewPass,
		}
		_, err := as.MerchantClients.UpdateMerchantPasswordByEmail(as.Context, reqUpdateNewPassword)
		if err != nil {
			return err
		}
	}

	return nil
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

func (as *AuthServiceImpl) validateForgotPasswordRequest(req *model.ForgotPasswordRequest) error {
	if !model.IsValidEmail.MatchString(req.Email) {
		return model.NewError(model.Validation, "invalid email")
	}

	return nil
}

func (as *AuthServiceImpl) validateResetPasswordRequest(req *model.ResetPasswordRequest) error {
	if req.Password == "" {
		return model.NewError(model.Validation, "password required")
	}

	ok := helper.IsValid(req.Password)
	if !ok {
		return model.NewError(model.Validation, "password must min length 7, and at least has 1 each upper,lower,number,special")
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
