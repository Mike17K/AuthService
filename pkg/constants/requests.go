package constants

// Request Context Keys
type ContextKey string

const (
	// ApplicationContextKey is used to store the application object in the request context
	ApplicationContextKey ContextKey = "application"
	// UserContextKey is used to store the user object in the request context
	UserContextKey ContextKey = "user"
)

// Request Headers
const (
	// ServiceToServiceAuthorizationHeader is for Auth service to Application For Private Routes
	ServiceToServiceAuthorizationHeader = "Auth-Service-Authorization"
	// ApplicationAuthorizationHeader is for validating application identity the request is from
	ApplicationAuthorizationHeader = "Application-Secret"
	// AuthorizationHeader is for validating user identity the request is from
	Authorization = "Authorization"
)
