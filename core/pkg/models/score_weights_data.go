package models

type Weights struct {
	RegDate          float64 `json:"reg_date"`
	RegCountry       float64 `json:"reg_country"`
	RegCity          float64 `json:"reg_city"`
	FirstEmail       float64 `json:"first_email"`
	Phone            float64 `json:"phone"`
	FirstDevice      float64 `json:"first_device"`
	Devices          float64 `json:"devices"`
	FirstTransaction float64 `json:"first_transaction,omitempty"`
}
