package database

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb() (*sql.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/dev_book?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return db, err
}
