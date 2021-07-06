package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Comment todo validate struct
type Comment struct {
	Id        		primitive.ObjectID `bson:"_id" json:"-"`
	StoryId        	primitive.ObjectID `bson:"storyId" json:"-"`
	Content string  `bson:"content" json:"content"`
	AuthorUsername string  `bson:"authorUsername" json:"-"`
	Likes  []string   			`bson:"likes" json:"-"`
	Dislikes []string  			 `bson:"dislikes" json:"-"`
	LikeCount int                `bson:"likeCount" json:"-"`
	DislikeCount int             `bson:"dislikeCount" json:"-"`
	FlagCount int				`bson:"flagCount" json:"-"`
	Flags []Flag				`bson:"flags" json:"-"`
	Replies []Comment 			`bson:"replies" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"-"`
}

type CommentDto struct {
	Content string  		`bson:"content" json:"content"`
	AuthorUsername string  `bson:"authorUsername" json:"authorUsername"`
	LikeCount int                `bson:"likeCount" json:"likeCount"`
	DislikeCount int             `bson:"dislikeCount" json:"dislikeCount"`
	Replies []CommentDto 		`bson:"replies" json:"replies"`
	CurrentUserLiked bool        `bson:"currentUserLiked" json:"currentUserLiked"`
	CurrentUserDisLiked bool        `bson:"currentUserDisLiked" json:"currentUserDisLiked"`
	CreatedDate string 				`json:"createdDate"`
}
