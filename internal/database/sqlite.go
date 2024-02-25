package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"polaris_food/internal/database/models"
)

type SQLiteConnector struct{}

func (connector *SQLiteConnector) InitDB(connectionString string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func (connector *SQLiteConnector) Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.Coupon{},
		&models.FoodCategory{},
		&models.MenuItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Rating{},
		&models.Restaurant{},
		&models.RiderProfile{},
		&models.User{},
		&models.UserCoupon{},
	)
}
