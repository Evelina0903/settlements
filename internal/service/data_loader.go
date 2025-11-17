package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"TP_Andreev/internal/models"
	"gorm.io/gorm"
)

type DataLoaderService struct {
	db *gorm.DB
}

func NewDataLoaderService(db *gorm.DB) *DataLoaderService {
	return &DataLoaderService{db: db}
}

func (ds *DataLoaderService) LoadEmployeeTravelData(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Skip header row (index 0)
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 7 {
			fmt.Printf("Skipping row %d: insufficient columns\n", i)
			continue
		}

		err := ds.processRow(record)
		if err != nil {
			fmt.Printf("Error processing row %d: %v\n", i, err)
			continue
		}
	}

	return nil
}

func (ds *DataLoaderService) processRow(record []string) error {
	// CSV columns: Department, Employee, Travel Start Date, Travel End Date, Destination(s), Purpose Of Travel, Actual Total Expenses
	employeeName := strings.TrimSpace(record[1])
	destination := strings.TrimSpace(record[4])
	startDateStr := strings.TrimSpace(record[2])
	endDateStr := strings.TrimSpace(record[3])
	moneySpentStr := strings.TrimSpace(record[6])

	if employeeName == "" || destination == "" {
		return nil // Skip incomplete records
	}

	// Parse dates
	startDate, err := time.Parse("2006/01/02", startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start date format: %s", startDateStr)
	}

	endDate, err := time.Parse("2006/01/02", endDateStr)
	if err != nil {
		return fmt.Errorf("invalid end date format: %s", endDateStr)
	}

	// Parse money spent
	moneySpent := 0
	if moneySpentStr != "" && moneySpentStr != "0.00" {
		// Convert to cents (int) by multiplying by 100
		floatVal := 0.0
		_, err := fmt.Sscanf(moneySpentStr, "%f", &floatVal)
		if err == nil {
			moneySpent = int(floatVal * 100)
		}
	}

	// Find or create employee
	var employee models.Employee
	result := ds.db.Where("name = ?", employeeName).FirstOrCreate(&employee, models.Employee{Name: employeeName})
	if result.Error != nil {
		return fmt.Errorf("failed to create/find employee: %w", result.Error)
	}

	// Find or create business trip
	var trip models.BusinessTrip
	result = ds.db.Where(
		"destination = ? AND start_at = ? AND end_at = ?",
		destination, startDate, endDate,
	).FirstOrCreate(&trip, models.BusinessTrip{
		Destination: destination,
		StartAt:     startDate,
		EndAt:       endDate,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to create/find business trip: %w", result.Error)
	}

	// Create assignment
	assignment := models.AssignmentToTrip{
		EmployeeID:     employee.ID,
		BusinessTripID: trip.ID,
		MoneySpent:     moneySpent,
	}

	result = ds.db.Create(&assignment)
	if result.Error != nil {
		return fmt.Errorf("failed to create assignment: %w", result.Error)
	}

	return nil
}
