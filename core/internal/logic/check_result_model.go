// Package logic provides calculation engines for scoring user data.
package logic

// rawCheckResult represents the intermediate output of a single validation rule.
// It carries the numerical score, a status code for flow control, and a 
// human-readable audit trail.
type rawCheckResult struct {
	// Value is the calculated score or multiplier (e.g., 0.0 to 1.0).
	Value   float64
	// Status indicates the outcome of the check (e.g., "match", "partial", "anomaly").
	Status  string
	// Comment provides a detailed explanation for why this specific score was assigned.
	Comment string
}

// RawTxResult is an alias for rawCheckResult, specifically used in transaction validation.
// Using an alias maintains semantic clarity when dealing with financial logic 
// while reusing the underlying data structure.
type RawTxResult = rawCheckResult