// Package models contains the data structures used throughout the scoring engine.
package models

import (
	"techsupport/core/internal/ipchecker"
	"time"
)

// UserData represents the environmental and technical metadata captured 
// during the current user session.
type UserData struct {
	// IP is the current network address of the user.
	IP         string    `json:"ip"`
	// ASN is the Autonomous System Number (ISP) associated with the current IP.
	ASN        string    `json:"asn"`
	// City is the geographic city detected for the current session.
	City       string    `json:"city"`
	// Country is the geographic country code detected for the current session.
	Country    string    `json:"country"`
	// DeviceID is the unique hardware identifier of the current device.
	DeviceID   string    `json:"device_id"`
	// DeviceName is the human-readable name of the current device (e.g., "Android 14").
	DeviceName string    `json:"device_name"`
	// UserClaim contains the identity and historical assertions provided by the user.
	UserClaim  UserClaim `json:"user_claim"`
}

// UserClaim encapsulates the specific identity attributes and historical 
// markers that the user is currently asserting.
type UserClaim struct {
	// AccTag is the internal identifier for the account being claimed.
	AccTag           string                   `json:"acc_tag"`
	// RegCountry is the original country of registration claimed by the user.
	RegCountry       string                   `json:"reg_country"`
	// RegCity is the original city of registration claimed by the user.
	RegCity          string                   `json:"reg_city"`
	// FirstEmail is the primary email address linked to the account.
	FirstEmail       string                   `json:"first_email"`
	// Phone is the current phone number associated with the account.
	Phone            string                   `json:"phone"`
	// FirstDevice is the original device ID used during registration.
	FirstDevice      string                   `json:"first_device"`
	// Devices is a list of known hardware fingerprints associated with this user.
	Devices          []string                 `json:"devices"`
	// IPInfo contains detailed third-party intelligence regarding the current IP address.
	IPInfo           *ipchecker.IpApiResponse `json:"ip_info"`
	// FirstTransaction holds data about the very first financial event of this user.
	FirstTransaction Transaction              `json:"first_transaction"`
	// RegDate is the timestamp of the account's creation.
	RegDate          time.Time                `json:"reg_date"`
}

// UserHistory provides a longitudinal view of user behavior by grouping 
// sessions into specific time-based windows.
type UserHistory struct {
	// FirstWindow represents the first 300 sessions, establishing the "baseline" behavior.
	FirstWindow []Session `json:"first_300_times"`
	// LastWindow represents the most recent 300 sessions, used to detect sudden shifts in behavior.
	LastWindow []Session `json:"last_300_times"`
	// AllPayments is a chronological record of every financial transaction made by the user.
	AllPayments []Transaction `json:"all_payments"`
}