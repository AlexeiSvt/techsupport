// Package models contains the data structures used throughout the scoring engine.
package models

import (
	"time"
)

// SupportContext represents the environmental and technical metadata captured
// at the moment a user contacts support. It helps verify if the person
// submitting a claim is using a known or suspicious environment.
type SupportContext struct {
	// TicketID is the unique identifier for the support request or chat session.
	TicketID string `json:"ticket_id" cypher:"id"`

	// UBTicketID is the internal Bot-Session identifier (the "UB" nuance).
	UBTicketID string `json:"ub_ticket_id" cypher:"ub_ticket_id"`

	// Country is the geographic origin of the support request.
	Country string `json:"country" cypher:"country"`

	// City is the specific city from which the user is contacting support.
	City string `json:"city" cypher:"city"`

	// DeviceID is the hardware identifier of the device used to contact support.
	DeviceID string `json:"device_id" cypher:"device_id"`

	// DeviceName is the human-readable model name (e.g., "Pixel 8 Pro").
	DeviceName string `json:"device_name" cypher:"device_name"`
}

// UserClaim encapsulates the identity attributes and historical markers
// asserted by the user. These are subjective values provided by the person
// attempting to prove ownership of an account.
type UserClaim struct {
	// ClaimID is a unique identifier for this specific set of assertions.
	ClaimID string `json:"claim_id" cypher:"id"`

	// We link this to the Bot session ID.
	UBTicketID string `json:"ub_ticket_id" cypher:"ub_ticket_id"`

	// AccTag is the identifier for the account the user is attempting to recover.
	AccTag string `json:"acc_tag" cypher:"acc_tag"`

	// RegCountry is the original registration country claimed by the user.
	RegCountry string `json:"reg_country" cypher:"reg_country"`

	// RegCity is the original registration city claimed by the user.
	RegCity string `json:"reg_city" cypher:"reg_city"`

	// Phone is the phone number the user claims is linked to the account.
	Phone string `json:"phone" cypher:"phone"`

	// FirstDeviceID is the ID of the hardware used for account creation (claimed).
	FirstDeviceID string `json:"first_device_id" cypher:"first_device_id"`
	// RegDate is the estimated account creation timestamp provided by the user.
	RegDate time.Time `json:"reg_date" cypher:"reg_date"`
}
