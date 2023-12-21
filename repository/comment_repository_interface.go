package repository

import (
	"context"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
)

// Repository provides access to the website storage.
type CommentRepository interface {
	MigrateComment(ctx context.Context) error
	CreateComment(ctx context.Context, comment models.GormComment) (*models.GormComment, error)
	AllComments(ctx context.Context) ([]models.GormComment, error)
	GetCommentByID(ctx context.Context, id uint) (*models.GormComment, error)
	GetCommentByUserID(ctx context.Context, userid uint) ([]models.GormComment, error)
	GetCommentByPostID(ctx context.Context, postid uint) ([]models.GormComment, error)
	GetCommentByUserIDPostID(ctx context.Context, userid uint, postid uint) (*models.GormComment, error)
	UpdateComment(ctx context.Context, id uint, updated models.GormComment) (*models.GormComment, error)
	DeleteComment(ctx context.Context, id uint) error
}
