package services

import (
	"errors"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type StockService interface {
	Create(stock *models.Stock) error
	GetTotalStockQuantity(productID uint) (int, error)
	FindByProductId(filter *models.StockFilter) []models.Stock
}

type stockService struct {
	stockRepository repositories.StockRepository
}

func NewStockService(stockRepository repositories.StockRepository) StockService {
	return &stockService{
		stockRepository: stockRepository,
	}
}

func (s *stockService) Create(stock *models.Stock) error {
	if stock == nil {
		return errors.New("o estoque deve ser informado")
	}
	if err := stock.Validate(); err != nil {
		return err
	}

	if err := s.stockRepository.Create(stock); err != nil {
		return errors.New("falha ao criar o estoque do produto")
	}

	return nil
}

func (s *stockService) GetTotalStockQuantity(productID uint) (int, error) {
	total, err := s.stockRepository.SumProductStockQuantity(productID)

	if err != nil {
		return 0, errors.New("falha ao obter a quantidade de estoque do produto")
	}

	return total, nil
}

func (s *stockService) FindByProductId(filter *models.StockFilter) []models.Stock {
	return s.stockRepository.FindByProductId(filter)
}
