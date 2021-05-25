package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DisLikeRepo interface {
	CreateDisLikeForStory(string, primitive.ObjectID) error
	CreateDisLikeForComment(string, primitive.ObjectID) error
	DeleteByUsername(string) error
}
