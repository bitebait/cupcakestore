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
		return errors.New("quantidade deve ser maior que zero")
	}

	if stock.ProfileID == 0 {
		return errors.New("o id do perfil deve ser informado")
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
