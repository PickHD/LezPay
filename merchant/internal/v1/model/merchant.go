package model

type (
	// CreateMerchantRequest consist request creating merchants
	CreateMerchantRequest struct {
		FullName    string
		Email       string
		PhoneNumber string
		Password    string
	}

	// CreateMerchantResponse consist data response of creating merchants
	CreateMerchantResponse struct {
		ID         int64
		IsVerified bool
	}

	// UpdateVerifiedMerchantResponse consist data response of updating verified merchants
	UpdateVerifiedMerchantResponse struct {
		IsVerified bool
	}
)
