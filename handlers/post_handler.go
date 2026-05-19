package handlers

import (
	"dev-book-api/dtos"
	"dev-book-api/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var req dtos.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	post, err := h.PostService.Create(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post": dtos.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			AuthorID:  post.UserID,
			Author:    post.User.Username,
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}

func (h *Handler) GetFeed(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	posts, err := h.PostService.GetFeed(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dtos.PostResponse, len(posts))
	for i, p := range posts {
		response[i] = dtos.PostResponse{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			AuthorID:  p.UserID,
			Author:    p.User.Username,
			CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"posts": response})
}

func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := h.PostService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": dtos.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			AuthorID:  post.UserID,
			Author:    post.User.Username,
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req dtos.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	actorID := c.MustGet("user_id").(uint)

	post, err := h.PostService.Update(uint(id), req, actorID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": dtos.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			AuthorID:  post.UserID,
			Author:    post.User.Username,
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}

func (h *Handler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	actorID := c.MustGet("user_id").(uint)

	err = h.PostService.Delete(uint(id), actorID)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetUserPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	posts, err := h.PostService.GetUserPosts(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dtos.PostResponse, len(posts))
	for i, p := range posts {
		response[i] = dtos.PostResponse{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			AuthorID:  p.UserID,
			Author:    p.User.Username,
			CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"posts": response})
}

func (h *Handler) LikePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	err = h.PostService.Like(userID, uint(id))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if errors.Is(err, services.ErrAlreadyLiked) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "liked successfully"})
}

func (h *Handler) UnlikePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	userID := c.MustGet("user_id").(uint)

	err = h.PostService.Unlike(userID, uint(id))
	if err != nil {
		if errors.Is(err, services.ErrNotLiked) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unliked successfully"})
}
