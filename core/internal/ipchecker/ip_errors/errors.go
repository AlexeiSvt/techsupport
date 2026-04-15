// Package iperrors defines standard error constants for the IP verification process.
// These sentinel errors allow for precise error handling and auditing throughout the system.
package iperrors

import "errors"

var (
	// ErrIPEmpty is returned when the input IP string is empty or contains only whitespace.
	ErrIPEmpty         = errors.New("IP_EMPTY")

	// ErrIPInvalidFormat is returned when the string is not a valid IPv4 or IPv6 address.
	ErrIPInvalidFormat = errors.New("IP_INVALID_FORMAT")

	// ErrApiKeyMissing indicates that the authentication token for the external API is not configured.
	ErrApiKeyMissing   = errors.New("IP_API_KEY_MISSING")

	// ErrApiUrlInvalid occurs if the base URL or constructed request URL is malformed.
	ErrApiUrlInvalid   = errors.New("IP_API_URL_INVALID")

	// ErrApiRequest occurs when the network call fails (e.g., DNS issues, timeouts, or connection refused).
	ErrApiRequest      = errors.New("IP_API_REQUEST_FAILED")

	// ErrApiStatus is returned when the API responds with a non-200 HTTP status code.
	ErrApiStatus       = errors.New("IP_API_BAD_STATUS")

	// ErrApiDecode is returned when the API response body is not valid JSON or doesn't match the expected struct.
	ErrApiDecode       = errors.New("IP_API_DECODE_ERROR")

	// ErrApiResponseEmpty indicates that the API returned a successful response but with no meaningful data.
	ErrApiResponseEmpty = errors.New("IP_API_RESPONSE_EMPTY")
)