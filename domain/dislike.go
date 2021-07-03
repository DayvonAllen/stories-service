package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Dislike todo validate struct
type Dislike struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	AuthorUsername  string `bson:"authorUsername" json:"authorUsername"`
	ContentId  primitive.ObjectID `bson:"contentId" json:"contentId"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
}
