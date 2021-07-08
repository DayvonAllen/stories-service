package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                   primitive.ObjectID `bson:"_id" json:"id"`
	Email                string             `bson:"email" json:"email"`
	Username             string             `bson:"username" json:"username"`
	CurrentTagLine       string             `bson:"currentTagLine" json:"currentTagLine"`
	UnlockedTagLine      []string           `bson:"unlockedTagLine" json:"unlockedTagLine"`
	ProfilePictureUrl    string             `bson:"profilePictureUrl" json:"profilePictureUrl"`
	CurrentBadgeUrl      string             `bson:"currentBadgeUrl" json:"currentBadgeUrl"`
	UnlockedBadgesUrls   []string           `bson:"unlockedBadgesUrls" json:"unlockedBadgesUrls"`
	Followers            []string           `bson:"followers" json:"followers"`
	Following            []string           `bson:"following" json:"following"`
	FollowerCount        int                `bson:"followerCount" json:"followerCount"`
	DisplayFollowerCount bool               `bson:"displayFollowerCount" json:"displayFollowerCount"`
	ProfileIsViewable    bool               `bson:"profileIsViewable" json:"profileIsViewable"`
	AcceptMessages       bool               `bson:"acceptMessages" json:"acceptMessages"`
	IsVerified           bool               `bson:"isVerified" json:"-"`
}

type ViewUserProfile struct {
	Username                    string     `json:"username"`
	CurrentTagLine              string     `json:"currentTagLine"`
	ProfilePictureUrl           string     `json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string     `json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string     `json:"currentBadgeUrl"`
	FollowerCount               int        `json:"followerCount"`
	ProfileIsViewable           bool       `json:"-"`
	DisplayFollowerCount        bool       `json:"displayFollowerCount"`
	IsFollowing                 bool       `json:"isFollowing"`
	Post                        []StoryDto `json:"post"`
}

type CurrentUserProfile struct {
	Username                    string     `json:"username"`
	CurrentTagLine              string     `json:"currentTagLine"`
	ProfilePictureUrl           string     `json:"profilePictureUrl"`
	ProfileBackgroundPictureUrl string     `json:"profileBackgroundPictureUrl"`
	CurrentBadgeUrl             string     `json:"currentBadgeUrl"`
	FollowerCount               int        `json:"followerCount"`
	Post                        []StoryDto `json:"post"`
}
