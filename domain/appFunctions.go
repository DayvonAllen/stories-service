package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

func CreateTag(tagName string) Tag {
	tag := new(Tag)
	tag.Id = primitive.NewObjectID()
	tag.TagName = tagName
	return *tag
}