package models

import "time"

type Follower struct {
	UserID     uint      `gorm:"primaryKey;autoIncrement:false"`
	FollowerID uint      `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt  time.Time
}
