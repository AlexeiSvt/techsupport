package iperrors

import "errors"

var (
	ErrIPEmpty         = errors.New("IP_EMPTY")
	ErrIPInvalidFormat = errors.New("IP_INVALID_FORMAT")
	ErrApiKeyMissing   = errors.New("IP_API_KEY_MISSING")
	ErrApiUrlInvalid   = errors.New("IP_API_URL_INVALID")
	ErrApiRequest      = errors.New("IP_API_REQUEST_FAILED")
	ErrApiStatus       = errors.New("IP_API_BAD_STATUS")
	ErrApiDecode       = errors.New("IP_API_DECODE_ERROR")
	ErrApiResponseEmpty = errors.New("IP_API_RESPONSE_EMPTY")
)