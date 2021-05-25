package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type LikeHandler struct {
	LikeService services.LikeService
}

func (l *LikeHandler) CreateLikeForStory(c *fiber.Ctx) error {
	c.Accepts("application/json")
	like := new(domain.Like)

	err := c.BodyParser(like)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = l.LikeService.CreateLikeForStory(like.AuthorUsername, like.ContentId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (l *LikeHandler) CreateLikeForComment(c *fiber.Ctx) error {
	c.Accepts("application/json")
	like := new(domain.Like)

	err := c.BodyParser(like)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = l.LikeService.CreateLikeForComment(like.AuthorUsername, like.ContentId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (l *LikeHandler) DeleteLikeByUsername(c *fiber.Ctx) error {
	c.Accepts("application/json")
	like := new(domain.Like)

	err := c.BodyParser(like)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = l.LikeService.DeleteLikeByUsername(like.AuthorUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}