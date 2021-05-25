package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepo interface {
	 Create(id primitive.ObjectID) error
	 UpdateById(primitive.ObjectID, string) (*domain.Story, error)
	 FindAll() (*[]domain.Story, error)
	 FindById(primitive.ObjectID) (*domain.Story, error)
	 DeleteById(primitive.ObjectID) error
}
