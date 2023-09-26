package model

// TransactionType consist data type of transactions (topup/payout)
type TransactionType string

const (
	TopupTransaction  TransactionType = "TOPUP"
	PayoutTransaction TransactionType = "PAYOUT"
)
