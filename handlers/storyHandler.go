package handlers

import (
	"example.com/app/domain"
	"example.com/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StoryHandler struct {
	StoryService services.StoryService
}

func (s *StoryHandler) CreateStory(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	c.Accepts("application/json")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	storyDto := new(domain.CreateStoryDto)

	err = c.BodyParser(storyDto)

	storyDto.AuthorUsername = u.Username
	storyDto.CreatedAt = time.Now()
	storyDto.UpdatedAt = time.Now()

	tagLength := len(storyDto.Tags)

	if tagLength < 1 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("story must have at least one tag")})
	}

	t := new(domain.Tag)
	for _, tag := range storyDto.Tags {
		err = tag.ValidateTag(t)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.Create(storyDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")

	stories, err := s.StoryService.FindAll(page)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": stories})
}

func (s *StoryHandler) UpdateStory(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	c.Accepts("application/json")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	storyDto := new(domain.UpdateStoryDto)

	err = c.BodyParser(storyDto)

	storyDto.UpdatedAt = time.Now()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	tagLength := len(storyDto.Tags)

	if tagLength < 1 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("story must have at least one tag")})
	}

	t := new(domain.Tag)
	for _, tag := range storyDto.Tags {
		err = tag.ValidateTag(t)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
	}

	data, err := s.StoryService.UpdateById(id, storyDto.Content, storyDto.Title, u.Username, &storyDto.Tags)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": &data})
}

func (s *StoryHandler) DeleteStory(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.DeleteById(id, u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}