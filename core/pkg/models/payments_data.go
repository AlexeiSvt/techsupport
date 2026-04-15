// Package models contains the data structures used throughout the scoring engine.
package models

// Transaction represents a financial or significant event within the system.
// It captures a snapshot of the environmental and technical metadata at the 
// moment of the operation, used to calculate risk and velocity metrics.
type Transaction struct {
	// TransactionID is the unique reference number provided by the payment gateway or internal system.
	TransactionID string `json:"transaction_id"`

	// Amount is the numerical value of the transaction.
	Amount float64 `json:"amount"`

	// Currency is the ISO 4217 code for the currency (e.g., "USD", "EUR").
	Currency string `json:"currency"`

	// PaymentMethod represents the type of payment (e.g., "credit_card", "e-wallet").
	PaymentMethod string `json:"payment_method"`

	// PaymentProvider is the third-party service handling the payment (e.g., "Stripe", "PayPal").
	PaymentProvider string `json:"payment_provider"`

	// DeviceID is the unique hardware fingerprint at the time of the transaction.
	DeviceID string `json:"device_id"`

	// DeviceModel is the human-readable description of the hardware (e.g., "iPhone 15 Pro").
	DeviceModel string `json:"device_model"`

	// IP is the network address from which the transaction was initiated.
	IP string `json:"ip"`

	// Country is the geographic origin of the transaction IP.
	Country string `json:"country"`

	// City is the specific city associated with the transaction IP.
	City string `json:"city"`

	// Timestamp is the RFC3339 formatted string indicating when the transaction occurred.
	Timestamp string `json:"timestamp"`

	// ASN (Autonomous System Number) identifies the Internet Service Provider (ISP) used.
	ASN string `json:"asn"` 

	// SessionContext provides the broader navigational context in which this transaction took place.
	SessionContext Session `json:"session_context"`
}

// Session represents a continuous period of user activity on the platform.
// Comparing Session metadata against Transaction metadata is a key technique 
// for detecting session hijacking or proxy usage.
type Session struct {
	// SessionID is the unique token identifying the current authenticated period.
	SessionID string `json:"session_id"`

	// SessionIP is the IP address used when the session was first established.
	SessionIP string `json:"session_ip"`

	// DeviceID is the hardware identifier used for the session login.
	DeviceID string `json:"device_id"`

	// ASN is the ISP identifier associated with the session's origin.
	ASN string `json:"asn"`

	// Country is the geographic origin of the session.
	Country string `json:"country"`

	// City is the specific city of the session's origin.
	City string `json:"city"`

	// StartTime is when the user logged in or the session began.
	StartTime string `json:"start_time"`

	// EndTime is when the session expired or the user logged out.
	EndTime string `json:"end_time"`
}