package models

import "time"

type Like struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement:false"`
	PostID    uint      `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
}
