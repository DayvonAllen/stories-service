package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryService interface {
	Create(dto *domain.CreateStoryDto) error
	UpdateById(primitive.ObjectID, string, string) (*domain.StoryDto, error)
	FindAll() (*[]domain.Story, error)
	FindById(primitive.ObjectID) (*domain.StoryDto, error)
	DeleteById(primitive.ObjectID, string) error
}

type DefaultStoryService struct {
	repo repo.StoryRepo
}

func (s DefaultStoryService) Create(story *domain.CreateStoryDto) error {
	err := s.repo.Create(story)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) UpdateById(id primitive.ObjectID, newContent string, newTitle string) (*domain.StoryDto, error) {
	story, err := s.repo.UpdateById(id, newContent, newTitle)
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

func (s DefaultStoryService) FindById(id primitive.ObjectID) (*domain.StoryDto, error) {
	story, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) DeleteById(id primitive.ObjectID, username string) error {
	err := s.repo.DeleteById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func NewStoryService(repository repo.StoryRepo) DefaultStoryService {
	return DefaultStoryService{repository}
}
