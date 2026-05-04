package models

import "time"

// UserNode represents a user in the system, containing essential information 
// for fraud detection and risk assessment.
// In a graph database, it serves as a central node connected to various claims, 
// tickets, and account history nodes.
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
	// such as ticket status and contact date, which are crucial for understanding the context 
	// of the user's issue.
	UserTicketInfo []*TicketNodeInfo `json:"user_ticket_info" cypher:"-"`
}


// Claim represents the assertions made by the user regarding their account and registration details.
// In a graph database, it is a node connected to the User node, providing insights into the 
// user's claims about their account, which can be used for fraud detection and risk assessment.
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

// AccountInfo represents the user's assertions about their financial standing and 
// identity verification status.
// In a graph database, it is a node connected to the Claim node, providing insights 
// into the user's claims about their account's financial and verification status, 
// which can be used for fraud detection and risk assessment.
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

// TicketNode represents the support ticket information related to 
// the user's claims and interactions with customer support.
// In a graph database, it is a node connected to the User node, 
// providing insights into the user's interactions with support and the status of their issues,
//  which can be used for fraud detection and risk assessment.
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

// Session represents a continuous authenticated period of user activity.
// In a graph, comparing Session nodes against Transaction nodes helps
// identify proxy usage or account takeovers (ATO).
type SessionNode struct {
	// SessionID is the unique token identifying this specific authenticated period.
	SessionID string `json:"session_id" cypher:"id"`

	// SessionIP is the network address captured when the session was established.
	SessionIP string `json:"session_ip" cypher:"ip"`

	// DeviceID is the hardware identifier. Used to create a relationship with a Device node.
	DeviceID string `json:"device_id" cypher:"device_id"`

	// ASN is the Autonomous System Number (ISP) associated with the session's origin.
	ASN string `json:"asn" cypher:"asn"`

	// Country is the geographic origin of the session establishment.
	Country string `json:"country" cypher:"country"`

	// City is the specific city where the session originated.
	City string `json:"city" cypher:"city"`

	// StartTime is the RFC3339 timestamp of the login or session creation.
	StartTime time.Time `json:"start_time" cypher:"start_time"`

	// EndTime is the RFC3339 timestamp of logout or session expiration.
	EndTime time.Time `json:"end_time" cypher:"end_time"`
}


// SupportContext represents the context of the user's support ticket, including the reason for contact and the communication channel used.
// In a graph database, it is a node connected to the Ticket node, providing insights into the user's interactions with support and the context of their issues, which can be used for fraud detection and risk assessment.
type SupportContextNodeInfo struct {
	// Reason is the reason for contacting support (e.g., "Account Recovery").
	Reason string `json:"reason" cypher:"reason"`

	// OperatorID is the identifier of the support agent handling the ticket, which can be used to analyze patterns in support interactions.
	OperatorID string `json:"operator_id" cypher:"operator_id"`

	// Source is the communication channel used (e.g., "chat", "email", "phone"), which can provide insights into user behavior and preferences.
	Source string `json:"source" cypher:"source"`
}