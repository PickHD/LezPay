package model

type (
	// PaymentChannelStatus consist type of payment channel status
	PaymentChannelStatus string

	// PaymentChannel representing of table payment_channel
	PaymentChannel struct {
		ID                 uint64               `json:"id"`
		Name               string               `json:"name"`
		Code               string               `json:"code"`
		ImageURL           string               `json:"image_url"`
		PaymentInstruction string               `json:"payment_instruction"`
		Status             PaymentChannelStatus `json:"status"`
	}
)

const (
	PaymentChannelStatusActive   PaymentChannelStatus = "ACTIVE"
	PaymentChannelStatusInactive PaymentChannelStatus = "INACTIVE"
)
