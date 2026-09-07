package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log/slog"
)

var ErrStoreConfigChanged = errors.New("a configuração mudou enquanto você editava; recarregue a página e tente novamente")

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
	if storeConfig == nil || storeConfig.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	result := r.db.Model(storeConfig).Where("updated_at = ?", storeConfig.UpdatedAt).
		Select("DeliveryPrice", "DeliveryIsActive", "PhysicalStoreEmail", "PhysicalStoreAddress",
			"PhysicalStoreCity", "PhysicalStoreState", "PhysicalStorePostalCode", "PhysicalStorePhoneNumber",
			"PaymentCashIsActive", "PaymentPixIsActive", "PixKey", "PixKeyType", "PixReceiverName").Updates(storeConfig)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrStoreConfigChanged
	}

	return nil
}
