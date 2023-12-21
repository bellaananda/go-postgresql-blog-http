package repository

import (
	"context"
	"errors"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (repo *PostgreSQLGORMRepository) MigrateComment(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormComment{})
	if err != nil {
		return err
	}
	return nil
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreateComment(ctx context.Context, comment models.GormComment) (*models.GormComment, error) {
	if err := repo.db.WithContext(ctx).Create(&comment).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &comment, nil
}

func (repo *PostgreSQLGORMRepository) AllComments(ctx context.Context) ([]models.GormComment, error) {
	var allComments []models.GormComment
	if err := repo.db.WithContext(ctx).Preload("User").Preload("Post").Find(&allComments).Error; err != nil {
		return nil, err
	}

	return allComments, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByID(ctx context.Context, id uint) (*models.GormComment, error) {
	var gormComment models.GormComment
	if err := repo.db.WithContext(ctx).Preload("User").Preload("Post").First(&gormComment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormComment, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByUserID(ctx context.Context, userid uint) ([]models.GormComment, error) {
	var gormComment []models.GormComment
	if err := repo.db.WithContext(ctx).Where("user_id = ?", userid).Find(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return gormComment, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByPostID(ctx context.Context, postid uint) ([]models.GormComment, error) {
	var gormComment []models.GormComment
	if err := repo.db.WithContext(ctx).Where("post_id = ?", postid).Find(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return gormComment, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByUserIDPostID(ctx context.Context, userid uint, postid uint) (*models.GormComment, error) {
	var gormComment models.GormComment
	if err := repo.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userid, postid).First(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormComment, nil
}

func (repo *PostgreSQLGORMRepository) UpdateComment(ctx context.Context, id uint, updated models.GormComment) (*models.GormComment, error) {
	updateRes := repo.db.WithContext(ctx).Where("id = ?", id).Save(&updated)
	if err := updateRes.Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := updateRes.RowsAffected
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &updated, nil
}

func (repo *PostgreSQLGORMRepository) DeleteComment(ctx context.Context, id uint) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormComment{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
