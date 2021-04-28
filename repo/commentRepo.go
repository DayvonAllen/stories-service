package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentRepo interface {
	Create(id primitive.ObjectID) error
	FindById(id primitive.ObjectID) (*domain.Comment, error)
	FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.Comment, error)
	UpdateById(id primitive.ObjectID, newContent string) (*domain.Comment, error)
	DeleteById(id primitive.ObjectID) error
}
