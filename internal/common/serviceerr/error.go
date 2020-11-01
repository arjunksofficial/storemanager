package serviceerr

// APIError is model for custom error with API handlers
type APIError struct {
	StatusCode int
	Error      error
}
