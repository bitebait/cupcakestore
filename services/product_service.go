package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type ProductService interface {
	Create(product *models.Product) error
	FindAll(p *models.Pagination, filter string) []models.Product
	FindById(id uint) (models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type productService struct {
	productRepository repositories.ProductRepository
}

func NewProductService(productRepository repositories.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) Create(product *models.Product) error {
	return s.productRepository.Create(product)
}

func (s *productService) FindAll(p *models.Pagination, filter string) []models.Product {
	return s.productRepository.FindAll(p, filter)
}

func (s *productService) FindById(id uint) (models.Product, error) {
	return s.productRepository.FindById(id)
}

func (s *productService) Update(product *models.Product) error {
	return s.productRepository.Update(product)
}

func (s *productService) Delete(id uint) error {
	product, err := s.FindById(id)
	if err != nil {
		return err
	}
	return s.productRepository.Delete(&product)
}
