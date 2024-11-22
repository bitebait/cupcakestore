package services

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"strings"
)

type StoreConfigService interface {
	GetStoreConfig() (models.StoreConfig, error)
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

func (s *storeConfigService) GetStoreConfig() (models.StoreConfig, error) {
	config, err := s.storeConfigRepository.GetStoreConfig()

	if err != nil {
		err = errors.New("falha ao obter a configuração da loja")
	}

	return config, err
}

func (s *storeConfigService) Update(storeConfig *models.StoreConfig) error {
	s.normalizeStoreConfig(storeConfig)

	if err := s.storeConfigRepository.Update(storeConfig); err != nil {
		return errors.New("falha ao atualizar a configuração da loja")
	}

	return nil
}

func (s *storeConfigService) normalizeStoreConfig(storeConfig *models.StoreConfig) {
	storeConfig.PhysicalStoreAddress = strings.Title(storeConfig.PhysicalStoreAddress)
	storeConfig.PhysicalStoreEmail = strings.ToLower(storeConfig.PhysicalStoreEmail)
	storeConfig.PhysicalStoreCity = strings.Title(storeConfig.PhysicalStoreCity)
	storeConfig.PhysicalStoreState = strings.Title(storeConfig.PhysicalStoreState)
}
