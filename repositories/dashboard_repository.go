package repositories

import (
	"github.com/gofiber/fiber/v2/log"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetInfo(lastNDays int) models.Dashboard
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(database *gorm.DB) DashboardRepository {
	return &dashboardRepository{
		db: database,
	}
}

func (r *dashboardRepository) GetInfo(lastNDays int) models.Dashboard {
	const dateFormat = "2006-01-02"
	lastDays := time.Now().AddDate(0, 0, -lastNDays).Format(dateFormat)

	dashboard := models.Dashboard{
		NewOrders: r.countOrderStatusInLastDays(lastDays, []models.ShoppingCartStatus{
			models.ActiveStatus,
			models.AwaitingPaymentStatus,
			models.DeliveredStatusAwaiting,
			models.DeliveredStatusDelivered,
			models.PaymentApprovedStatus,
			models.ProcessingStatus,
		}),
		Sales: r.countOrderStatusInLastDays(lastDays, []models.ShoppingCartStatus{
			models.DeliveredStatusAwaiting,
			models.DeliveredStatusDelivered,
			models.PaymentApprovedStatus,
		}),
		Users:    r.countRecords(&models.User{}),
		Products: r.countRecords(&models.Product{}),
	}

	return dashboard
}

func (r *dashboardRepository) countOrderStatusInLastDays(lastDays string, statuses []models.ShoppingCartStatus) int64 {
	var count int64

	if err := r.db.Model(&models.Order{}).Where("status IN ?", statuses).Where("DATE(created_at) >= ?", lastDays).Count(&count).Error; err != nil {
		log.Errorf("DashboardRepository countOrderStatusInLastDays: %s", err.Error())
		return 0
	}

	return count
}

func (r *dashboardRepository) countRecords(model interface{}) int64 {
	var count int64

	if err := r.db.Model(model).Count(&count).Error; err != nil {
		log.Errorf("DashboardRepository countRecords: %s", err.Error())
		return 0
	}

	return count
}
