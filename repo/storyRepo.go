package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepo interface {
	 Create(story *domain.CreateStoryDto) error
	 UpdateById(primitive.ObjectID, string, string, string, *[]domain.Tag) (*domain.StoryDto, error)
	 FindAll(string) (*[]domain.Story, error)
	 FindById(primitive.ObjectID) (*domain.StoryDto, error)
	 DeleteById(primitive.ObjectID, string) error
}
