package handlers

import (
	"bytes"
	"dev-book-api/dtos"
	"dev-book-api/models"
	"dev-book-api/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.CreatePost(c)
	})

	post := &models.Post{
		Title:  "Test Post",
		UserID: 1,
		User:   models.User{Username: "author"},
	}
	post.ID = 1

	postSvc.On("Create", dtos.CreatePostRequest{Title: "Test Post", Content: "Content"}, uint(1)).Return(post, nil)

	body, _ := json.Marshal(map[string]string{"title": "Test Post", "content": "Content"})
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	postSvc.AssertExpectations(t)
}

func TestCreatePost_InvalidBody(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.CreatePost(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFeed_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/posts", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetFeed(c)
	})

	posts := []models.Post{
		{Title: "Post1", UserID: 1, User: models.User{Username: "author1"}},
		{Title: "Post2", UserID: 2, User: models.User{Username: "author2"}},
	}

	postSvc.On("GetFeed", uint(1)).Return(posts, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestGetPost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetPost(c)
	})

	post := &models.Post{Title: "Post", UserID: 1, User: models.User{Username: "author"}}
	post.ID = 1

	postSvc.On("GetByID", uint(1)).Return(post, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestGetPost_NotFound(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetPost(c)
	})

	postSvc.On("GetByID", uint(999)).Return(nil, services.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/posts/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	postSvc.AssertExpectations(t)
}

func TestUpdatePost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.PUT("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UpdatePost(c)
	})

	post := &models.Post{Title: "Updated", Content: "Updated Content", UserID: 1, User: models.User{Username: "author"}}
	post.ID = 1

	postSvc.On("Update", uint(1), dtos.UpdatePostRequest{Title: "Updated", Content: "Updated Content"}, uint(1)).Return(post, nil)

	body, _ := json.Marshal(map[string]string{"title": "Updated", "content": "Updated Content"})
	req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestUpdatePost_Unauthorized(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.PUT("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 2)
		h.UpdatePost(c)
	})

	postSvc.On("Update", uint(1), dtos.UpdatePostRequest{Title: "Updated"}, uint(2)).Return(nil, services.ErrUnauthorized)

	body, _ := json.Marshal(map[string]string{"title": "Updated"})
	req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	postSvc.AssertExpectations(t)
}

func TestDeletePost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.DELETE("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.DeletePost(c)
	})

	postSvc.On("Delete", uint(1), uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	postSvc.AssertExpectations(t)
}

func TestDeletePost_Unauthorized(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.DELETE("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 2)
		h.DeletePost(c)
	})

	postSvc.On("Delete", uint(1), uint(2)).Return(services.ErrUnauthorized)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	postSvc.AssertExpectations(t)
}

func TestGetUserPosts_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/users/:userId/posts", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetUserPosts(c)
	})

	posts := []models.Post{{Title: "Post1", UserID: 1, User: models.User{Username: "author"}}}
	postSvc.On("GetUserPosts", uint(1)).Return(posts, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestGetUserPosts_NotFound(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/users/:userId/posts", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetUserPosts(c)
	})

	postSvc.On("GetUserPosts", uint(999)).Return(nil, services.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/users/999/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	postSvc.AssertExpectations(t)
}

func TestLikePost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts/:postId/like", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.LikePost(c)
	})

	postSvc.On("Like", uint(1), uint(1)).Return(nil)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/posts/1/like", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestLikePost_AlreadyLiked(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts/:postId/like", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.LikePost(c)
	})

	postSvc.On("Like", uint(1), uint(1)).Return(services.ErrAlreadyLiked)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/posts/1/like", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	postSvc.AssertExpectations(t)
}

func TestLikePost_NotFound(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts/:postId/like", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.LikePost(c)
	})

	postSvc.On("Like", uint(1), uint(999)).Return(services.ErrNotFound)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/posts/999/like", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	postSvc.AssertExpectations(t)
}

func TestUnlikePost_Success(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts/:postId/unlike", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UnlikePost(c)
	})

	postSvc.On("Unlike", uint(1), uint(1)).Return(nil)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/posts/1/unlike", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	postSvc.AssertExpectations(t)
}

func TestUnlikePost_NotLiked(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.POST("/posts/:postId/unlike", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UnlikePost(c)
	})

	postSvc.On("Unlike", uint(1), uint(1)).Return(services.ErrNotLiked)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/posts/1/unlike", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	postSvc.AssertExpectations(t)
}

func TestGetPost_InvalidID(t *testing.T) {
	userSvc := new(mockUserService)
	postSvc := new(mockPostService)
	h, r := setupHandlerTest(userSvc, postSvc)

	r.GET("/posts/:postId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetPost(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/posts/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
