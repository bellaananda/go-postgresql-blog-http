package repository

import (
	"context"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
)

// Repository provides access to the website storage.
type PostRepository interface {
	MigratePost(ctx context.Context) error
	CreatePost(ctx context.Context, post models.GormPost) (*models.GormPost, error)
	AllPosts(ctx context.Context) ([]models.GormPost, error)
	GetPostByID(ctx context.Context, id uint) (*models.GormPost, error)
	GetPostByTitle(ctx context.Context, title string) (*models.GormPost, error)
	GetPostByUserID(ctx context.Context, userid uint) ([]models.GormPost, error)
	UpdatePost(ctx context.Context, id uint, updated models.GormPost) (*models.GormPost, error)
	DeletePost(ctx context.Context, id uint) error
}
