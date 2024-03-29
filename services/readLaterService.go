package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadLaterService interface {
	Create(username string, storyId primitive.ObjectID) error
	GetByUsername(username string) (*domain.ReadLaterDto, error)
	Delete(id primitive.ObjectID, username string) error
}

type DefaultReadLaterService struct {
	repo repo.ReadLaterRepo
}

func (s DefaultReadLaterService) Create(username string, storyId primitive.ObjectID) error {
	err := s.repo.Create(username, storyId)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultReadLaterService) GetByUsername(username string) (*domain.ReadLaterDto, error) {
	readLaterItems, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return readLaterItems, nil
}

func (s DefaultReadLaterService) Delete(id primitive.ObjectID, username string) error {
	err := s.repo.Delete(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewReadLaterService(repository repo.ReadLaterRepo) DefaultReadLaterService {
	return DefaultReadLaterService{repository}
}