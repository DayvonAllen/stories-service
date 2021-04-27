package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Content string  `bson:"content" json:"content"`
	AuthorId primitive.ObjectID  `bson:"authorId" json:"authorId"`
	Likes  []primitive.ObjectID   `bson:"likes" json:"likes"`
	Dislikes []primitive.ObjectID  `bson:"dislikes" json:"dislikes"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	FlagCount []primitive.ObjectID	`bson:"flagCount" json:"-"`
	Replies []primitive.ObjectID `bson:"replies" json:"replies"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}
