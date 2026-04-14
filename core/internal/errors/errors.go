package errors

import "errors"

var (
    ErrEmptyEmailData = errors.New("CORE_SCORE_EMPTY_EMAIL")
    ErrEmailNoMatch   = errors.New("CORE_SCORE_EMAIL_MISMATCH")

    StatusNoData  = "no_data"
    StatusMatch   = "match"
    StatusNoMatch = "no_match"
)

var (
	ErrEmptyDeviceData = errors.New("CORE_SCORE_EMPTY_DEVICE")
	ErrDeviceNoMatch   = errors.New("CORE_SCORE_DEVICE_MISMATCH")
)

var (
    ErrEmptyDeviceList = errors.New("CORE_SCORE_DEVICE_LIST_EMPTY")
    ErrFutureRegDate   = errors.New("CORE_SCORE_REGDATE_IN_FUTURE")
)

var (
	ErrEmptyPhoneData = errors.New("CORE_SCORE_EMPTY_PHONE")
	ErrPhoneTooShort  = errors.New("CORE_SCORE_PHONE_TOO_SHORT")
	ErrPhoneNoMatch   = errors.New("CORE_SCORE_PHONE_MISMATCH")
)

const StatusPartial = "partial"

var (
	ErrEmptyRegDate    = errors.New("CORE_SCORE_EMPTY_REG_DATE")
	ErrRegDateAnomaly  = errors.New("CORE_SCORE_REG_DATE_ANOMALY")
)

const StatusAnomaly = "anomaly"

var (
    ErrEmptyLocationData = errors.New("CORE_SCORE_EMPTY_LOCATION")
    ErrLocationMismatch  = errors.New("CORE_SCORE_LOCATION_MISMATCH")
)

var (
	ErrTxTimeParse      = errors.New("TX_TIME_PARSE_ERROR")
	ErrHighFreqTx       = errors.New("TX_HIGH_FREQUENCY_DETECTED")
	ErrSuddenHighAmount = errors.New("TX_SUDDEN_HIGH_AMOUNT")
)

var (
    StatusAnomalyBlock = "anomaly_block"
    StatusSkipped      = "skipped"
)