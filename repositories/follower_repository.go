package repositories

import (
	"dev-book-api/models"

	"gorm.io/gorm"
)

type FollowRepository interface {
	Create(follower *models.Follower) error
	Delete(userID, followerID uint) error
	GetFollowers(userID uint) ([]models.User, error)
	GetFollowing(userID uint) ([]models.User, error)
	IsFollowing(followerID, followingID uint) (bool, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follower *models.Follower) error {
	return r.db.Create(follower).Error
}

func (r *followRepository) Delete(userID, followerID uint) error {
	return r.db.Where("user_id = ? AND follower_id = ?", userID, followerID).
		Delete(&models.Follower{}).Error
}

func (r *followRepository) GetFollowers(userID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.
		Joins("JOIN followers ON followers.follower_id = users.id").
		Where("followers.user_id = ?", userID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *followRepository) GetFollowing(userID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.
		Joins("JOIN followers ON followers.user_id = users.id").
		Where("followers.follower_id = ?", userID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *followRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Follower{}).
		Where("user_id = ? AND follower_id = ?", followingID, followerID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
