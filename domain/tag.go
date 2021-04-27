package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	TagName  string `bson:"tagName" json:"tagName"`
}
