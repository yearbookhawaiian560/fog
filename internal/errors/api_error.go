package errors

import "encoding/json"

type APIError struct {
	ErrorCode int `json:"error_code"`

	// Only present in case of internal errors (status 500 or 502)
	ErrorID string `json:"error_id,omitempty"`

	// Human-friendly description that can be shown to the user.
	Detail string `json:"detail"`

	Data json.RawMessage `json:"data"`

	// The original error that caused this response, not forwarded to the end user.
	// Can be nil
	RawError error `json:"-"`
	// HTTP status that the error should have.
	Status int `json:"-"`
}
