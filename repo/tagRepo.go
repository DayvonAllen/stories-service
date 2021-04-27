package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagRepo interface {
	Create(*domain.Tag) error
	CreateMany(*[]interface{}) error
	FindByTagName(string) (*domain.Tag, error)
	FindAll() (*[]domain.Tag, error)
	DeleteById(id primitive.ObjectID) error
}