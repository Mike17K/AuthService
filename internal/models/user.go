package models

// Provider is an enum for authentication providers
type Provider int

const (
	Email Provider = iota
	Google
	Facebook
)

type UserType int

const (
	Admin UserType = iota
	SimpleUser
)

// User represents a user in the system
type User struct {
	ID                              string `gorm:"type:varchar(36);primary_key"`
	Name                            string `gorm:"type:varchar(100);"`
	Email                           string `gorm:"type:varchar(100);unique_index;not null;"`
	EmailVerified                   bool
	Password                        string `gorm:"type:varchar(100);"`
	AuthProvider                    Provider
	TwoFactor                       bool
	TwoFactorValidated              bool
	TwoFactorInitialSecret          string `gorm:"type:varchar(100);"`
	IsBlocked                       bool
	OpenToChangePasswordWithTokenID string `gorm:"type:varchar(100);"`
	ApplicationID                   string `gorm:"type:varchar(36);"`
	UserType                        UserType
	RefreshToken                    string `gorm:"type:varchar(1000);"`
}
