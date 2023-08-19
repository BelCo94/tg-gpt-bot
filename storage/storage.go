package storage

import (
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func (storage *Storage) InitModels() {
	storage.DB.AutoMigrate(&User{}, &Chat{}, &Message{})
}
