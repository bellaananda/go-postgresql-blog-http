package service

import (
	"context"
	"errors"
	"log"

	// "fmt"
	// "log"
	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/bellaananda/go-postgresql-blog-http.git/repository"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		UserRepo: userRepo,
		db:       db,
	}
}

func (userService *UserService) CreateUser(ctx context.Context, user models.GormUser) (*models.GormUser, error) {
	_, err := userService.UserRepo.GetUserByEmail(ctx, user.Email)
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("user with this email already exists")
	}

	return userService.UserRepo.CreateUser(ctx, user)
}

func (userService *UserService) GetAllUsers(ctx context.Context) ([]models.GormUser, error) {
	var allUsers []models.GormUser
	if err := userService.db.WithContext(ctx).Find(&allUsers).Error; err != nil {
		return nil, err
	}

	return allUsers, nil
}

func (userService *UserService) GetUserByID(ctx context.Context, id uint) (*models.GormUser, error) {
	user, err := userService.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	// log.Printf("User found by id '%d': %+v\n", id, user)
	return user, nil
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (*models.GormUser, error) {
	user, err := userService.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.GormUser, error) {
	user, err := userService.UserRepo.GetUserByUsernameAndPassword(ctx, username, password)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (userService *UserService) UpdateUserByID(ctx context.Context, userID uint, user models.GormUser) (*models.GormUser, error) {
	// Check if the user ID in the URL matches the ID in the user object
	if userID != user.ID {
		return nil, errors.New("mismatched user ID in URL and request body")
	}

	existingUser, err := userService.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = userService.UserRepo.UpdateUser(ctx, userID, user)
	if err != nil {
		log.Printf("Error updating user with ID %d: %v", userID, err)
		return nil, err
	}

	return existingUser, nil
}

func (userService *UserService) DeleteUserByID(ctx context.Context, id uint) error {
	if err := userService.UserRepo.DeleteUser(ctx, id); err != nil {
		log.Printf("Error deleting user with ID %d: %v", id, err)
		return err
	}
	return nil
}
