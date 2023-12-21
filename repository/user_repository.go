package repository

import (
	"context"
	"errors"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type PostgreSQLGORMRepository struct {
	db *gorm.DB
}

func (repo *PostgreSQLGORMRepository) MigrateUser(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormUser{})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreateUser(ctx context.Context, user models.GormUser) (*models.GormUser, error) {
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &user, nil
}

func (repo *PostgreSQLGORMRepository) AllUsers(ctx context.Context) ([]models.GormUser, error) {
	var allUsers []models.GormUser
	if err := repo.db.WithContext(ctx).Find(&allUsers).Error; err != nil {
		return nil, err
	}

	return allUsers, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByID(ctx context.Context, id uint) (*models.GormUser, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("id = ?", id).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormUser, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByEmail(ctx context.Context, email string) (*models.GormUser, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("email = ?", email).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormUser, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.GormUser, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("username = ? AND password = ?", username, password).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &gormUser, nil
}

func (repo *PostgreSQLGORMRepository) UpdateUser(ctx context.Context, id uint, updated models.GormUser) (*models.GormUser, error) {
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

func (repo *PostgreSQLGORMRepository) DeleteUser(ctx context.Context, id uint) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormUser{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
