package router

import (
	"example.com/app/handlers"
	"example.com/app/middleware"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	ch := handlers.CommentHandler{CommentService: services.NewCommentService(repo.NewCommentRepoImpl())}
	sh := handlers.StoryHandler{StoryService: services.NewStoryService(repo.NewStoryRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	stories := api.Group("/stories")
	stories.Post("/", sh.CreateStory)
	stories.Put("/:id", sh.UpdateStory)
	stories.Put("/like/:id", sh.LikeStory)
	stories.Put("/dislike/:id", sh.DisLikeStory)
	stories.Get("/featured", sh.FeaturedStories)
	stories.Get("/:id", sh.FindStory)
	stories.Delete("/:id", sh.DeleteStory)
	stories.Get("/", middleware.IsLoggedIn, sh.FindAll)


	comments := api.Group("/comment")
	comments.Post("/", middleware.IsLoggedIn, ch.CreateComment)
	comments.Get("/", middleware.IsLoggedIn, ch.FindById)
	comments.Get("/story", middleware.IsLoggedIn, ch.FindAllCommentsByStoryId)
	comments.Put("/", middleware.IsLoggedIn, ch.UpdateById)
	comments.Delete("/", ch.DeleteById)
}

func Setup() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		ExposeHeaders: "Authorization",
	}))

	SetupRoutes(app)
	return app
}

