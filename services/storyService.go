package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryService interface {
	Create(id primitive.ObjectID) error
	UpdateById(primitive.ObjectID, string) (*domain.Story, error)
	FindAll() (*[]domain.Story, error)
	FindById(primitive.ObjectID) (*domain.Story, error)
	DeleteById(primitive.ObjectID) error
}

type DefaultStoryService struct {
	repo repo.StoryRepo
}

func (s DefaultStoryService) Create(id primitive.ObjectID) error {
	err := s.repo.Create(id)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) UpdateById(id primitive.ObjectID, newContent string) (*domain.Story, error) {
	story, err := s.repo.UpdateById(id, newContent)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FindAll() (*[]domain.Story, error) {
	story, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FindById(id primitive.ObjectID) (*domain.Story, error) {
	story, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) DeleteById(id primitive.ObjectID) error {
	err := s.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func NewStoryService(repository repo.StoryRepo) DefaultStoryService {
	return DefaultStoryService{repository}
}
