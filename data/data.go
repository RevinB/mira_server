package data

import (
	"context"
	"fmt"
	"github.com/RevinB/mira_server/data/upload"
	"github.com/RevinB/mira_server/data/user"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
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
	Cache() *redis.Client

	User() *user.Store
	Upload() *upload.Store
}

type storeImpl struct {
	client *gorm.DB
	cache  *redis.Client

	user   *user.Store
	upload *upload.Store
}

func NewStore() (Store, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	intDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rc := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       intDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := rc.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	si := &storeImpl{
		client: db,
		cache:  rc,
		user:   user.NewStore(db),
		upload: upload.NewStore(db),
	}

	return si, nil
}

func (si *storeImpl) Client() *gorm.DB {
	return si.client
}

func (si *storeImpl) Cache() *redis.Client {
	return si.cache
}

func (si *storeImpl) User() *user.Store {
	return si.user
}

func (si *storeImpl) Upload() *upload.Store {
	return si.upload
}
