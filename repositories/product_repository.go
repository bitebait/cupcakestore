package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(filter *models.ProductFilter) []models.Product
	FindById(id uint) (models.Product, error)
	Update(product *models.Product) error
	Delete(product *models.Product) error
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
func (r *productRepository) FindAll(filter *models.ProductFilter) []models.Product {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.Product{})

	if filter.Product.Name != "" {
		filterPattern := "%" + filter.Product.Name + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", filterPattern, filterPattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var products []models.Product
	query.Offset(offset).Limit(filter.Pagination.Limit).Order("name").Find(&products)
	return products
}

func (r *productRepository) FindById(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return product, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *models.Product) error {
	return r.db.Delete(product).Error
}
