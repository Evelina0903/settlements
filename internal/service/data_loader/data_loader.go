package data_loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"settlements/internal/models"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type DataLoader struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DataLoader {
	return &DataLoader{db: db}
}

func (dl *DataLoader) LoadCityData(filePath string) error {
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

	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 14 {
			fmt.Printf("Skipping row %d: insufficient columns\n", i)
			continue
		}

		err := dl.processRow(record)
		if err != nil {
			fmt.Printf("Error processing row %d: %v\n", i, err)
			continue
		}
	}

	return nil
}

func (dl *DataLoader) processRow(record []string) error {
	// CSV columns: Region, Settlement, Type, Population, Children, Latitude, Longitude,
	region := strings.TrimSpace(record[1])
	settlement := strings.TrimSpace(record[3])
	typ := strings.TrimSpace(record[4])
	populationStr := strings.TrimSpace(record[5])
	childrenStr := strings.TrimSpace(record[6])
	latitudeStr := strings.TrimSpace(record[10])
	longitudeStr := strings.TrimSpace(record[11])

	population, _ := strconv.Atoi(populationStr)
	childrens, _ := strconv.Atoi(childrenStr)
	latitudeF64, _ := strconv.ParseFloat(latitudeStr, 32)
	longitudeF64, _ := strconv.ParseFloat(longitudeStr, 32)
	latitude := float32(latitudeF64)
	longitude := float32(longitudeF64)

	if population == 0 {
		return nil
	}

	var typeM models.Type
	err := dl.db.Where("name=?", typ).FirstOrCreate(&typeM, models.Type{Name: typ}).Error
	if err != nil {
		return fmt.Errorf("failed to create/find type: %w", err)
	}

	var district models.District
	err = dl.db.Where("name=?", region).FirstOrCreate(&district, models.District{Name: region}).Error
	if err != nil {
		return fmt.Errorf("failed to create/find district: %w", err)
	}

	err = dl.db.Create(
		&models.City{
			Name:       settlement,
			TypeID:     typeM.ID,
			DistrictID: district.ID,
			Population: population,
			Childrens:  childrens,
			Latitude: latitude,
			Longitude: longitude,
		},
	). Error
	if err != nil {
		return fmt.Errorf("failed to create city: %w", err)
	}

	return nil
}
