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

func (r *dashboardRepository) GetInfo(lastNDays int) *models.Dashboard {
	var dashboard models.Dashboard

	lastDays := time.Now().AddDate(0, 0, -lastNDays).Format("2006-01-02")

	r.db.Model(&models.Order{}).Where("status IN ?", []models.ShoppingCartStatus{
		models.ActiveStatus,
		models.AwaitingPaymentStatus,
	}).Where("DATE(created_at) >= ?", lastDays).Count(&dashboard.NewOrders)

	r.db.Model(&models.Order{}).Where("status IN ?", []models.ShoppingCartStatus{
		models.PaymentApprovedStatus,
		models.DeliveredStatusDelivered,
		models.DeliveredStatusAwaiting,
		models.ProcessingStatus,
	}).Where("DATE(created_at) >= ?", lastDays).Count(&dashboard.Sales)

	r.db.Model(&models.User{}).Count(&dashboard.Users)

	r.db.Model(&models.Product{}).Count(&dashboard.Products)

	return &dashboard
}
