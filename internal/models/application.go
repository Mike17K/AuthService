package models

// User represents a user in the system
type Application struct {
	ID            string `gorm:"type:varchar(36);primary_key"`
	Name          string `gorm:"type:varchar(100);unique_index;not null;"`
	Password      string `gorm:"type:varchar(100);"`
	Description   string `gorm:"type:varchar(100);"`
	IsBlocked     bool
	BaseSecretKey string `gorm:"type:varchar(100);"` // This is the secret key that will be used to generate the JWT token
}
