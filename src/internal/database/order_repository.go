package database

import (
	"gorm.io/gorm"
	"myapp/src/internal/order_model"
)

type OrderRepository struct{
	db *gorm.DB
}


func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	err := r.db.Transaction(func(tx *gorm.DB) error{
		return tx.Create(order).Error
	})
	return err
}

func (r *OrderRepository) GetByTrackNumber(trackNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.Transaction(func(tx *gorm.DB) error{
		if err := tx.Preload("Items").First(&order, "track_number = ?", trackNumber).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindForCache(size int) ([]model.Order , error) {
	var latestOrders []model.Order
	err := r.db.Transaction(func(tx *gorm.DB) error{
		if err := tx.Preload("Items").Order("created_at desc").Limit(size).Find(&latestOrders).Error ; err != nil {
			return err
		}
		return nil
	})
	if err != nil {	
		return nil, err
	}
	return latestOrders, nil
}