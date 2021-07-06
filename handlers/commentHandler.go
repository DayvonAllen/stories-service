package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentHandler struct {
	CommentService services.CommentService
}

func (ch *CommentHandler) CreateCommentOnStory(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	comment := new(domain.Comment)

	err = c.BodyParser(comment)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	comment.Id = primitive.NewObjectID()
	comment.AuthorUsername = u.Username
	comment.StoryId = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	err = ch.CommentService.Create(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) FindById(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	commentData, err := ch.CommentService.FindById(comment.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": commentData})
}

//func (ch *CommentHandler) FindAllCommentsByStoryId(c *fiber.Ctx) error {
//	id, err := primitive.ObjectIDFromHex(c.Params("id"))
//
//	if err != nil {
//		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//	}
//
//	comments, err := ch.CommentService.FindAllCommentsByStoryId(id)
//
//	if err != nil {
//		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
//	}
//
//	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": comments})
//}

func (ch *CommentHandler) UpdateById(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	commentData, err := ch.CommentService.UpdateById(comment.Id, comment.Content)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": commentData})
}

func (ch *CommentHandler) DeleteById(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	err = ch.CommentService.DeleteById(u.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}



