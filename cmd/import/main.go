package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lecid-code/carvault/internal/database"
)

const (
	KiaUUID   = "8432be0a-80d4-4e53-85c3-9c91a2934cca"
	RoverUUID = "a28da350-c686-42e5-b191-f0b6ce315852"
	UserID    = 1
	KiaID     = 1
	RoverID   = 2
)

func main() {
	// Connect to database
	db, err := database.New("data/carvault.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create vehicle mapping
	vehicleMap := map[string]int{
		KiaUUID:   KiaID,
		RoverUUID: RoverID,
	}

	// Import CSV data
	if err := importCSV(db, vehicleMap, "data/expenses.csv"); err != nil {
		log.Fatalf("Failed to import CSV: %v", err)
	}

	fmt.Println("Import completed successfully!")
}

func importCSV(db *database.DB, vehicleMap map[string]int, csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	count := 0
	for i, record := range records[1:] { // Skip header
		if len(record) < 7 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		// Parse fields
		vehicleUUID := record[1]
		category := record[2]
		dateStr := record[3]
		amountCents, err := strconv.Atoi(record[4])
		if err != nil {
			log.Printf("Skipping row %d: invalid amount: %v", i+2, err)
			continue
		}
		
		var mileage *int
		if record[5] != "" && record[5] != "0" {
			m, err := strconv.Atoi(record[5])
			if err != nil {
				log.Printf("Warning: invalid mileage in row %d: %v", i+2, err)
			} else {
				mileage = &m
			}
		}
		
		details := record[6]

		// Get vehicle ID
		vehicleID, exists := vehicleMap[vehicleUUID]
		if !exists {
			log.Printf("Skipping row %d: unknown vehicle UUID: %s", i+2, vehicleUUID)
			continue
		}

		// Parse date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("Skipping row %d: invalid date format: %v", i+2, err)
			continue
		}

		// Convert cents to dollars
		amount := float64(amountCents) / 100.0

		// Validate JSON details
		var detailsPtr *string
		if details != "" && details != "{}" {
			var temp interface{}
			if err := json.Unmarshal([]byte(details), &temp); err != nil {
				log.Printf("Warning: invalid JSON in row %d: %v", i+2, err)
			} else {
				detailsPtr = &details
			}
		}

		// Insert expense
		_, err = db.Conn().NamedExec(`
			INSERT INTO expenses (user_id, vehicle_id, date, mileage, amount, expense_type, details, created_at)
			VALUES (:user_id, :vehicle_id, :date, :mileage, :amount, :expense_type, :details, :created_at)`,
			map[string]interface{}{
				"user_id":      UserID,
				"vehicle_id":   vehicleID,
				"date":         date,
				"mileage":      mileage,
				"amount":       amount,
				"expense_type": category,
				"details":      detailsPtr,
				"created_at":   time.Now(),
			})
		if err != nil {
			return fmt.Errorf("failed to insert expense from row %d: %w", i+2, err)
		}

		count++
		if count%50 == 0 {
			fmt.Printf("Imported %d expenses...\n", count)
		}
	}

	fmt.Printf("Successfully imported %d expenses\n", count)
	return nil
}