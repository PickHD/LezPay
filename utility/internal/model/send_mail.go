package model

import "regexp"

type (
	// UserType consist type of user (CUSTOMER/MERCHANT)
	UserType string

	// VerificationType consist type of mails verification (Registration/)
	VerificationType string

	// SendMailRequest consist request send mails
	SendMailRequest struct {
		Subject          string           `json:"subject"`
		To               []string         `json:"to"`
		Cc               string           `json:"cc"`
		CcTitle          string           `json:"cc_title"`
		Username         string           `json:"username"`
		Link             string           `json:"link"`
		UserType         UserType         `json:"user_type"`
		VerificationType VerificationType `json:"verification_type"`
	}
)

const (
	Customer UserType = "customer"
	Merchant UserType = "merchant"

	RegisterVerification       VerificationType = "register_verification"
	ForgotPasswordVerification VerificationType = "forgot_password_verification"
)

var (
	IsValidEmail, _ = regexp.Compile(`^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
)
