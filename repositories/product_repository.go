package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(p *models.Pagination, filter string) []models.Product
	FindById(id uint) (models.Product, error)
	Delete(user *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &productRepository{
		db: database,
	}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(p *models.Pagination, filter string) []models.Product {
	offset := (p.Page - 1) * p.Limit

	query := r.db.Model(&models.Product{})

	if filter != "" {
		filterPattern := "%" + filter + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", filterPattern, filterPattern)
	}

	query.Count(&p.Total)

	var products []models.Product
	query.Offset(offset).Limit(p.Limit).Order("name").Find(&products)

	return products
}

func (r *productRepository) FindById(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return product, err
}

func (r *productRepository) Delete(user *models.Product) error {
	return r.db.Delete(user).Error
}
