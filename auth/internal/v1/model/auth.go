package model

import (
	"regexp"
	"time"
)

type (
	// UserType determinate users is customer or merchant
	UserType string

	// VerificationType consist type of verification
	VerificationType string

	// RegisterRequest consist request of register customer/merchant
	RegisterRequest struct {
		Fullname    string   `json:"full_name,omitempty"`
		Email       string   `json:"email"`
		PhoneNumber string   `json:"phone_number,omitempty"`
		Password    string   `json:"password"`
		Pin         string   `json:"pin"`
		UserType    UserType `json:"user_type"`
	}

	// VerifyCodeResponse consist response of verify code
	VerifyCodeResponse struct {
		IsVerified bool `json:"is_verified"`
	}

	// LoginRequest consist request of login customer/merchant
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// LoginResponse consist response of login customer/merchant
	LoginResponse struct {
		AccessToken string    `json:"access_token"`
		ExpiredAt   time.Time `json:"expired_at"`
		Type        string    `json:"type"`
	}
)

const (
	Customer UserType = "customer"
	Merchant UserType = "merchant"

	RegisterVerification       VerificationType = "register_verification"
	ForgotPasswordVerification VerificationType = "forgot_password_verification"
)

var (
	IsValidEmail, _    = regexp.Compile(`^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
	IsValidPin         = regexp.MustCompile(`^\d{6}$`)
	IsValidPhoneNumber = regexp.MustCompile(`^(?:\+62|0)[0-9]{9,15}$`)
)
