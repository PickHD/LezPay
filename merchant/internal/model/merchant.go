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

	// GetMerchantDetailsByEmailRequest consist request getting merchant details by email
	GetMerchantDetailsByEmailRequest struct {
		Email string
	}

	// GetMerchantDetailsByEmailResponse consist data response getting merchant details by email
	GetMerchantDetailsByEmailResponse struct {
		ID          int64
		FullName    string
		Email       string
		PhoneNumber string
		Password    string
	}

	// UpdateMerchantPasswordByEmailRequest consist request update password merchant by email
	UpdateMerchantPasswordByEmailRequest struct {
		Email    string
		Password string
	}

	// UpdateMerchantPasswordByEmailResponse consist data response updating password merchant
	UpdateMerchantPasswordByEmailResponse struct {
		Email string
	}
)
