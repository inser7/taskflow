package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func Migrate(db *gorm.DB) error {
	// AutoMigrate создаст таблицу, если она не существует
	// и обновит схему, если есть различия
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return err
	}

	// Дополнительные операции миграции, если необходимо
	// Например, создание индексов, начальных записей и т.д.

	return nil
}
