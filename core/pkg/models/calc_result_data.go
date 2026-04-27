// Package models contains the data structures used throughout the scoring engine.
package models

// CalcResult represents the output of a single scoring rule or logic unit.
// It provides transparency, allowing developers and auditors to see exactly 
// how a specific rule influenced the final verdict.
type CalcResult struct {
    // Name is the descriptive name of the rule (e.g., "IP Geolocation Match").
    Name string `json:"name" cypher:"name"`

    // Code is the machine-readable constant identifying the rule (e.g., "geo_ip_match").
    Code string `json:"code" cypher:"code"`

    // Value is the raw score produced by the logic before weight adjustment.
    Value float64 `json:"value" cypher:"value"`

    // Weight represents the relative importance of this rule in the total score.
    Weight float64 `json:"weight" cypher:"weight"`

    // Result is the final contribution to the total score (Value * Weight).
    Result float64 `json:"result" cypher:"result"`

    // Comment explains the reasoning behind the score (e.g., "City mismatch: Moscow vs Tula").
    Comment string `json:"comment" cypher:"comment"`

    // Status categorizes the outcome for quick filtering (e.g., "MATCH", "MISMATCH", "FRAUD").
    Status string `json:"status" cypher:"status"`
}