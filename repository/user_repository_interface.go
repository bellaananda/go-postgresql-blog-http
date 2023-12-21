package repository

import (
	"context"
	"errors"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

// Repository provides access to the website storage.
type UserRepository interface {
	MigrateUser(ctx context.Context) error
	CreateUser(ctx context.Context, user models.GormUser) (*models.GormUser, error)
	AllUsers(ctx context.Context) ([]models.GormUser, error)
	GetUserByID(ctx context.Context, id uint) (*models.GormUser, error)
	GetUserByEmail(ctx context.Context, email string) (*models.GormUser, error)
	GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (*models.GormUser, error)
	UpdateUser(ctx context.Context, id uint, updated models.GormUser) (*models.GormUser, error)
	DeleteUser(ctx context.Context, id uint) error
}
