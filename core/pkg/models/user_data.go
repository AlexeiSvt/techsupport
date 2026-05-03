// // Package models contains the data structures used throughout the scoring engine.
// package models

// import (
// 	"time"
// )

// // SupportContext represents the environmental and technical metadata captured
// // at the moment a user contacts support.
// type SupportContext struct {
// 	// TicketID is the unique identifier for the support request or chat session.
// 	TicketID string `json:"ticket_id" cypher:"id"`

// 	// UBTicketID is the internal Bot-Session identifier (the "UB" nuance).
// 	UBTicketID string `json:"ub_ticket_id" cypher:"ub_ticket_id"`

// 	// Country is the geographic origin of the support request.
// 	Country string `json:"country" cypher:"country"`

// 	// City is the specific city from which the user is contacting support.
// 	City string `json:"city" cypher:"city"`

// 	// DeviceID is the hardware identifier of the device used to contact support.
// 	DeviceID string `json:"device_id" cypher:"device_id"`

// 	// DeviceName is the human-readable model name (e.g., "Pixel 8 Pro").
// 	DeviceName string `json:"device_name" cypher:"device_name"`

// 	// History contains pre-fetched historical data from Neo4j for comparison.
// 	// This field is excluded from Cypher mapping as it represents relationships.
// 	History UserHistoryContext `json:"-" cypher:"-"`
// }

// // UserClaim encapsulates the identity attributes and historical markers asserted by the user.
// type UserClaim struct {
// 	// ClaimID is a unique identifier for this specific set of assertions.
// 	ClaimID string `json:"claim_id" cypher:"id"`

// 	// We link this to the Bot session ID.
// 	UBTicketID string `json:"ub_ticket_id" cypher:"ub_ticket_id"`

// 	// AccTag is the identifier for the account the user is attempting to recover.
// 	AccTag string `json:"acc_tag" cypher:"acc_tag"`

// 	// RegCountry is the original registration country claimed by the user.
// 	RegCountry string `json:"reg_country" cypher:"reg_country"`

// 	// RegCity is the original registration city claimed by the user.
// 	RegCity string `json:"reg_city" cypher:"reg_city"`

// 	// Phone is the phone number the user claims is linked to the account.
// 	FirstPhone string `json:"phone" cypher:"phone"`

// 	// Email is the email the user claims is linked to the account.
// 	FirstEmail string `json:"email" cypher:"email"`

// 	// FirstDeviceName is the ID of the hardware used for account creation (claimed).
// 	FirstDeviceName string `json:"first_device_id" cypher:"first_device_id"`

// 	// RegDate is the estimated account creation timestamp provided by the user.
// 	RegDate time.Time `json:"reg_date" cypher:"reg_date"`

// 	// FirstTransaction is the specific payment detail the user is using to prove ownership.
// 	FirstTransaction Transaction `json:"first_transaction" cypher:"-"`
// }

// // UserHistoryContext provides flattened historical data windows for validators.
// type UserHistoryContext struct {
// 	FirstWindow []Session `json:"first_window"`
// 	LastWindow  []Session `json:"last_window"`
// 	AllPayments []Payment `json:"all_payments"`
// }

package models

import "time"

type UserNode struct {
	// UserID is the unique identifier for the user in company's system.
	UserID string `json:"user_id" cypher:"id"`

	//When the user's account, who has written about the issue, 
	// was created in the company's system. 
	// This is a critical attribute for fraud detection, as it helps 
	// establish the age of the account 
	// and can be used to identify patterns of fraudulent behavior.
	CreationDate time.Time `json:"creation_date" cypher:"creation_date"`

	//Refers to the Claim node that shows what the user claims about their account, 
	// such as registration details and first transaction information.
	UserClaim []*ClaimNode `json:"user_claims" cypher:"-"`
}

type ClaimNode struct {
	// ClaimID is a unique identifier for this specific set of assertions made by the user.
	ClaimID string `json:"claim_id" cypher:"id"`

	//Field that the user claims about his registration country.
	RegistrationCountry string `json:"registration_country" cypher:"registration_country"`

	//Field that the user claims about his registration city.
	RegistrationCity string `json:"registration_city" cypher:"registration_city"`

	//Field that the user claims about the first phone number linked to the account.
	FirstPhone string `json:"first_phone" cypher:"first_phone"`

	//Field that the user claims about the first email linked to the account.
	FirstEmail string `json:"first_email" cypher:"first_email"`

	//Field that the user claims about the first device name linked to the account.
	FirstDeviceName string `json:"first_device_name" cypher:"first_device_name"`

	//Field that the user claims about the registration date of the account.
	RegistrationDate time.Time `json:"registration_date" cypher:"registration_date"`

	AccountInfo []*AccountInfoNode `json:"account_info" cypher:"-"`
}

type AccountInfoNode struct {
	// AccountInfoID is a unique identifier for this specific account information entry.
	AccountInfoID string `json:"account_info_id" cypher:"id"`

	//Field that the user claims about the first transaction linked to the account.
	FirstTransaction Transaction `json:"first_transaction" cypher:"-"`
}