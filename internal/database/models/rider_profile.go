package models

import "time"

type RiderProfile struct {
	ID                 uint `gorm:"primaryKey"`
	AvailabilityStatus bool
	IsDelivering       bool
	UserID             uint `gorm:"unique"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Latitude           float64
	Longitude          float64
}
