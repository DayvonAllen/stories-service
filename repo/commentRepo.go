package repo

import (
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentRepo interface {
	Create(comment *domain.Comment) error
	FindById(id primitive.ObjectID) (*domain.Comment, error)
	FindAllCommentsByStoryId(id primitive.ObjectID) (*[]domain.CommentDto, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time) (*domain.Comment, error)
	DeleteById(id primitive.ObjectID) error
}
