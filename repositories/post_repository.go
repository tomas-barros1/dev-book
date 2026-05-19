package repositories

import (
	"dev-book-api/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindByUserID(userID uint) ([]models.Post, error)
	FindFeed(userID uint) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindFeed(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.
		Joins("LEFT JOIN followers ON followers.user_id = posts.user_id").
		Where("followers.follower_id = ? OR posts.user_id = ?", userID, userID).
		Group("posts.id").
		Preload("User").
		Order("posts.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
