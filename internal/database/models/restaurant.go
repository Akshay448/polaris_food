package models

import (
	"gorm.io/gorm"
)

type Restaurant struct {
	gorm.Model
	Name         string
	Address      string
	DeliveryTime int64 // average delivery time in minutes
	IsOpen       bool
	Latitude     float64
	Longitude    float64
}
