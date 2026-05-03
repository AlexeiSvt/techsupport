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
	UserClaimInfo []*ClaimInfoNode `json:"user_claims" cypher:"-"`

	//Refers to the Ticket node that contains information about the support ticket,
	// such as ticket status and contact date, which are crucial for understanding the context of the user's issue.
	UserTicketInfo []*TicketNodeInfo `json:"user_ticket_info" cypher:"-"`
}


type ClaimInfoNode struct {
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
	FirstDeviceModel string `json:"first_device_name" cypher:"first_device_name"`

	//Field that the user claims about the registration date of the account.
	RegistrationDate time.Time `json:"registration_date" cypher:"registration_date"`

	AccountInfo []*AccountInfoNode `json:"account_info" cypher:"-"`
}

type AccountInfoNode struct {
	// AccountInfoID is a unique identifier for this specific account information entry.
	AccountInfoID string `json:"account_info_id" cypher:"id"`

	//IsSolvent indicates whether the user claims to have a positive financial standing, which can be a factor in assessing the legitimacy of the account.
	IsSolvent bool `json:"is_solvent" cypher:"is_solvent"`

	//IsVerified indicates whether the user claims to have completed identity verification processes, which can be a strong indicator of account legitimacy.
	IsVerified bool `json:"is_verified" cypher:"is_verified"`

	//AccountHistory contains historical data about the user's account, such as past sessions and payments, which can be used to verify the user's claims and assess the risk of fraud.
	AccountHistory *[]AccountHistoryNodeInfo `json:"account_history" cypher:"-"`
}

type TicketNodeInfo struct {
	// TicketID is the unique identifier for the support ticket.
	TicketID string `json:"ticket_id" cypher:"id"`

	//TicketStatus indicates the current state of the support ticket (e.g., "open", "pending", "resolved").
	TicketStatus string `json:"ticket_status" cypher:"ticket_status"`

	// ContactDate is the timestamp when the user contacted support, crucial for temporal analysis of user behavior.
	ContactDate time.Time `json:"contact_date" cypher:"contact_date"`

	//UpdatedDate is the timestamp when the ticket was last updated, which can indicate ongoing issues or resolution status.
	UpdatedDate time.Time `json:"updated_date" cypher:"updated_date"`
}

type AccountHistoryNodeInfo struct {
	// FistWindowSessions contains the sessions from the first historical window (e.g., 30 days after account creation).
	FirstWindowSessions []*SessionNode `json:"first_window_sessions" cypher:"-"`

	// LastWindowSessions contains the sessions from the last historical window (e.g., 30 days before the support contact).
	LastWindowSessions []*SessionNode `json:"last_window_sessions" cypher:"-"`

	// AllPayments contains all the payment transactions linked to the account, used for ownership proof and fraud detection.
	AllPayments []*PaymentInfoNode `json:"all_payments" cypher:"-"`
}
