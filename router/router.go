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
	lh := handlers.LikeHandler{LikeService: services.NewLikeService(repo.NewLikeRepoImpl())}
	dh := handlers.DisLikeHandler{DisLikeService: services.NewDisLikeService(repo.NewDisLikeRepoImpl())}
	ch := handlers.CommentHandler{CommentService: services.NewCommentService(repo.NewCommentRepoImpl())}
	sh := handlers.StoryHandler{StoryService: services.NewStoryService(repo.NewStoryRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	stories := api.Group("/stories")
	stories.Post("/", sh.CreateStory)
	stories.Put("/:id", sh.UpdateStory)
	stories.Get("/", middleware.IsLoggedIn, sh.FindAll)
	stories.Delete("/:id", sh.DeleteStory)

	likes := api.Group("/likes")
	likes.Post("/story", middleware.IsLoggedIn, lh.CreateLikeForStory)
	likes.Post("/comment", middleware.IsLoggedIn, lh.CreateLikeForComment)
	likes.Delete("/", middleware.IsLoggedIn, lh.DeleteLikeByUsername)

	disLikes := api.Group("/dislike")
	disLikes.Post("/story", middleware.IsLoggedIn, dh.CreateDisLikeForStory)
	disLikes.Post("/comment", middleware.IsLoggedIn, dh.CreateDisLikeForComment)
	disLikes.Delete("/", middleware.IsLoggedIn, dh.DeleteDisLikeByUsername)

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

