package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadLaterRepo interface {
	Create(username string, storyId primitive.ObjectID) error
	GetByUsername(username string) (*domain.ReadLaterDto, error)
	Delete(id primitive.ObjectID, username string) error
}
