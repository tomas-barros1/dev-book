package repositories

import (
	"dev-book-api/models"

	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, postID uint) error
	HasLiked(userID, postID uint) (bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.Like{}).Error
}

func (r *likeRepository) HasLiked(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
