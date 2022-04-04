package data

import (
	"fmt"
	"github.com/RevinB/mira_server/data/file"
	"github.com/RevinB/mira_server/data/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Store struct {
	Client *gorm.DB

	Users *user.Store
	Files *file.Store
}

func NewStore() (*Store, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	si := &Store{
		Client: db,
		Users:  user.NewStore(db),
		Files:  file.NewStore(db),
	}

	return si, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		user.Model{},
		file.Model{},
	)
}
