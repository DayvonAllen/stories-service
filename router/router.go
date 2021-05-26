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
	th := handlers.TagHandler{TagService: services.NewTagService(repo.NewTagRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	//stories := api.Group("/stories")
	////stories.Get("/")
	//fmt.Println(stories)

	tags := api.Group("/tags")
	tags.Get("/", middleware.IsLoggedIn, th.FindAll)
	tags.Get("/:tagName", middleware.IsLoggedIn, th.FindByTagName)
	tags.Post("/", middleware.IsLoggedIn, th.Create)
	tags.Delete("/", th.DeleteById)

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

