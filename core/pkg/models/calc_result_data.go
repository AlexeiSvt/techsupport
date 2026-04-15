package models

type CalcResult struct {
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	Value   float64 `json:"value"`
	Weight  float64 `json:"weight"`
	Result  float64 `json:"result"`
	Comment string  `json:"comment"`
	Status  string  `json:"status"`
}
