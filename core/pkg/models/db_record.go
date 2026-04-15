// Package models contains the data structures used throughout the scoring engine.
package models

import "time"

// DBRecord represents the historical state of a user account retrieved from the database.
// It serves as the "Gold Standard" or baseline for comparing new claims and 
// calculating consistency scores.
type DBRecord struct {
	// AccTag is a unique internal identifier or label for the account.
	AccTag           string      `json:"acc_tag"`

	// RegCountry is the ISO country code provided during account registration.
	RegCountry       string      `json:"reg_country"`

	// RegCity is the city name recorded at the time of account creation.
	RegCity          string      `json:"reg_city"`

	// FirstEmail is the primary email address linked to the account's history.
	FirstEmail       string      `json:"first_email"`

	// Phone is the verified phone number associated with the user profile.
	Phone            string      `json:"phone"`

	// FirstDevice is the hardware identifier (UUID/IMEI) of the original registration device.
	FirstDevice      string      `json:"first_device"`

	// Devices is a historical list of all unique device fingerprints authorized by the user.
	Devices          []string    `json:"devices"`

	// IsDonator indicates premium status, used to determine the logic/weighting for financial rules.
	IsDonator        bool        `json:"is_donator"`

	// FirstTransaction contains metadata regarding the user's initial financial activity.
	FirstTransaction Transaction `json:"first_transaction"`

	// UserHistory encapsulates behavioral metrics and event logs for the account's lifecycle.
	UserHistory      UserHistory `json:"user_history"`

	// RegDate is the timestamp when the account was officially created.
	RegDate          time.Time   `json:"reg_date"`
}