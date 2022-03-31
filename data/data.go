package data

import (
	"fmt"
	"github.com/RevinB/mira_server/data/upload"
	"github.com/RevinB/mira_server/data/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var _ Store = (*storeImpl)(nil)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		user.Model{},
		upload.Model{},
	)
}

type Store interface {
	Client() *gorm.DB

	User() user.Store
	Upload() upload.Store
}

type storeImpl struct {
	client *gorm.DB

	user   user.Store
	upload upload.Store
}

func NewStore() (Store, error) {
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

	si := storeImpl{
		client: db,
		user:   user.NewStore(db),
		upload: upload.NewStore(db),
	}

	return &si, nil
}

func (si *storeImpl) Client() *gorm.DB {
	return si.client
}

func (si *storeImpl) User() user.Store {
	return si.user
}

func (si *storeImpl) Upload() upload.Store {
	return si.upload
}
