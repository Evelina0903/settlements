package repo

import (
	"log"
	"settlements/internal/dto"
	"settlements/internal/models"

	"gorm.io/gorm"
)

type CityRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CityRepo {
	return &CityRepo{db: db}
}

func (r *CityRepo) All() *[]dto.CityDTO {
	var cities []models.City
	err := r.db.Model(&models.City{}).Preload("Type").Preload("District").Find(&cities).Error
	if err != nil {
		log.Fatal(err)
	}

	res := []dto.CityDTO{}
	for _, c := range cities {
		cityDTO := dto.CityDTO{
			ID:         c.ID,
			Name:       c.Name,
			Type:       c.Type.Name,
			District:   c.District.Name,
			Population: c.Population,
			Childrens:  c.Childrens,
			Latitude:   c.Latitude,
			Longitude:  c.Longitude,
		}
		res = append(res, cityDTO)
	}

	return &res
}

func (r *CityRepo) MinLongitude() float64 {
	var city models.City
	err := r.db.Order("longitude").Limit(1).Find(&city).Error
	if err != nil {
		log.Fatal(err)
	}

	return city.Longitude
}

func (r *CityRepo) MaxLongitude() float64 {
	var city models.City
	err := r.db.Order("longitude desc").Limit(1).Find(&city).Error
	if err != nil {
		log.Fatal(err)
	}

	return city.Longitude
}

func (r *CityRepo) GetCitiesInLongitudeGap(lMin, lMax float64) *[]dto.CityDTO {
	var cities []models.City
	err := r.db.Where("longitude >= ? AND longitude < ?", lMin, lMax).Preload("Type").Preload("District").Find(&cities).Error
	if err != nil {
		log.Fatal(err)
	}

	res := []dto.CityDTO{}
	for _, c := range cities {
		cityDTO := dto.CityDTO{
			ID:         c.ID,
			Name:       c.Name,
			Type:       c.Type.Name,
			District:   c.District.Name,
			Population: c.Population,
			Childrens:  c.Childrens,
			Latitude:   c.Latitude,
			Longitude:  c.Longitude,
		}
		res = append(res, cityDTO)
	}

	return &res
}
