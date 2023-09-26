package model

type (
	// TransactionStatus consist enum of transaction status
	TransactionStatus string

	// CreateTransactionRequest consist request creating transactions customer/merchant
	CreateTransactionRequest struct {
		CustomerID       uint64
		PaymentChannelID uint64
		Amount           int64
		TypeTransaction  TransactionType
	}

	// CreateTransactionResponse consist data response of creating transactions customer/merchant
	CreateTransactionResponse struct {
		ID     uint64
		Status string
	}
)

const (
	TransactionPending TransactionStatus = "PENDING"
	//...
)
