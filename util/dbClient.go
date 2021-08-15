package util

import "gorm.io/gorm"

type dbClient struct {
	DB *gorm.DB
}

func NewGormDB()
