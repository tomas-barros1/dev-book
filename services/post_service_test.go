package services

import (
	"dev-book-api/dtos"
	"dev-book-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockPostRepo struct {
	mock.Mock
}

func (m *mockPostRepo) Create(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *mockPostRepo) FindByID(id uint) (*models.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostRepo) FindByUserID(userID uint) ([]models.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *mockPostRepo) FindFeed(userID uint) ([]models.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *mockPostRepo) Update(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *mockPostRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockLikeRepo struct {
	mock.Mock
}

func (m *mockLikeRepo) Create(like *models.Like) error {
	args := m.Called(like)
	return args.Error(0)
}

func (m *mockLikeRepo) Delete(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *mockLikeRepo) HasLiked(userID, postID uint) (bool, error) {
	args := m.Called(userID, postID)
	return args.Bool(0), args.Error(1)
}

func TestPostService_Create_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	postRepo.On("Create", mock.AnythingOfType("*models.Post")).Return(nil)
	postRepo.On("FindByID", uint(0)).Return(&models.Post{
		Title:  "Test Title",
		UserID: 1,
		User:   models.User{Username: "author"},
	}, nil)

	post, err := svc.Create(dtos.CreatePostRequest{Title: "Test Title", Content: "Content"}, 1)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Title", post.Title)
	postRepo.AssertExpectations(t)
}

func TestPostService_GetByID_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{
		Title:  "Test",
		UserID: 1,
		User:   models.User{Username: "author"},
	}
	postRepo.On("FindByID", uint(1)).Return(post, nil)

	result, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Title)
	postRepo.AssertExpectations(t)
}

func TestPostService_GetByID_NotFound(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	postRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostService_GetFeed_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	posts := []models.Post{{Title: "Post1"}, {Title: "Post2"}}
	postRepo.On("FindFeed", uint(1)).Return(posts, nil)

	result, err := svc.GetFeed(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	postRepo.AssertExpectations(t)
}

func TestPostService_GetUserPosts_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	userRepo.On("FindByID", uint(1)).Return(&models.User{}, nil)
	posts := []models.Post{{Title: "Post1", UserID: 1}}
	postRepo.On("FindByUserID", uint(1)).Return(posts, nil)

	result, err := svc.GetUserPosts(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	postRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestPostService_GetUserPosts_UserNotFound(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	userRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.GetUserPosts(999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, result)
}

func TestPostService_Update_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Update", mock.AnythingOfType("*models.Post")).Return(nil)

	result, err := svc.Update(1, dtos.UpdatePostRequest{Title: "New Title", Content: "New Content"}, 1)

	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Content", result.Content)
	postRepo.AssertExpectations(t)
}

func TestPostService_Update_Unauthorized(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Title", UserID: 1}
	postRepo.On("FindByID", uint(1)).Return(post, nil)

	result, err := svc.Update(1, dtos.UpdatePostRequest{Title: "New Title"}, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
	assert.Nil(t, result)
}

func TestPostService_Update_NotFound(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	postRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.Update(999, dtos.UpdatePostRequest{Title: "New Title"}, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, result)
}

func TestPostService_Delete_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Title", UserID: 1}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)

	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostService_Delete_Unauthorized(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Title", UserID: 1}
	postRepo.On("FindByID", uint(1)).Return(post, nil)

	err := svc.Delete(1, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestPostService_Like_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Title", UserID: 2}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	likeRepo.On("HasLiked", uint(1), uint(1)).Return(false, nil)
	likeRepo.On("Create", mock.AnythingOfType("*models.Like")).Return(nil)

	err := svc.Like(1, 1)

	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestPostService_Like_AlreadyLiked(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	post := &models.Post{Title: "Title", UserID: 2}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	likeRepo.On("HasLiked", uint(1), uint(1)).Return(true, nil)

	err := svc.Like(1, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrAlreadyLiked, err)
}

func TestPostService_Like_PostNotFound(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	postRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	err := svc.Like(1, 999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestPostService_Unlike_Success(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	likeRepo.On("HasLiked", uint(1), uint(1)).Return(true, nil)
	likeRepo.On("Delete", uint(1), uint(1)).Return(nil)

	err := svc.Unlike(1, 1)

	assert.NoError(t, err)
	likeRepo.AssertExpectations(t)
}

func TestPostService_Unlike_NotLiked(t *testing.T) {
	postRepo := new(mockPostRepo)
	likeRepo := new(mockLikeRepo)
	userRepo := new(mockUserRepo)
	svc := NewPostService(postRepo, likeRepo, userRepo)

	likeRepo.On("HasLiked", uint(1), uint(1)).Return(false, nil)

	err := svc.Unlike(1, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrNotLiked, err)
}
