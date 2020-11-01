package middleware

// ContextKey is type for various request context keys
type ContextKey string

const (
	// ParsedRequest is key used to identify parsed request from request context
	ParsedRequest = ContextKey("parsedRequest")
)
