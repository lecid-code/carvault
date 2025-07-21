package models

import (
	"time"
)

// Database model
type Expense struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	VehicleID   int       `db:"vehicle_id"`
	Date        time.Time `db:"date"`
	Mileage     *int      `db:"mileage"`
	Amount      float64   `db:"amount"`
	ExpenseType string    `db:"expense_type"`
	Details     *string   `db:"details"`
	CreatedAt   time.Time `db:"created_at"`
}

// API response model
type ExpenseResponse struct {
	ID          int       `json:"id"`
	VehicleID   int       `json:"vehicle_id"`
	Date        time.Time `json:"date"`
	Mileage     *int      `json:"mileage"`
	Amount      float64   `json:"amount"`
	ExpenseType string    `json:"expense_type"`
	Details     *string   `json:"details"`
	CreatedAt   time.Time `json:"created_at"`
}