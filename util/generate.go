package util

import (
	"example.com/app/domain"
)

func GenerateTags() *[]interface{} {
	tags := []string{"creepypasta", "paranormal", "ghost story", "true scary story", "campfire", "other"}
	tagsSlice := make([]interface{}, 0, len(tags))

	for _, tag := range tags{
		tagsSlice = append(tagsSlice, domain.CreateTag(tag))
	}
	return &tagsSlice
}
