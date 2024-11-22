package services

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"strings"
)

type ProductService interface {
	Create(product *models.Product) error
	FindAll(filter *models.ProductFilter) []models.Product
	FindActiveWithStock(filter *models.ProductFilter) []models.Product
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
	s.normalizeProduct(product)

	if err := s.productRepository.Create(product); err != nil {
		return errors.New("falha ao cadastrar o produto")
	}

	return nil
}

func (s *productService) FindAll(filter *models.ProductFilter) []models.Product {
	return s.productRepository.FindAll(filter)
}

func (s *productService) FindActiveWithStock(filter *models.ProductFilter) []models.Product {
	return s.productRepository.FindActiveWithStock(filter)
}

func (s *productService) FindById(id uint) (models.Product, error) {
	product, err := s.productRepository.FindById(id)

	if err != nil {
		err = errors.New("falha ao encontrar o produto com o id informado")
	}

	return product, err
}

func (s *productService) Update(product *models.Product) error {
	s.normalizeProduct(product)

	if err := s.productRepository.Update(product); err != nil {
		return errors.New("falha ao atualizar o produto")
	}

	return nil
}

func (s *productService) Delete(id uint) error {
	product, err := s.FindById(id)

	if err != nil {
		return errors.New("falha ao encontrar o produto com o id informado")
	}

	if err := s.productRepository.Delete(&product); err != nil {
		return errors.New("falha ao deletar o produto")
	}

	return nil
}

func (s *productService) normalizeProduct(product *models.Product) {
	product.Name = strings.Title(product.Name)
}
