// Package errors defines the shared sentinel errors and status constants 
// used across the scoring engine and its calculators.
package errors

import "errors"

// --- Email Scoring Errors & Statuses ---

var (
	// ErrEmptyEmailData is returned when either the claim or the DB record lacks an email address.
	ErrEmptyEmailData = errors.New("CORE_SCORE_EMPTY_EMAIL")
	// ErrEmailNoMatch indicates a complete mismatch between provided email addresses.
	ErrEmailNoMatch   = errors.New("CORE_SCORE_EMAIL_MISMATCH")

	// Standard status codes for basic matching logic.
	StatusNoData  = "no_data"
	StatusMatch   = "match"
	StatusNoMatch = "no_match"
)

// --- Device Scoring Errors ---

var (
	// ErrEmptyDeviceData is returned for single device comparison failures.
	ErrEmptyDeviceData = errors.New("CORE_SCORE_EMPTY_DEVICE")
	// ErrDeviceNoMatch indicates that the specific device ID does not match the records.
	ErrDeviceNoMatch   = errors.New("CORE_SCORE_DEVICE_MISMATCH")
)

// --- Device List & History Errors ---

var (
	// ErrEmptyDeviceList occurs when a user's device history is completely empty.
	ErrEmptyDeviceList = errors.New("CORE_SCORE_DEVICE_LIST_EMPTY")
	// ErrFutureRegDate is a logical error where a registration date is set in the future.
	ErrFutureRegDate   = errors.New("CORE_SCORE_REGDATE_IN_FUTURE")
)

// --- Phone Scoring Errors ---

var (
	// ErrEmptyPhoneData is returned when phone fields are missing.
	ErrEmptyPhoneData = errors.New("CORE_SCORE_EMPTY_PHONE")
	// ErrPhoneTooShort indicates the phone number does not meet minimum length requirements.
	ErrPhoneTooShort  = errors.New("CORE_SCORE_PHONE_TOO_SHORT")
	// ErrPhoneNoMatch indicates no full or partial (suffix) match was found.
	ErrPhoneNoMatch   = errors.New("CORE_SCORE_PHONE_MISMATCH")
)

// StatusPartial is used when data matches partially (e.g., phone suffix or fuzzy match).
const StatusPartial = "partial"

// --- Registration Date Errors ---

var (
	// ErrEmptyRegDate occurs when registration timestamps are missing.
	ErrEmptyRegDate    = errors.New("CORE_SCORE_EMPTY_REG_DATE")
	// ErrRegDateAnomaly indicates a significant time discrepancy (e.g., > 1 year).
	ErrRegDateAnomaly  = errors.New("CORE_SCORE_REG_DATE_ANOMALY")
)

// StatusAnomaly indicates that the result is technically valid but suspicious.
const StatusAnomaly = "anomaly"

// --- Location (City/Country) Errors ---

var (
	// ErrEmptyLocationData occurs when city or country fields are blank.
	ErrEmptyLocationData = errors.New("CORE_SCORE_EMPTY_LOCATION")
	// ErrLocationMismatch indicates a geographic discrepancy.
	ErrLocationMismatch  = errors.New("CORE_SCORE_LOCATION_MISMATCH")
)

// --- Transaction & Fraud Detection Errors ---

var (
	// ErrTxTimeParse indicates a malformed timestamp string in transaction data.
	ErrTxTimeParse      = errors.New("TX_TIME_PARSE_ERROR")
	// ErrHighFreqTx is triggered when multiple transactions occur in a short window (velocity).
	ErrHighFreqTx       = errors.New("TX_HIGH_FREQUENCY_DETECTED")
	// ErrSuddenHighAmount is triggered by amounts exceeding historical averages significantly.
	ErrSuddenHighAmount = errors.New("TX_SUDDEN_HIGH_AMOUNT")
)

// Additional workflow statuses.
var (
	// StatusAnomalyBlock is used to immediately reject a request due to multiple red flags.
	StatusAnomalyBlock = "anomaly_block"
	// StatusSkipped is used when a calculation is bypassed (e.g., weight is zero).
	StatusSkipped      = "skipped"
)