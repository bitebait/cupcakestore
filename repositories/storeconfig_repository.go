package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log/slog"
)

type StoreConfigRepository interface {
	GetStoreConfig() (models.StoreConfig, error)
	Update(storeConfig *models.StoreConfig) error
}

type storeConfigRepository struct {
	db *gorm.DB
}

func NewStoreConfigRepository(database *gorm.DB) StoreConfigRepository {
	return &storeConfigRepository{
		db: database,
	}
}

func (r *storeConfigRepository) GetStoreConfig() (models.StoreConfig, error) {
	var storeConfig models.StoreConfig
	err := r.db.First(&storeConfig).Error

	if err != nil {
		slog.Error("StoreConfigRepository GetStoreConfig", "error", err)
	}

	return storeConfig, err
}

func (r *storeConfigRepository) Update(storeConfig *models.StoreConfig) error {
	if err := r.db.Save(storeConfig).Error; err != nil {
		slog.Error("StoreConfigRepository Update", "error", err)
		return err
	}

	return nil
}
