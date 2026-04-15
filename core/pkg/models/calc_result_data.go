// Package models contains the data structures used throughout the scoring engine.
package models

// CalcResult represents the detailed output of a single scoring rule or calculator.
// It serves as an audit trail, providing transparency into how a specific 
// component contributed to the overall score.
type CalcResult struct {
	// Name is the human-readable identifier for the rule (e.g., "Email Consistency").
	Name    string  `json:"name"`

	// Code is a machine-readable unique identifier for the rule (e.g., "email_match").
	Code    string  `json:"code"`

	// Value is the raw numerical result or score before weighting is applied.
	Value   float64 `json:"value"`

	// Weight is the importance factor assigned to this rule (e.g., 0.5 for 50%).
	Weight  float64 `json:"weight"`

	// Result is the final weighted score (Value * Weight) contributed to the total.
	Result  float64 `json:"result"`

	// Comment provides a natural-language explanation for the result, 
	// essential for auditing and troubleshooting.
	Comment string  `json:"comment"`

	// Status represents the categorical outcome (e.g., "match", "partial", "anomaly").
	Status  string  `json:"status"`
}