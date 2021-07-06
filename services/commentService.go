package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	Create(comment *domain.Comment) error
	FindById(id primitive.ObjectID) (*domain.Comment, error)
	FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.CommentDto, error)
	UpdateById(id primitive.ObjectID, newContent string) (*domain.Comment, error)
	DeleteById(id primitive.ObjectID) error
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

func (c DefaultCommentService) FindById(id primitive.ObjectID) (*domain.Comment, error) {
	comment, err := c.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (c DefaultCommentService) FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.CommentDto, error) {
	comment, err := c.repo.FindAllCommentsByStoryId(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (c DefaultCommentService) UpdateById(id primitive.ObjectID, newContent string) (*domain.Comment, error) {
	comment, err := c.repo.UpdateById(id, newContent)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (c DefaultCommentService) DeleteById(id primitive.ObjectID) error {
	err := c.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func NewCommentService(repository repo.CommentRepo) DefaultCommentService {
	return DefaultCommentService{repository}
}