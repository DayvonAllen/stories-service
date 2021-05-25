package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type DisLikeHandler struct {
	DisLikeService services.DisLikeService
}

func (d *DisLikeHandler) CreateDisLikeForStory(c *fiber.Ctx) error {
	c.Accepts("application/json")
	disLike := new(domain.Dislike)

	err := c.BodyParser(disLike)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = d.DisLikeService.CreateDisLikeForStory(disLike.AuthorUsername, disLike.ContentId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (d *DisLikeHandler) CreateDisLikeForComment(c *fiber.Ctx) error {
	c.Accepts("application/json")
	disLike := new(domain.Dislike)

	err := c.BodyParser(disLike)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = d.DisLikeService.CreateDisLikeForComment(disLike.AuthorUsername, disLike.ContentId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (d *DisLikeHandler) DeleteDisLikeByUsername(c *fiber.Ctx) error {
	c.Accepts("application/json")
	disLike := new(domain.Dislike)

	err := c.BodyParser(disLike)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = d.DisLikeService.DeleteDisLikeByUsername(disLike.AuthorUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
