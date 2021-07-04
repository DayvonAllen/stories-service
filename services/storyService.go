package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryService interface {
	Create(dto *domain.CreateStoryDto) error
	UpdateById(primitive.ObjectID, string, string, string, *[]domain.Tag)  error
	FindAll(string) (*[]domain.Story, error)
	LikeStoryById(primitive.ObjectID, string) error
	DisLikeStoryById(primitive.ObjectID, string) error
	FindById(primitive.ObjectID, string) (*domain.StoryDto, error)
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

func (s DefaultStoryService) UpdateById(id primitive.ObjectID, newContent string, newTitle string, username string, tags *[]domain.Tag) error {
	err := s.repo.UpdateById(id, newContent, newTitle, username, tags)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) FindAll(page string) (*[]domain.Story, error) {
	story, err := s.repo.FindAll(page)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) FindById(id primitive.ObjectID, username string) (*domain.StoryDto, error) {
	story, err := s.repo.FindById(id, username)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (s DefaultStoryService) LikeStoryById(id primitive.ObjectID, username string) error {
	err := s.repo.LikeStoryById(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultStoryService) DisLikeStoryById(id primitive.ObjectID, username string) error {
	err := s.repo.DisLikeStoryById(id, username)
	if err != nil {
		return err
	}
	return nil
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
