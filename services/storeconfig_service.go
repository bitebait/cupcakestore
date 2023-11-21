package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"strings"
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

func (s *storeConfigService) GetStoreConfig() (*models.StoreConfig, error) {
	return s.storeConfigRepository.GetStoreConfig()
}

func (s *storeConfigService) Update(storeConfig *models.StoreConfig) error {
	s.normalizeStoreConfig(storeConfig)
	return s.storeConfigRepository.Update(storeConfig)
}

func (s *storeConfigService) normalizeStoreConfig(storeConfig *models.StoreConfig) {
	storeConfig.PhysicalStoreAddress = strings.Title(storeConfig.PhysicalStoreAddress)
	storeConfig.PhysicalStoreEmail = strings.ToLower(storeConfig.PhysicalStoreEmail)
	storeConfig.PhysicalStoreCity = strings.Title(storeConfig.PhysicalStoreCity)
	storeConfig.PhysicalStoreState = strings.Title(storeConfig.PhysicalStoreState)
}
