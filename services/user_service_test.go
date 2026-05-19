package services

import (
	"dev-book-api/dtos"
	"dev-book-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepo) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) Search(query string) ([]models.User, error) {
	args := m.Called(query)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserRepo) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockFollowRepo struct {
	mock.Mock
}

func (m *mockFollowRepo) Create(follower *models.Follower) error {
	args := m.Called(follower)
	return args.Error(0)
}

func (m *mockFollowRepo) Delete(userID, followerID uint) error {
	args := m.Called(userID, followerID)
	return args.Error(0)
}

func (m *mockFollowRepo) GetFollowers(userID uint) ([]models.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockFollowRepo) GetFollowing(userID uint) ([]models.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockFollowRepo) IsFollowing(followerID, followingID uint) (bool, error) {
	args := m.Called(followerID, followingID)
	return args.Bool(0), args.Error(1)
}

func TestUserService_Create_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	userRepo.On("FindByUsername", "newuser").Return(nil, gorm.ErrRecordNotFound)
	userRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.Create(dtos.CreateUserRequest{Username: "newuser", Password: "password123"})

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.NotEmpty(t, user.Password)
	userRepo.AssertExpectations(t)
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	existing := &models.User{Username: "existing"}
	userRepo.On("FindByUsername", "existing").Return(existing, nil)

	user, err := svc.Create(dtos.CreateUserRequest{Username: "existing", Password: "password123"})

	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateUsername, err)
	assert.Nil(t, user)
	userRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{Username: "testuser", Password: string(hashed)}
	userRepo.On("FindByUsername", "testuser").Return(user, nil)

	result, err := svc.Login(dtos.LoginRequest{Username: "testuser", Password: "password123"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	userRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &models.User{Username: "testuser", Password: string(hashed)}
	userRepo.On("FindByUsername", "testuser").Return(user, nil)

	result, err := svc.Login(dtos.LoginRequest{Username: "testuser", Password: "wrongpassword"})

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	userRepo.On("FindByUsername", "nonexistent").Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.Login(dtos.LoginRequest{Username: "nonexistent", Password: "password123"})

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	user := &models.User{Username: "testuser"}
	userRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
	userRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	userRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

func TestUserService_Search(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	users := []models.User{{Username: "joao"}, {Username: "joana"}}
	userRepo.On("Search", "joa").Return(users, nil)

	result, err := svc.Search("joa")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	userRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	user := &models.User{Username: "oldname"}
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("FindByUsername", "newname").Return(nil, gorm.ErrRecordNotFound)
	userRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	result, err := svc.Update(1, dtos.UpdateUserRequest{Username: "newname"}, 1)

	assert.NoError(t, err)
	assert.Equal(t, "newname", result.Username)
	userRepo.AssertExpectations(t)
}

func TestUserService_Update_Unauthorized(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	result, err := svc.Update(1, dtos.UpdateUserRequest{Username: "newname"}, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
	assert.Nil(t, result)
}

func TestUserService_Delete_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	user := &models.User{Username: "testuser"}
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_Delete_Unauthorized(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	err := svc.Delete(1, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestUserService_Follow_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	user := &models.User{Username: "target"}
	userRepo.On("FindByID", uint(2)).Return(user, nil)
	followRepo.On("IsFollowing", uint(1), uint(2)).Return(false, nil)
	followRepo.On("Create", mock.AnythingOfType("*models.Follower")).Return(nil)

	err := svc.Follow(1, 2)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	followRepo.AssertExpectations(t)
}

func TestUserService_Follow_Self(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	err := svc.Follow(1, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrCannotFollowSelf, err)
}

func TestUserService_Follow_AlreadyFollowing(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	user := &models.User{Username: "target"}
	userRepo.On("FindByID", uint(2)).Return(user, nil)
	followRepo.On("IsFollowing", uint(1), uint(2)).Return(true, nil)

	err := svc.Follow(1, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrAlreadyFollowing, err)
}

func TestUserService_Unfollow_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	followRepo.On("IsFollowing", uint(1), uint(2)).Return(true, nil)
	followRepo.On("Delete", uint(2), uint(1)).Return(nil)

	err := svc.Unfollow(1, 2)

	assert.NoError(t, err)
	followRepo.AssertExpectations(t)
}

func TestUserService_Unfollow_NotFollowing(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	followRepo.On("IsFollowing", uint(1), uint(2)).Return(false, nil)

	err := svc.Unfollow(1, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFollowing, err)
}

func TestUserService_GetFollowers_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	followers := []models.User{{Username: "follower1"}}
	followRepo.On("GetFollowers", uint(1)).Return(followers, nil)

	result, err := svc.GetFollowers(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	followRepo.AssertExpectations(t)
}

func TestUserService_GetFollowing_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	following := []models.User{{Username: "following1"}}
	followRepo.On("GetFollowing", uint(1)).Return(following, nil)

	result, err := svc.GetFollowing(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	followRepo.AssertExpectations(t)
}

func TestUserService_UpdatePassword_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	user := &models.User{Username: "testuser", Password: string(hashed)}
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	err := svc.UpdatePassword(1, dtos.UpdatePasswordRequest{OldPassword: "oldpass", NewPassword: "newpass"}, 1)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_UpdatePassword_WrongOldPassword(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	user := &models.User{Username: "testuser", Password: string(hashed)}
	userRepo.On("FindByID", uint(1)).Return(user, nil)

	err := svc.UpdatePassword(1, dtos.UpdatePasswordRequest{OldPassword: "wrongold", NewPassword: "newpass"}, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestUserService_UpdatePassword_Unauthorized(t *testing.T) {
	userRepo := new(mockUserRepo)
	followRepo := new(mockFollowRepo)
	svc := NewUserService(userRepo, followRepo)

	err := svc.UpdatePassword(1, dtos.UpdatePasswordRequest{OldPassword: "oldpass", NewPassword: "newpass"}, 2)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}
