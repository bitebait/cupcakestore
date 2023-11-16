package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type StoreConfigService interface {
	GetStoreConfig() (*models.StoreConfig, error)
	Update(storeConfig *models.StoreConfig) error
}

type storeConfigService struct {
	storeConfigRepository repositories.StoreConfigRepository
}

func NewStoreConfigService(storeConfigRepository repositories.StoreConfigRepository) StoreConfigService {
	return &storeConfigService{
		storeConfigRepository: storeConfigRepository,
	}
}

func (s storeConfigService) GetStoreConfig() (*models.StoreConfig, error) {
	return s.storeConfigRepository.GetStoreConfig()
}

func (s *storeConfigService) Update(storeConfig *models.StoreConfig) error {
	return s.storeConfigRepository.Update(storeConfig)
}
