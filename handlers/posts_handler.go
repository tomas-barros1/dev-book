package handlers

import (
	"net/http"

	"dev-book-api/models"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetPosts(c *gin.Context) {
	var posts []models.Post

	result := h.DB.Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
