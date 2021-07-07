package repo

import (
	"example.com/app/database"
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CommentRepo interface {
	Create(comment *domain.Comment, mongoCollection *mongo.Collection, conn *database.Connection, dbType string) error
	FindAllCommentsByResourceId(id primitive.ObjectID, username string) (*[]domain.CommentDto, error)
	UpdateById(id primitive.ObjectID, newContent string, edited bool, updatedTime time.Time, username string) (*domain.Comment, error)
	LikeCommentById(primitive.ObjectID, string) error
	DisLikeCommentById(primitive.ObjectID, string) error
	UpdateFlagCount(flag *domain.Flag) error
	DeleteById(id primitive.ObjectID, username string) error
}
