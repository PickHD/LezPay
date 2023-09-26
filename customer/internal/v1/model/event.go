package model

// EventName consist data event name in proto TopupTransactionMessage
type EventName string

const (
	EventPaymentRequested EventName = "PaymentRequested"
)
