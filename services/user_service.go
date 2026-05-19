package services

import (
	"dev-book-api/dtos"
	"dev-book-api/models"
	"dev-book-api/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req dtos.CreateUserRequest) (*models.User, error)
	Login(req dtos.LoginRequest) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Search(query string) ([]models.User, error)
	Update(id uint, req dtos.UpdateUserRequest, actorID uint) (*models.User, error)
	Delete(id uint, actorID uint) error
	Follow(followerID, followingID uint) error
	Unfollow(followerID, followingID uint) error
	GetFollowers(userID uint) ([]models.User, error)
	GetFollowing(userID uint) ([]models.User, error)
	UpdatePassword(id uint, req dtos.UpdatePasswordRequest, actorID uint) error
}

type userService struct {
	userRepo   repositories.UserRepository
	followRepo repositories.FollowRepository
}

func NewUserService(userRepo repositories.UserRepository, followRepo repositories.FollowRepository) UserService {
	return &userService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *userService) Create(req dtos.CreateUserRequest) (*models.User, error) {
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, ErrDuplicateUsername
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req dtos.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Search(query string) ([]models.User, error) {
	return s.userRepo.Search(query)
}

func (s *userService) Update(id uint, req dtos.UpdateUserRequest, actorID uint) (*models.User, error) {
	if id != actorID {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if req.Username != "" {
		existing, _ := s.userRepo.FindByUsername(req.Username)
		if existing != nil && existing.ID != id {
			return nil, ErrDuplicateUsername
		}
		user.Username = req.Username
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint, actorID uint) error {
	if id != actorID {
		return ErrUnauthorized
	}

	if _, err := s.userRepo.FindByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	return s.userRepo.Delete(id)
}

func (s *userService) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}

	if _, err := s.userRepo.FindByID(followingID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	already, err := s.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if already {
		return ErrAlreadyFollowing
	}

	return s.followRepo.Create(&models.Follower{
		UserID:     followingID,
		FollowerID: followerID,
	})
}

func (s *userService) Unfollow(followerID, followingID uint) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}

	already, err := s.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if !already {
		return ErrNotFollowing
	}

	return s.followRepo.Delete(followingID, followerID)
}

func (s *userService) GetFollowers(userID uint) ([]models.User, error) {
	return s.followRepo.GetFollowers(userID)
}

func (s *userService) GetFollowing(userID uint) ([]models.User, error) {
	return s.followRepo.GetFollowing(userID)
}

func (s *userService) UpdatePassword(id uint, req dtos.UpdatePasswordRequest, actorID uint) error {
	if id != actorID {
		return ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}
