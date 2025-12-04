package dto

import (
	"testing"
)

func TestCityDTOCreation(t *testing.T) {
	city := CityDTO{
		ID:         1,
		Name:       "Moscow",
		Type:       "город",
		District:   "Moscow Region",
		Population: 12500000,
		Childrens:  2500000,
		Latitude:   55.7558,
		Longitude:  37.6173,
	}

	if city.ID != 1 {
		t.Errorf("Expected ID 1, got %d", city.ID)
	}

	if city.Name != "Moscow" {
		t.Errorf("Expected name 'Moscow', got %s", city.Name)
	}

	if city.Type != "город" {
		t.Errorf("Expected type 'город', got %s", city.Type)
	}

	if city.District != "Moscow Region" {
		t.Errorf("Expected district 'Moscow Region', got %s", city.District)
	}

	if city.Population != 12500000 {
		t.Errorf("Expected population 12500000, got %d", city.Population)
	}

	if city.Childrens != 2500000 {
		t.Errorf("Expected childrens 2500000, got %d", city.Childrens)
	}

	if city.Latitude != 55.7558 {
		t.Errorf("Expected latitude 55.7558, got %f", city.Latitude)
	}

	if city.Longitude != 37.6173 {
		t.Errorf("Expected longitude 37.6173, got %f", city.Longitude)
	}
}

func TestCityDTOZeroValues(t *testing.T) {
	city := CityDTO{}

	if city.ID != 0 {
		t.Errorf("Expected zero ID, got %d", city.ID)
	}

	if city.Name != "" {
		t.Errorf("Expected empty name, got %s", city.Name)
	}

	if city.Population != 0 {
		t.Errorf("Expected zero population, got %d", city.Population)
	}
}
