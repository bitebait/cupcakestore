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
	if stock.Quantity <= 0 {
		return errors.New("quantidade inválida")
	}

	if stock.ProfileID == 0 {
		return errors.New("ProfileID não fornecido")
	}

	return s.stockRepository.Create(stock)
}

func (s *stockService) GetTotalStockQuantity(productID uint) (int, error) {
	return s.stockRepository.SumProductStockQuantity(productID)
}

func (s *stockService) FindByProductId(filter *models.StockFilter) []models.Stock {
	return s.stockRepository.FindByProductId(filter)
}
