package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
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
		log.Errorf("StoreConfigRepository GetStoreConfig: %s", err.Error())
	}

	return storeConfig, err
}

func (r *storeConfigRepository) Update(storeConfig *models.StoreConfig) error {
	if err := r.db.Save(storeConfig).Error; err != nil {
		log.Errorf("StoreConfigRepository Update: %s", err.Error())
		return err
	}

	return nil
}
