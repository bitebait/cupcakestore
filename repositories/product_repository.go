package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(filter *models.ProductFilter) []models.Product
	FindActiveWithStock(filter *models.ProductFilter) []models.Product
	FindById(id uint) (models.Product, error)
	Update(product *models.Product) error
	Delete(product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		log.Errorf("ProductRepository Create: %s", err.Error())
		return err
	}

	return nil
}

func (r *productRepository) FindAll(filter *models.ProductFilter) []models.Product {
	return r.findProducts(filter, "")
}

func (r *productRepository) FindActiveWithStock(filter *models.ProductFilter) []models.Product {
	return r.findProducts(filter, "is_active = 1 AND current_stock > 0")
}

func (r *productRepository) findProducts(filter *models.ProductFilter, additionalCondition string) []models.Product {
	query := r.db.Model(&models.Product{})

	if additionalCondition != "" {
		query = query.Where(additionalCondition)
	}

	if filter.Product.Name != "" {
		filterPattern := "%" + filter.Product.Name + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", filterPattern, filterPattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("ProductRepository findProducts: %s", err.Error())
		return nil
	}
	filter.Pagination.Total = total

	var products []models.Product
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit
	query.Offset(offset).Limit(filter.Pagination.Limit).Order("created_at desc").Find(&products)

	return products
}

func (r *productRepository) FindById(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error

	if err != nil {
		log.Errorf("ProductRepository FindOrCreateById: %s", err.Error())
	}

	return product, err
}

func (r *productRepository) Update(product *models.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		log.Errorf("ProductRepository Update: %s", err.Error())
		return err
	}

	return nil
}

func (r *productRepository) Delete(product *models.Product) error {
	if err := r.db.Delete(product).Error; err != nil {
		log.Errorf("ProductRepository Delete: %s", err.Error())
		return err
	}

	return nil
}
