package services

import (
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeService interface {
	CreateLikeForStory(string, primitive.ObjectID) error
	CreateLikeForComment(string, primitive.ObjectID) error
	DeleteLikeByUsername(string) error
}

type DefaultLikeService struct {
	repo repo.LikeRepo
}

func (c DefaultLikeService) CreateLikeForStory(username string, id primitive.ObjectID) error {
	err := c.repo.CreateLikeForStory(username,id)
	if err != nil {
		return err
	}
	return nil
}

func (c DefaultLikeService) CreateLikeForComment(username string, id primitive.ObjectID) error {
	err := c.repo.CreateLikeForComment(username,id)
	if err != nil {
		return err
	}
	return nil
}

func (c DefaultLikeService) DeleteLikeByUsername(username string) error {
	err := c.repo.DeleteByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func NewLikeService(repository repo.LikeRepo) DefaultLikeService {
	return DefaultLikeService{repository}
}
