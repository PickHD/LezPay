package model

type (
	// CreateCustomerRequest consist request creating customers
	CreateCustomerRequest struct {
		FullName    string
		Email       string
		PhoneNumber string
		Password    string
		Pin         string
	}

	// CreateCustomerResponse consist data response of creating customers
	CreateCustomerResponse struct {
		ID         int64
		IsVerified bool
	}

	// UpdateVerifiedCustomerResponse consist data response of updating verified customers
	UpdateVerifiedCustomerResponse struct {
		IsVerified bool
	}

	// GetCustomerIDByEmailRequest consist request getting customer id by email
	GetCustomerIDByEmailRequest struct {
		Email string
	}

	// GetCustomerIDByEmailResponse consist data response getting customer id by email
	GetCustomerIDByEmailResponse struct {
		CustomerID int64
	}

	// GetCustomerDetailsByEmailRequest consist request getting customer details by email
	GetCustomerDetailsByEmailRequest struct {
		Email string
	}

	// GetCustomerDetailsByEmailResponse consist data response getting customer details by email
	GetCustomerDetailsByEmailResponse struct {
		ID          int64
		FullName    string
		Email       string
		PhoneNumber string
		Password    string
		Pin         string
	}

	// UpdateCustomerPasswordByEmailRequest consist request update password customer by email
	UpdateCustomerPasswordByEmailRequest struct {
		Email    string
		Password string
	}

	// UpdateCustomerPasswordByEmailResponse consist data response updating password customer
	UpdateCustomerPasswordByEmailResponse struct {
		Email string
	}

	// GetCustomerDashboardResponse consist data response get customer dashboard
	GetCustomerDashboardResponse struct {
		WalletID uint64 `json:"wallet_id"`
		Balance  int64  `json:"balance"`
	}
)
