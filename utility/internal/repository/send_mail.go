package repository

import (
	"context"

	"github.com/PickHD/LezPay/utility/internal/config"
	"github.com/PickHD/LezPay/utility/internal/model"
	"github.com/matcornic/hermes/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/gomail.v2"
)

type (
	// SendMailRepository is an interface that has all the function to be implemented inside send mail repository
	SendMailRepository interface {
		Send(ctx context.Context, req *model.SendMailRequest) error
	}

	// SendMailRepositoryImpl is an app send mail struct that consists of all the dependencies needed for health check repository
	SendMailRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		Mailer  *gomail.Dialer
	}
)

// NewSendMailRepository return new instances send mail repository
func NewSendMailRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, mailer *gomail.Dialer) *SendMailRepositoryImpl {
	return &SendMailRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		Mailer:  mailer,
	}
}

func (sr *SendMailRepositoryImpl) Send(ctx context.Context, req *model.SendMailRequest) error {
	tr := sr.Tracer.Tracer("Utility-Send Repository")
	_, span := tr.Start(sr.Context, "Start Send")
	defer span.End()

	result, err := sr.constructEmailTemplate(req.Username, req.Link, req.VerificationType)
	if err != nil {
		sr.Logger.Error("SendMailRepositoryImpl.Send constructEmailTemplate ERROR ", err)

		return err
	}

	err = sr.sendMail(sr.Config.Mailer.Sender, []string{req.To[0]}, req.Cc, req.CcTitle, req.Subject, result)
	if err != nil {
		sr.Logger.Error("SendMailRepositoryImpl.Send sendMail ERROR ", err)

		return err
	}

	return nil

}

func (sr *SendMailRepositoryImpl) sendMail(from string, to []string, cc string, ccTitle string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetAddressHeader("Cc", cc, ccTitle)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	err := sr.Mailer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SendMailRepositoryImpl) constructEmailTemplate(userName string, emailLink string, verificationType model.VerificationType) (string, error) {
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
		sr.Logger.Error(err)

		return "", err
	}

	return result, nil
}
