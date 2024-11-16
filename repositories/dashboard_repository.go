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
	const dateFormat = "2006-01-02"
	var dashboard models.Dashboard

	lastDays := time.Now().AddDate(0, 0, -lastNDays).Format(dateFormat)

	dashboard.NewOrders = r.countOrderStatusInLastDays(lastDays, []models.ShoppingCartStatus{
		models.ActiveStatus,
		models.AwaitingPaymentStatus,
	})
	dashboard.Sales = r.countOrderStatusInLastDays(lastDays, []models.ShoppingCartStatus{
		models.PaymentApprovedStatus,
		models.DeliveredStatusDelivered,
		models.DeliveredStatusAwaiting,
		models.ProcessingStatus,
	})
	dashboard.Users = r.countRecords(&models.User{})
	dashboard.Products = r.countRecords(&models.Product{})

	return &dashboard
}

func (r *dashboardRepository) countOrderStatusInLastDays(lastDays string, statuses []models.ShoppingCartStatus) int64 {
	var count int64
	r.db.Model(&models.Order{}).
		Where("status IN ?", statuses).
		Where("DATE(created_at) >= ?", lastDays).
		Count(&count)
	return count
}

func (r *dashboardRepository) countRecords(model interface{}) int64 {
	var count int64
	r.db.Model(model).Count(&count)
	return count
}
