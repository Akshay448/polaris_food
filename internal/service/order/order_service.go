package order

import (
	"errors"
	"gorm.io/gorm"
	"polaris_food/internal/database/models"
)

type OrderService interface {
	CreateOrder(order models.Order, orderItems []models.OrderItem) error
	UpdateOrder(orderID uint, status string) error
}

type orderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderService {
	return &orderService{db}
}

// CreateOrder creates a new order and puts all the order items in order_items table
func (s *orderService) CreateOrder(order models.Order, orderItems []models.OrderItem) error {
	if err := s.db.Create(&order).Error; err != nil {
		return err
	}
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.db.Create(&orderItems[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateOrder updates the status of a given order
func (s *orderService) UpdateOrder(orderID uint, status string) error {
	result := s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no order found with the given ID")
	}
	return nil
}
