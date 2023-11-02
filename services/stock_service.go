package services

import (
	"errors"
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
	if stock.Quantity <= 0 {
		return errors.New("quantidade inválida")
	}

	return s.stockRepository.Create(stock)
}
