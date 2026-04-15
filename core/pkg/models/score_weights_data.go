// Package models contains the data structures used throughout the scoring engine.
package models

// Weights defines the influence or importance of each specific scoring rule.
// These floating-point multipliers (typically between 0.0 and 1.0) allow the engine
// to be tuned for different scenarios (e.g., stricter checks for high-value accounts).
type Weights struct {
	// RegDate is the weight applied to the account age and registration date consistency.
	RegDate          float64 `json:"reg_date"`

	// RegCountry is the weight for verifying the origin country against historical records.
	RegCountry       float64 `json:"reg_country"`

	// RegCity is the weight for verifying the specific city against historical records.
	RegCity          float64 `json:"reg_city"`

	// FirstEmail is the weight assigned to the match between the current and original email.
	FirstEmail       float64 `json:"first_email"`

	// Phone is the weight for phone number verification and consistency.
	Phone            float64 `json:"phone"`

	// FirstDevice is the weight for matching the current device to the registration hardware.
	FirstDevice      float64 `json:"first_device"`

	// Devices is the weight applied to the overall device history and fingerprint analysis.
	Devices          float64 `json:"devices"`

	// FirstTransaction is the weight specifically for financial activity verification.
	// It is omitted if the user has no transaction history (e.g., non-donators).
	FirstTransaction float64 `json:"first_transaction,omitempty"`
}