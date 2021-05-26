package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagService interface {
	Create(*domain.Tag) error
	CreateMany(*[]interface{}) error
	FindByTagName(string) (*domain.Tag, error)
	FindAll() (*[]domain.Tag, error)
	DeleteById(id primitive.ObjectID) error
}

type DefaultTagService struct {
	repo repo.TagRepo
}

func (t DefaultTagService) Create(tag *domain.Tag) error {
	err := t.repo.Create(tag)
	if err != nil {
		return err
	}
	return nil
}

func (t DefaultTagService) CreateMany(tags *[]interface{}) error {
	err := t.repo.CreateMany(tags)
	if err != nil {
		return err
	}
	return nil
}

func (t DefaultTagService) FindByTagName(tagName string) (*domain.Tag, error) {
	tag, err := t.repo.FindByTagName(tagName)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t DefaultTagService) FindAll() (*[]domain.Tag, error) {
	tag, err := t.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t DefaultTagService) DeleteById(id primitive.ObjectID) error {
	err := t.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func NewTagService(repository repo.TagRepo) DefaultTagService {
	return DefaultTagService{repository}
}
