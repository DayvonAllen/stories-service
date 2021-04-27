package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Like struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	AuthorId  primitive.ObjectID `bson:"authorId" json:"authorId"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
}