package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"time"
)

type DashboardRepository interface {
	GetInfo(lastNDays int) *models.Dashboard
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(database *gorm.DB) DashboardRepository {
	return &dashboardRepository{
		db: database,
	}
}

func (r dashboardRepository) GetInfo(lastNDays int) *models.Dashboard {
	var newOrders, sales, users, products int64

	newOrdersStatuses := []models.ShoppingCartStatus{
		models.ActiveStatus,
		models.AwaitingPaymentStatus,
	}

	salesStatuses := []models.ShoppingCartStatus{
		models.PaymentApprovedStatus,
		models.DeliveredStatusDelivered,
		models.DeliveredStatusAwaiting,
		models.ProcessingStatus,
	}

	lastDays := time.Now().AddDate(0, 0, -lastNDays)

	r.db.Model(&models.Order{}).Where("status IN (?) AND created_at >= ?", newOrdersStatuses, lastDays).Count(&newOrders)
	r.db.Model(&models.Order{}).Where("status IN (?) AND created_at >= ?", salesStatuses, lastDays).Count(&sales)
	r.db.Model(&models.User{}).Where("is_staff = ? AND created_at >= ?", false, lastDays).Count(&users)
	r.db.Model(&models.Product{}).Where("is_active = ? AND created_at >= ?", true, lastDays).Count(&products)

	dashboard := &models.Dashboard{
		NewOrders: newOrders,
		Sales:     sales,
		Users:     users,
		Products:  products,
	}

	return dashboard
}
