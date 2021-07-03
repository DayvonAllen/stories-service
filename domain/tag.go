package domain

import (
	"fmt"
	"strings"
)

type Tag struct {
	Value string `json:"value"`
	CreepyPasta bool `json:"-"`
	TrueScaryStory bool `json:"-"`
	CampFire bool `json:"-"`
	Paranormal bool `json:"-"`
	GhostStory bool `json:"-"`
	Other bool `json:"-"`
}

func (t Tag) ValidateTag(tagValidator *Tag) (string, error) {
	switch tag := strings.ToLower(t.Value); tag {
	case "creepypasta":
		if !tagValidator.CreepyPasta {
			tagValidator.CreepyPasta = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	case "truescarystory":
		if !tagValidator.TrueScaryStory {
			tagValidator.TrueScaryStory = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	case "campfire":
		if !tagValidator.CampFire {
			tagValidator.CampFire = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	case "ghoststory":
		if !tagValidator.GhostStory {
			tagValidator.GhostStory = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	case "paranormal":
		if !tagValidator.Paranormal {
			tagValidator.Paranormal = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	case "other":
		if !tagValidator.Other {
			tagValidator.Other = true
			return tag, nil
		}
		return "", fmt.Errorf("no duplicate tags")
	default:
		return "", fmt.Errorf("invalid tag")
	}
}