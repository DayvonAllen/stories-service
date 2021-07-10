package repo

import (
	"example.com/app/config"
	"example.com/app/domain"
	"example.com/app/events"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProcessMessage(message domain.Message) error {

	if message.ResourceType == "user" {
		// 201 is the created messageType
		if message.MessageType == 201 {
			user := message.User
			err := UserRepoImpl{}.Create(&user)

			if err != nil {
				return err
			}
			return nil
		}

		// 200 is the updated messageType
		if message.MessageType == 200 {
			user := message.User

			err := UserRepoImpl{}.UpdateByID(&user)
			if err != nil {
				return err
			}
			return nil
		}

		// 204 is the deleted messageType
		if message.MessageType == 204 {
			user := message.User

			err := UserRepoImpl{}.DeleteByID(&user)

			if err != nil {
				return err
			}
			return nil
		}
	}


	return fmt.Errorf("cannot process this message")
}

func PushUserToQueue(message []byte) error {

	producer := events.GetInstance()

	msg := &sarama.ProducerMessage{
		Topic: config.Config("PRODUCER_TOPIC"),
		Value: sarama.StringEncoder(message),
	}


	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		err = producer.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Failed to send message to the queue")
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "user", partition, offset)
	return nil
}

func SendKafkaMessage(story *domain.Story, eventType int) error {
	um := new(domain.Message)
	um.Story = *story

	// user created/updated event
	um.MessageType = eventType
	um.ResourceType = "story"

	fmt.Println(um.Story)
	//turn user struct into a byte array
	b, err := msgpack.Marshal(um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b)

	if err != nil {
		return err
	}

	return nil
}

func HandleKafkaMessage(err error, story *domain.Story, messageType int) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	err = SendKafkaMessage(story, messageType)

	if err != nil {
		fmt.Println("Failed to publish new user")
	}

	return nil
}