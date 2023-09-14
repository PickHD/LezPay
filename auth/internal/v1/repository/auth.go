package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/LezPay/auth/internal/v1/config"
	"github.com/PickHD/LezPay/auth/internal/v1/helper"
	"github.com/PickHD/LezPay/auth/internal/v1/model"
	"github.com/matcornic/hermes/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/gomail.v2"
)

type (
	// AuthRepository is an interface that has all the function to be implemented inside auth repository
	AuthRepository interface {
		SendMailVerification(req interface{}, verificationType model.VerificationType) error
		GetVerificationByCode(ctx context.Context, code string, verificationType model.VerificationType) (string, error)
	}

	// AuthRepositoryImpl is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Redis   *redis.Client
		Tracer  *trace.TracerProvider
		Mailer  *gomail.Dialer
	}
)

// NewAuthRepository return new instances auth repository
func NewAuthRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, rds *redis.Client, tracer *trace.TracerProvider, mailer *gomail.Dialer) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Redis:   rds,
		Tracer:  tracer,
		Mailer:  mailer,
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

		result, err := ar.constructEmailTemplate(reqRegister.Fullname, emailLink, model.RegisterVerification)
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification constructEmailTemplate ERROR ", err)

			return err
		}

		err = ar.sendMail(ar.Config.Mailer.Sender, []string{reqRegister.Email}, reqRegister.Email, "Registration Confirmations", "Please Complete the Verification of your Request Registration", result)
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification sendMail ERROR ", err)

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

		result, err := ar.constructEmailTemplate("", emailLink, model.ForgotPasswordVerification)
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification constructEmailTemplate ERROR ", err)

			return err
		}

		err = ar.sendMail(ar.Config.Mailer.Sender, []string{reqForgotPassRequest.Email}, reqForgotPassRequest.Email, "Forgot Password Confirmations", "Please Complete the Verification of your Request Forgot Password", result)
		if err != nil {
			ar.Logger.Error("AuthRepositoryImpl.SendMailVerification sendMail ERROR ", err)

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

func (ar *AuthRepositoryImpl) sendMail(from string, to []string, cc string, ccTitle string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetAddressHeader("Cc", cc, ccTitle)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	err := ar.Mailer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepositoryImpl) constructEmailTemplate(userName string, emailLink string, verificationType model.VerificationType) (string, error) {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "LezPay",
			Link: "https://github.com/PickHD/LezPay",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: userName,
			Intros: []string{
				"Selamat Datang di LezPay! Tinggal sedikit lagi nih kamu bisa menggunakan wallet nya",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Konfirmasi akunmu disini :",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Konfirmasi",
						Link:  emailLink,
					},
				},
			},
			Outros: []string{
				"Butuh bantuan? balas email ini, akan kami bantu sebisa mungkin",
			},
		},
	}

	switch verificationType {
	case model.RegisterVerification:
	case model.ForgotPasswordVerification:
		email.Body.Intros = append(email.Body.Intros[:0], email.Body.Intros[0+1:]...)

		email.Body.Intros = append(email.Body.Intros, "Lupa kata sandimu? silahkan konfirmasi dengan menekan tombol konfirmasi untuk melanjutkan penggantian kata sandi")
		email.Body.Actions[0].Instructions = "Konfirmasi ganti kata sandimu disini :"
	}

	result, err := h.GenerateHTML(email)
	if err != nil {
		ar.Logger.Error(err)

		return "", err
	}

	return result, nil
}
