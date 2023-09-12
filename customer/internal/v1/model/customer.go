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
)
