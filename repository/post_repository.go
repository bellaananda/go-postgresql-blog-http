package repository

import (
	"context"
	"errors"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (repo *PostgreSQLGORMRepository) MigratePost(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormPost{})
	if err != nil {
		return err
	}
	return nil
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreatePost(ctx context.Context, post models.GormPost) (*models.GormPost, error) {
	if err := repo.db.WithContext(ctx).Create(&post).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &post, nil
}

func (repo *PostgreSQLGORMRepository) AllPosts(ctx context.Context) ([]models.GormPost, error) {
	var allPosts []models.GormPost
	if err := repo.db.WithContext(ctx).Preload("User").Find(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByID(ctx context.Context, id uint) (*models.GormPost, error) {
	var gormPost models.GormPost
	if err := repo.db.WithContext(ctx).Preload("User").First(&gormPost, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormPost, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByTitle(ctx context.Context, title string) (*models.GormPost, error) {
	var gormPost models.GormPost
	if err := repo.db.WithContext(ctx).Where("title = ?", title).First(&gormPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormPost, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByUserID(ctx context.Context, userid uint) ([]models.GormPost, error) {
	var gormPost []models.GormPost
	if err := repo.db.WithContext(ctx).Where("user_id = ?", userid).Find(&gormPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return gormPost, nil
}

func (repo *PostgreSQLGORMRepository) UpdatePost(ctx context.Context, id uint, updated models.GormPost) (*models.GormPost, error) {
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

func (repo *PostgreSQLGORMRepository) DeletePost(ctx context.Context, id uint) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormPost{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
