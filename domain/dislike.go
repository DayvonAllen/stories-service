package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Dislike struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	AuthorUsername  string `bson:"authorUsername" json:"authorUsername"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
}
