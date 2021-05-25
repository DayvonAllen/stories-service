package services

import (
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DisLikeService interface {
	CreateDisLikeForStory(string, primitive.ObjectID) error
	CreateDisLikeForComment(string, primitive.ObjectID) error
	DeleteDisLikeByUsername(string) error
}

type DefaultDisLikeService struct {
	repo repo.DisLikeRepo
}

func (d DefaultDisLikeService) CreateDisLikeForStory(username string, id primitive.ObjectID) error {
	err := d.repo.CreateDisLikeForStory(username,id)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultDisLikeService) CreateDisLikeForComment(username string, id primitive.ObjectID) error {
	err := d.repo.CreateDisLikeForComment(username,id)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultDisLikeService) DeleteDisLikeByUsername(username string) error {
	err := d.repo.DeleteByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func NewDisLikeService(repository repo.DisLikeRepo) DefaultDisLikeService {
	return DefaultDisLikeService{repository}
}
