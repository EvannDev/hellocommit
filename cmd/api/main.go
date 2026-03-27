package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/hellocommit/api/internal/database"
	"github.com/hellocommit/api/internal/handlers"
	"github.com/hellocommit/api/internal/middleware"
	"github.com/hellocommit/api/internal/repositories"
	"github.com/hellocommit/api/internal/services"
	"github.com/hellocommit/api/pkg/github"
)

func main() {
	db, err := database.NewSQLite("hellocommit.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	githubClient := github.NewClient(githubToken)

	userRepo := repositories.NewUserRepository(db)
	repoRepo := repositories.NewRepoRepository(db)
	issueRepo := repositories.NewIssueRepository(db)

	userService := services.NewUserService(userRepo, githubClient)
	starredService := services.NewStarredService(userRepo, repoRepo, githubClient)
	issueService := services.NewIssueService(issueRepo, repoRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	starredHandler := handlers.NewStarredHandler(starredService, issueService)
	issueHandler := handlers.NewIssueHandler(issueService)
	feedHandler := handlers.NewFeedHandler(issueService)

	app := fiber.New(fiber.Config{
		AppName: "HelloCommit API",
	})

	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{allowedOrigin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Registration — no auth required
	api.Post("/users", userHandler.Create)

	// All routes below require a valid Bearer token
	authMiddleware := middleware.Auth(userRepo)

	users := api.Group("/users", authMiddleware)
	users.Get("/:id", userHandler.Get)
	users.Delete("/:id", userHandler.Delete)
	users.Get("/:id/starred", starredHandler.GetStarred)
	users.Post("/:id/sync", starredHandler.SyncStarred)

	repos := api.Group("/repos", authMiddleware)
	repos.Get("/:owner/:name/issues", issueHandler.GetIssues)
	repos.Post("/:owner/:name/sync", issueHandler.SyncIssues)

	issues := api.Group("/issues", authMiddleware)
	issues.Get("/good-first", issueHandler.GetGoodFirstIssues)
	issues.Post("/sync-all/:userId", issueHandler.SyncAllIssues)
	issues.Get("/:id", issueHandler.GetIssue)
	issues.Post("/:issueId/dismiss", issueHandler.DismissIssue)
	issues.Post("/:issueId/bookmark", issueHandler.ToggleBookmark)

	feeds := api.Group("/feeds", authMiddleware)
	feeds.Get("/rss", feedHandler.GetRSS)

	sync := api.Group("/sync", authMiddleware)
	sync.Post("/all/:userId", starredHandler.SyncAll)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
