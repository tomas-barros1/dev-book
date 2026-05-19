package main

import (
	"dev-book-api/config"
	"dev-book-api/database"
	"dev-book-api/handlers"
	"dev-book-api/repositories"
	"dev-book-api/routes"
	"dev-book-api/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	likeRepo := repositories.NewLikeRepository(db)

	userService := services.NewUserService(userRepo, followRepo)
	postService := services.NewPostService(postRepo, likeRepo, userRepo)

	handler := &handlers.Handler{
		UserService: userService,
		PostService: postService,
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwt_secret", cfg.JWTSecret)
		c.Next()
	})

	routes.RegisterRoutes(handler, r, cfg.JWTSecret)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
