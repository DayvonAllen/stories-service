package events

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"fmt"
)


func ProcessMessage(message domain.UserMessage) error {

	// 201 is the created messageType
	if message.MessageType == 201 {
		err := repo.UserRepoImpl{}.Create(&message.User)

		if err != nil {
			return err
		}
		return nil
	}

	// 200 is the updated messageType
	if message.MessageType == 200 {
		err := repo.UserRepoImpl{}.UpdateByID(&message.User)
		if err != nil {
			return err
		}
		return  nil
	}

	// 204 is the deleted messageType
	if message.MessageType == 204 {
		err := repo.UserRepoImpl{}.DeleteByID(&message.User)

		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("cannot process this message")
}
