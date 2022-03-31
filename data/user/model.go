package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	// ID is a UUID, cannot be changed
	// links uploads to users
	// this can be used to regenerate the secret key and JWT ALONE
	// essentially a "master password"
	ID string `json:"id" gorm:"primaryKey"`

	// SecretKey is only used for revoking JWTs
	SecretKey string `json:"secret_key" gorm:"unique"`

	// IsAdmin allows user to have admin perms
	IsAdmin bool `json:"isAdmin"`
}

func (Model) TableName() string {
	return "users"
}

func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
