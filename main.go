package main

import (
	"log"

	"dev-book-api/database"
	"dev-book-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal("Could not connect to database!")
	}

	handler := handlers.Handler{
		DB: db,
	}

	r.Run()
}
