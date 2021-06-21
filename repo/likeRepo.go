package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeRepo interface {
	CreateLikeForStory(string, primitive.ObjectID) error
	CreateLikeForComment(string, primitive.ObjectID) error
	DeleteByUsername(string) error
}
