package model

type (
	// CreateWalletRequest consist request creating wallet
	CreateWalletRequest struct {
		CustomerID uint64
	}

	// CreateWalletResponse consist data response of creating wallet
	CreateWalletResponse struct {
		ID uint64
	}

	// Wallet consist model of representation wallet table
	Wallet struct {
		ID         uint64
		CustomerID uint64
		Balance    int64
	}
)
