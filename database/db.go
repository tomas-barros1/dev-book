package database

import (
	"dev-book-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/dev_book?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&models.Post{}, &models.User{})

	return db, err
}
