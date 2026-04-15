package models

import "time"

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