package main

import (
	"log"

	"dev-book-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal("Could not connect to database!")
	}

	r.Run()
}
