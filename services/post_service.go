package services

import (
	"dev-book-api/cache"
	"dev-book-api/dtos"
	"dev-book-api/models"
	"dev-book-api/repositories"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	postCacheTTL = 5 * time.Minute
	postCacheKey = "post:%d"
)

type PostService interface {
	Create(req dtos.CreatePostRequest, userID uint) (*models.Post, error)
	GetByID(id uint) (*models.Post, error)
	GetFeed(userID uint) ([]models.Post, error)
	GetUserPosts(userID uint) ([]models.Post, error)
	Update(id uint, req dtos.UpdatePostRequest, actorID uint) (*models.Post, error)
	Delete(id uint, actorID uint) error
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
}

type postService struct {
	postRepo repositories.PostRepository
	likeRepo repositories.LikeRepository
	userRepo repositories.UserRepository
	cache    cache.Cache
}

func NewPostService(postRepo repositories.PostRepository, likeRepo repositories.LikeRepository, userRepo repositories.UserRepository, c cache.Cache) PostService {
	return &postService{
		postRepo: postRepo,
		likeRepo: likeRepo,
		userRepo: userRepo,
		cache:    c,
	}
}

func (s *postService) Create(req dtos.CreatePostRequest, userID uint) (*models.Post, error) {
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	return s.postRepo.FindByID(post.ID)
}

func (s *postService) GetByID(id uint) (*models.Post, error) {
	key := fmt.Sprintf(postCacheKey, id)
	var post models.Post
	if err := s.cache.Get(key, &post); err == nil {
		return &post, nil
	}

	postDB, err := s.postRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	s.cache.Set(key, postDB, postCacheTTL)
	return postDB, nil
}

func (s *postService) GetFeed(userID uint) ([]models.Post, error) {
	return s.postRepo.FindFeed(userID)
}

func (s *postService) GetUserPosts(userID uint) ([]models.Post, error) {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.postRepo.FindByUserID(userID)
}

func (s *postService) Update(id uint, req dtos.UpdatePostRequest, actorID uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if post.UserID != actorID {
		return nil, ErrUnauthorized
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	s.cache.Del(fmt.Sprintf(postCacheKey, id))
	return post, nil
}

func (s *postService) Delete(id uint, actorID uint) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	if post.UserID != actorID {
		return ErrUnauthorized
	}

	s.cache.Del(fmt.Sprintf(postCacheKey, id))
	return s.postRepo.Delete(id)
}

func (s *postService) Like(userID, postID uint) error {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	already, err := s.likeRepo.HasLiked(userID, postID)
	if err != nil {
		return err
	}
	if already {
		return ErrAlreadyLiked
	}

	return s.likeRepo.Create(&models.Like{
		UserID: userID,
		PostID: postID,
	})
}

func (s *postService) Unlike(userID, postID uint) error {
	already, err := s.likeRepo.HasLiked(userID, postID)
	if err != nil {
		return err
	}
	if !already {
		return ErrNotLiked
	}

	return s.likeRepo.Delete(userID, postID)
}
