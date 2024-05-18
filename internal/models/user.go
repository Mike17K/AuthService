package models

// Provider is an enum for authentication providers
type Provider int

const (
	Email Provider = iota
	Google
	Facebook
)

// User represents a user in the system
type User struct {
	ID                              uint   `gorm:"primary_key"`
	Name                            string `gorm:"type:varchar(100);"`
	Email                           string `gorm:"type:varchar(100);"`
	EmailVerified                   bool
	Password                        string `gorm:"type:varchar(100);"`
	AuthProvider                    Provider
	TwoFactor                       bool
	TwoFactorValidated              bool
	TwoFactorInitialSecret          string `gorm:"type:varchar(100);"`
	IsBlocked                       bool
	OpenToChangePasswordWithTokenID string `gorm:"type:varchar(100);"`
}
