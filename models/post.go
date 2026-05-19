package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title"`
	Content string `json:"content"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	User    User   `json:"-"`
}
