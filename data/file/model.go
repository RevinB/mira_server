package file

import "time"

type Model struct {
	// ID short 5-8 character long file identifier (same as name in URL)
	ID string `json:"id" gorm:"primaryKey"`

	// FileExtension without dot separator
	FileExtension string `json:"file_extension"`

	MIMEType string `json:"mime_type"`

	// Owner foreign id to users table
	Owner string `json:"owner" gorm:"not null"`

	// CreatedAt automatically handled by gorm
	CreatedAt time.Time `json:"created_at" gorm:"not null"`

	// DeletedAt automatically handled by gorm
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Model) TableName() string {
	return "files"
}
