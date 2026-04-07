package models

import "time"

type InputData struct {
	UserData UserData `json:"user_data"`
	DBRecord DBRecord `json:"db_record"`
}

type UserData struct {
	IP         string    `json:"ip"`
	ASN        string    `json:"asn"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	UserClaim  UserClaim `json:"user_claim"`
}

type OutputData struct {
	TicketID       string    `json:"ticket_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FinaPercentage string    `json:"final_percentage"`
	Metrics        Metrics   `json:"metrics"`
}

type UserClaim struct {
	AccTag           string      `json:"acc_tag"`
	RegCountry       string      `json:"reg_country"`
	RegCity          string      `json:"reg_city"`
	FirstEmail       string      `json:"first_email"`
	Phone            string      `json:"phone"`
	FirstDevice      string      `json:"first_device"`
	Devices          []string    `json:"devices"`
	FirstTransaction Transaction `json:"first_transaction"`
	RegDate          time.Time   `json:"reg_date"`
}

type DBRecord struct {
	AccTag           string      `json:"acc_tag"`
	RegCountry       string      `json:"reg_country"`
	RegCity          string      `json:"reg_city"`
	FirstEmail       string      `json:"first_email"`
	Phone            string      `json:"phone"`
	FirstDevice      string      `json:"first_device"`
	Devices          []string    `json:"devices"`
	IsDonator        bool        `json:"is_donator"`
	FirstTransaction Transaction `json:"first_transaction"`
	UserHistory      UserHistory `json:"user_history"`
	RegDate          time.Time   `json:"reg_date"`
}

type Metrics struct {
	Knowledge    float64 `json:"knowledge"`
	PenaltyScore float64 `json:"penalty_score"`
}
