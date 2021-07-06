package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentService interface {
	Create(comment *domain.Comment) error
	FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.CommentDto, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) (*domain.Comment, error)
	LikeCommentById(primitive.ObjectID, string) error
	DisLikeCommentById(primitive.ObjectID, string) error
	DeleteById(id primitive.ObjectID, username string) error
}

type DefaultCommentService struct {
	repo repo.CommentRepo
}

func (c DefaultCommentService) Create(comment *domain.Comment) error {
	err := c.repo.Create(comment)
	if err != nil {
		return err
	}
	return nil
}

func (c DefaultCommentService) FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.CommentDto, error) {
	comment, err := c.repo.FindAllCommentsByStoryId(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (c DefaultCommentService) UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) (*domain.Comment, error) {
	comment, err := c.repo.UpdateById(id, newContent, edited, updatedTime, username)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (c DefaultCommentService) LikeCommentById(id primitive.ObjectID, username string) error {
	err := c.repo.LikeCommentById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (c DefaultCommentService) DisLikeCommentById(id primitive.ObjectID, username string) error {
	err := c.repo.DisLikeCommentById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (c DefaultCommentService) DeleteById(id primitive.ObjectID, username string) error {
	err := c.repo.DeleteById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewCommentService(repository repo.CommentRepo) DefaultCommentService {
	return DefaultCommentService{repository}
}