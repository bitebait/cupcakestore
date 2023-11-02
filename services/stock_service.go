package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type StockService interface {
	Create(stock *models.Stock) error
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
	return s.stockRepository.Create(stock)
}
