package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Email string                 `json:"email"`
	Username string              `json:"username"`
	CurrentTagLine  string       `json:"currentTagLine"`
	UnlockedTagLine  []string    `json:"unlockedTagLine"`
	ProfilePictureUrl  string    `json:"profilePictureUrl"`
	CurrentBadgeUrl  string      `json:"currentBadgeUrl"`
	UnlockedBadgesUrls  []string `json:"unlockedBadgesUrls"`
	ProfileIsViewable  bool      `json:"profileIsViewable"`
	AcceptMessages  bool         `json:"acceptMessages"`
	IsVerified  bool             `bson:"isVerified" json:"-"`
}
