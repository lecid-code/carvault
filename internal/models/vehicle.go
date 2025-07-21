package models

import (
	"time"
)

// Database model
type Vehicle struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	Name            string    `db:"name"`
	Make            *string   `db:"make"`
	Model           *string   `db:"model"`
	Year            *int      `db:"year"`
	VIN             *string   `db:"vin"`
	LicensePlate    *string   `db:"license_plate"`
	DatePurchased   *time.Time `db:"date_purchased"`
	PurchaseMileage *int      `db:"purchase_mileage"`
	CreatedAt       time.Time `db:"created_at"`
}

// API response model
type VehicleResponse struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Make            *string    `json:"make"`
	Model           *string    `json:"model"`
	Year            *int       `json:"year"`
	VIN             *string    `json:"vin"`
	LicensePlate    *string    `json:"license_plate"`
	DatePurchased   *time.Time `json:"date_purchased"`
	PurchaseMileage *int       `json:"purchase_mileage"`
	CreatedAt       time.Time  `json:"created_at"`
}