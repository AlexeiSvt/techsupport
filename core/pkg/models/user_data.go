package models

import (
	"techsupport/core/internal/ipchecker"
	"time"
)

type UserData struct {
	IP         string    `json:"ip"`
	ASN        string    `json:"asn"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	UserClaim  UserClaim `json:"user_claim"`
}

type UserClaim struct {
	AccTag           string                   `json:"acc_tag"`
	RegCountry       string                   `json:"reg_country"`
	RegCity          string                   `json:"reg_city"`
	FirstEmail       string                   `json:"first_email"`
	Phone            string                   `json:"phone"`
	FirstDevice      string                   `json:"first_device"`
	Devices          []string                 `json:"devices"`
	IPInfo           *ipchecker.IpApiResponse `json:"ip_info"`
	FirstTransaction Transaction              `json:"first_transaction"`
	RegDate          time.Time                `json:"reg_date"`
}


type UserHistory struct {
	FirstWindow []Session `json:"first_300_times"`
	LastWindow []Session `json:"last_300_times"`
	AllPayments []Transaction `json:"all_payments"`
}