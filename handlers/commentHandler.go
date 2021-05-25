package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	CommentService services.CommentService
}

func (d *CommentHandler) CreateComment(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = d.CommentService.Create(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (d *CommentHandler) FindById(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	commentData, err := d.CommentService.FindById(comment.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": commentData})
}

func (d *CommentHandler) FindAllCommentsByStoryId(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	comments, err := d.CommentService.FindAllCommentsByStoryId(comment.StoryId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": comments})
}

func (d *CommentHandler) UpdateById(c *fiber.Ctx) error {
	c.Accepts("application/json")
	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	commentData, err := d.CommentService.UpdateById(comment.Id, comment.Content)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": commentData})
}

func (d *CommentHandler) DeleteById(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	err = d.CommentService.DeleteById(u.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}



