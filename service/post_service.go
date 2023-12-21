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

type PostService struct {
	PostRepo repository.PostRepository
	db       *gorm.DB
}

func NewPostService(postRepo repository.PostRepository, db *gorm.DB) *PostService {
	return &PostService{
		PostRepo: postRepo,
		db:       db,
	}
}

func (postService *PostService) CreatePost(ctx context.Context, post models.GormPost) (*models.GormPost, error) {
	_, err := postService.PostRepo.GetPostByTitle(ctx, post.Title)
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("a post with this title already exists")
	}

	return postService.PostRepo.CreatePost(ctx, post)
}

func (postService *PostService) GetAllPosts(ctx context.Context) ([]models.GormPost, error) {
	posts, err := postService.PostRepo.AllPosts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (postService *PostService) GetPostByID(ctx context.Context, id uint) (*models.GormPost, error) {
	post, err := postService.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}

func (postService *PostService) GetPostByTitle(ctx context.Context, title string) (*models.GormPost, error) {
	post, err := postService.PostRepo.GetPostByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}

func (postService *PostService) GetPostByUserID(ctx context.Context, userid uint) ([]models.GormPost, error) {
	post, err := postService.PostRepo.GetPostByUserID(ctx, userid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}
	return post, nil
}

func (postService *PostService) UpdatePostByID(ctx context.Context, postID uint, post models.GormPost) (*models.GormPost, error) {
	// Check if the post ID in the URL matches the ID in the post object
	if postID != post.ID {
		return nil, errors.New("mismatched post ID in URL and request body")
	}

	existingPost, err := postService.PostRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	_, err = postService.PostRepo.UpdatePost(ctx, postID, post)
	if err != nil {
		log.Printf("Error updating post with ID %d: %v", postID, err)
		return nil, err
	}

	return existingPost, nil
}

func (postService *PostService) DeletePostByID(ctx context.Context, id uint) error {
	if err := postService.PostRepo.DeletePost(ctx, id); err != nil {
		log.Printf("Error deleting post with ID %d: %v", id, err)
		return err
	}
	return nil
}
