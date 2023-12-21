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

type CommentService struct {
	CommentRepo repository.CommentRepository
	db          *gorm.DB
}

func NewCommentService(commentRepo repository.CommentRepository, db *gorm.DB) *CommentService {
	return &CommentService{
		CommentRepo: commentRepo,
		db:          db,
	}
}

func (commentService *CommentService) CreateComment(ctx context.Context, comment models.GormComment) (*models.GormComment, error) {
	_, err := commentService.CommentRepo.GetCommentByUserIDPostID(ctx, comment.UserID, comment.PostID)
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("a comment with the post already exists")
	}

	return commentService.CommentRepo.CreateComment(ctx, comment)
}

func (commentService *CommentService) GetAllComments(ctx context.Context) ([]models.GormComment, error) {
	post, err := commentService.CommentRepo.AllComments(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}

func (commentService *CommentService) GetCommentByID(ctx context.Context, id uint) (*models.GormComment, error) {
	comment, err := commentService.CommentRepo.GetCommentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return comment, nil
}

func (commentService *CommentService) GetCommentByUserID(ctx context.Context, userid uint) ([]models.GormComment, error) {
	comment, err := commentService.CommentRepo.GetCommentByUserID(ctx, userid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}
	return comment, nil
}

func (commentService *CommentService) GetCommentByPostID(ctx context.Context, postid uint) ([]models.GormComment, error) {
	comment, err := commentService.CommentRepo.GetCommentByPostID(ctx, postid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return comment, nil
}

func (commentService *CommentService) GetCommentByUserIDPostID(ctx context.Context, userid uint, postid uint) (*models.GormComment, error) {
	comment, err := commentService.CommentRepo.GetCommentByUserIDPostID(ctx, userid, postid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			log.Printf("Post with user id '%d' and post id '%d' does not exist in the repository\n", userid, postid)
			return nil, err
		}
		log.Printf("Error getting post by user id '%d' and post id '%d': %v\n", userid, postid, err)
		return nil, err
	}
	log.Printf("Post found by user id '%d' and post id '%d': %+v\n", userid, postid, comment)
	return comment, nil
}

func (commentService *CommentService) UpdateCommentByID(ctx context.Context, commentID uint, comment models.GormComment) (*models.GormComment, error) {
	// Check if the comment ID in the URL matches the ID in the comment object
	if commentID != comment.ID {
		return nil, errors.New("mismatched comment ID in URL and request body")
	}

	existingComment, err := commentService.CommentRepo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	if _, err := commentService.CommentRepo.UpdateComment(ctx, commentID, comment); err != nil {
		log.Printf("Error updating comment with ID %d: %v", commentID, err)
		return nil, err
	}
	return existingComment, nil
}

func (commentService *CommentService) DeleteCommentByID(ctx context.Context, id uint) error {
	if err := commentService.CommentRepo.DeleteComment(ctx, id); err != nil {
		log.Printf("Error deleting post with ID %d: %v", id, err)
		return err
	}
	return nil
}
