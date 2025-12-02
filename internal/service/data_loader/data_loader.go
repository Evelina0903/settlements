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

var settlementsTypes = map[string]string{
	"г":             "город",
	"гп":            "городской поселок",
	"д":             "деревня",
	"дп":            "дачный поселок",
	"ж/д блокпост":  "железнодорожный блокпост",
	"ж/д будка":     "железнодорожная будка",
	"ж/д ветка":     "железнодорожная ветка",
	"ж/д казарма":   "железнодорожная казарма",
	"ж/д комбинат":  "железнодорожный комбинат",
	"ж/д оп":        "железнодорожный остановочный пункт",
	"ж/д платформа": "железнодорожная платформа",
	"ж/д площадка":  "железнодорожная площадка",
	"ж/д пост":      "железнодорожный путевой пост",
	"ж/д рзд":       "железнодорожный разъезд",
	"ж/д ст":        "железнодорожная станция",
	"зим":           "зимовье",
	"к":             "кишлак",
	"кп":            "курортный поселок",
	"л/п":           "лесной поселок",
	"м":             "местечко",
	"мкр":           "микрорайон",
	"нп":            "населенный пункт",
	"п ж/д ст":      "поселок при железнодорожной станции",
	"п":             "поселок",
	"п/ст":          "поселок при станции",
	"пгт":           "поселок городского типа",
	"р-н":           "район",
	"рзд":           "разъезд",
	"рп":            "рабочий поселок",
	"с":             "село",
	"с/п":           "сельское поселение",
	"сл":            "слобода",
	"снт":           "садоводческое некоммерческое товарищество",
	"ст-ца":         "станица",
	"ст":            "станция",
	"тер":           "территория",
	"у":             "улус",
	"х":             "хутор",
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
	// CSV columns: Region, Settlement, Type, Population, Children, Latitude, Longitude
	region := strings.TrimSpace(record[1])
	settlement := strings.TrimSpace(record[3])
	typeShort := strings.TrimSpace(record[4])
	populationStr := strings.TrimSpace(record[5])
	childrenStr := strings.TrimSpace(record[6])
	latitudeStr := strings.TrimSpace(record[9])
	longitudeStr := strings.TrimSpace(record[10])

	population, _ := strconv.Atoi(populationStr)
	childrens, _ := strconv.Atoi(childrenStr)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	if population == 0 {
		return nil
	}

	if longitude < 0 {
		longitude = 180 - longitude
	}

	typ := ""
	v, ok := settlementsTypes[typeShort]
	if ok {
		typ = v
	} else {
		typ = typeShort
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

	if region == settlement {
		var city models.City
		err = dl.db.Where("name=? AND latitude=? AND longitude=?", settlement, latitude, longitude).FirstOrCreate(
			&city,
			models.City{
				Name:       settlement,
				TypeID:     typeM.ID,
				DistrictID: district.ID,
				Population: 0,
				Childrens:  0,
				Latitude:   latitude,
				Longitude:  longitude,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to create/find city: %w", err)
		}
		city.Population += population
		city.Childrens += childrens

		err = dl.db.Save(&city).Error
		if err != nil {
			return fmt.Errorf("failed to save updated city: %w", err)
		}

		return nil
	}

	city := models.City{
		Name:       settlement,
		TypeID:     typeM.ID,
		DistrictID: district.ID,
		Population: population,
		Childrens:  childrens,
		Latitude:   latitude,
		Longitude:  longitude,
	}

	err = dl.db.Create(&city).Error
	if err != nil {
		return fmt.Errorf("failed to create city: %w", err)
	}

	return nil
}
