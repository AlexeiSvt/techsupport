package models

type Transaction struct {
	TransactionID string `json:"transaction_id"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	PaymentProvider string `json:"payment_provider"`
	DeviceID string `json:"device_id"`
	DeviceModel string `json:"device_model"`
	IP string `json:"ip"`
	Country string `json:"country"`
	City string `json:"city"`
	Timestamp string `json:"timestamp"`
	ASN string `json:"asn"`	
	SessionContext Session `json:"session_context"`
}

type Session struct {
	SessionID string `json:"session_id"`
	SessionIP string `json:"session_ip"`
	DeviceID string `json:"device_id"`
	ASN string `json:"asn"`
	Country string `json:"country"`
	City string `json:"city"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}