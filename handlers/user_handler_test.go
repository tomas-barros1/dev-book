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
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Create(req dtos.CreateUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) Login(req dtos.LoginRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) Search(query string) ([]models.User, error) {
	args := m.Called(query)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserService) Update(id uint, req dtos.UpdateUserRequest, actorID uint) (*models.User, error) {
	args := m.Called(id, req, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) Delete(id uint, actorID uint) error {
	args := m.Called(id, actorID)
	return args.Error(0)
}

func (m *mockUserService) Follow(followerID, followingID uint) error {
	args := m.Called(followerID, followingID)
	return args.Error(0)
}

func (m *mockUserService) Unfollow(followerID, followingID uint) error {
	args := m.Called(followerID, followingID)
	return args.Error(0)
}

func (m *mockUserService) GetFollowers(userID uint) ([]models.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserService) GetFollowing(userID uint) ([]models.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserService) UpdatePassword(id uint, req dtos.UpdatePasswordRequest, actorID uint) error {
	args := m.Called(id, req, actorID)
	return args.Error(0)
}

type mockPostService struct {
	mock.Mock
}

func (m *mockPostService) Create(req dtos.CreatePostRequest, userID uint) (*models.Post, error) {
	args := m.Called(req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetFeed(userID uint) ([]models.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *mockPostService) GetUserPosts(userID uint) ([]models.Post, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *mockPostService) Update(id uint, req dtos.UpdatePostRequest, actorID uint) (*models.Post, error) {
	args := m.Called(id, req, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) Delete(id uint, actorID uint) error {
	args := m.Called(id, actorID)
	return args.Error(0)
}

func (m *mockPostService) Like(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *mockPostService) Unlike(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func setupHandlerTest(userSvc *mockUserService, postSvc *mockPostService) (*Handler, *gin.Engine) {
	if userSvc == nil {
		userSvc = new(mockUserService)
	}
	if postSvc == nil {
		postSvc = new(mockPostService)
	}

	h := &Handler{
		UserService: userSvc,
		PostService: postSvc,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwt_secret", "test-secret")
		c.Next()
	})

	return h, r
}

func authenticatedRequest(c *gin.Context, userID uint) {
	c.Set("user_id", userID)
	c.Set("username", "testuser")
}

func TestCreateUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users", h.CreateUser)

	createdUser := &models.User{Username: "newuser"}
	createdUser.ID = 1

	mockSvc.On("Create", dtos.CreateUserRequest{Username: "newuser", Password: "password123"}).Return(createdUser, nil)

	body, _ := json.Marshal(map[string]string{"username": "newuser", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateUser_Duplicate(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users", h.CreateUser)

	mockSvc.On("Create", dtos.CreateUserRequest{Username: "existing", Password: "password123"}).Return(nil, services.ErrDuplicateUsername)

	body, _ := json.Marshal(map[string]string{"username": "existing", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateUser_InvalidBody(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users", h.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{"invalid`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/login", h.LoginUser)

	user := &models.User{Username: "testuser"}
	user.ID = 1

	mockSvc.On("Login", dtos.LoginRequest{Username: "testuser", Password: "password123"}).Return(user, nil)

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/login", h.LoginUser)

	mockSvc.On("Login", dtos.LoginRequest{Username: "testuser", Password: "wrongpass"}).Return(nil, services.ErrInvalidCredentials)

	body, _ := json.Marshal(map[string]string{"username": "testuser", "password": "wrongpass"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestSearchUsers(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.GET("/users", h.SearchUsers)

	users := []models.User{{Username: "joao"}, {Username: "joana"}}
	mockSvc.On("Search", "joa").Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users?query=joa", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.GET("/users/:userId", h.GetUser)

	user := &models.User{Username: "testuser"}
	user.ID = 1

	mockSvc.On("GetByID", uint(1)).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.GET("/users/:userId", h.GetUser)

	mockSvc.On("GetByID", uint(999)).Return(nil, services.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.PUT("/users/:userId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UpdateUser(c)
	})

	updatedUser := &models.User{Username: "newname"}
	updatedUser.ID = 1

	mockSvc.On("Update", uint(1), dtos.UpdateUserRequest{Username: "newname"}, uint(1)).Return(updatedUser, nil)

	body, _ := json.Marshal(map[string]string{"username": "newname"})
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.PUT("/users/:userId", func(c *gin.Context) {
		authenticatedRequest(c, 2)
		h.UpdateUser(c)
	})

	mockSvc.On("Update", uint(1), dtos.UpdateUserRequest{Username: "newname"}, uint(2)).Return(nil, services.ErrUnauthorized)

	body, _ := json.Marshal(map[string]string{"username": "newname"})
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.DELETE("/users/:userId", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.DeleteUser(c)
	})

	mockSvc.On("Delete", uint(1), uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFollowUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/follow", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.FollowUser(c)
	})

	mockSvc.On("Follow", uint(1), uint(2)).Return(nil)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/users/2/follow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFollowUser_AlreadyFollowing(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/follow", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.FollowUser(c)
	})

	mockSvc.On("Follow", uint(1), uint(2)).Return(services.ErrAlreadyFollowing)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/users/2/follow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUnfollowUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/unfollow", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UnfollowUser(c)
	})

	mockSvc.On("Unfollow", uint(1), uint(2)).Return(nil)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/users/2/unfollow", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetFollowers(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.GET("/users/:userId/followers", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetFollowers(c)
	})

	users := []models.User{{Username: "follower1"}}
	mockSvc.On("GetFollowers", uint(1)).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1/followers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetFollowing(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.GET("/users/:userId/following", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.GetFollowing(c)
	})

	users := []models.User{{Username: "following1"}}
	mockSvc.On("GetFollowing", uint(1)).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1/following", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdatePassword_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/update-password", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UpdatePassword(c)
	})

	mockSvc.On("UpdatePassword", uint(1), dtos.UpdatePasswordRequest{OldPassword: "oldpass", NewPassword: "newpass"}, uint(1)).Return(nil)

	body, _ := json.Marshal(map[string]string{"old_password": "oldpass", "new_password": "newpass"})
	req := httptest.NewRequest(http.MethodPost, "/users/1/update-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdatePassword_WrongOldPassword(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/update-password", func(c *gin.Context) {
		authenticatedRequest(c, 1)
		h.UpdatePassword(c)
	})

	mockSvc.On("UpdatePassword", uint(1), dtos.UpdatePasswordRequest{OldPassword: "wrong", NewPassword: "newpass"}, uint(1)).Return(services.ErrInvalidCredentials)

	body, _ := json.Marshal(map[string]string{"old_password": "wrong", "new_password": "newpass"})
	req := httptest.NewRequest(http.MethodPost, "/users/1/update-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdatePassword_Unauthorized(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.POST("/users/:userId/update-password", func(c *gin.Context) {
		authenticatedRequest(c, 2)
		h.UpdatePassword(c)
	})

	mockSvc.On("UpdatePassword", uint(1), mock.Anything, uint(2)).Return(services.ErrUnauthorized)

	body, _ := json.Marshal(map[string]string{"old_password": "oldpass", "new_password": "newpass"})
	req := httptest.NewRequest(http.MethodPost, "/users/1/update-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockSvc := new(mockUserService)
	h, r := setupHandlerTest(mockSvc, nil)

	r.DELETE("/users/:userId", func(c *gin.Context) {
		authenticatedRequest(c, 999)
		h.DeleteUser(c)
	})

	mockSvc.On("Delete", uint(999), uint(999)).Return(services.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}
