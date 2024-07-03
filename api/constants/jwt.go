package constants

// JWTTokenType defines the type for JWT token types
type JWTTokenType string

// Predefined token types
const (
	AccessToken  JWTTokenType = "access_token"
	RefreshToken JWTTokenType = "refresh_token"
)

// IsValid checks if the value is one of the predefined token types
func (t JWTTokenType) IsValid() bool {
	switch t {
	case AccessToken, RefreshToken:
		return true
	}
	return false
}

const JWTUserIdField string = "user_id"
const JWTUserTypeField string = "user_type"
const JWTTokenTypeField string = "token_type"
