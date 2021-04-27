package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRepo interface {
	 Create(id primitive.ObjectID) error
	 UpdateById(primitive.ObjectID, primitive.ObjectID) (*domain.Story, error)
	 FindAll(id primitive.ObjectID) (*[]domain.Story, error)
	 FindById(primitive.ObjectID, primitive.ObjectID) (*domain.Story, error)
	 DeleteById(primitive.ObjectID, primitive.ObjectID) error
}
