package models

import (
	"testing"
)

func TestCityModel(t *testing.T) {
	city := City{
		ID:         1,
		Name:       "Moscow",
		TypeID:     1,
		DistrictID: 1,
		Population: 12000000,
		Childrens:  2000000,
		Latitude:   55.7558,
		Longitude:  37.6173,
	}

	if city.ID != 1 {
		t.Errorf("Expected ID 1, got %d", city.ID)
	}

	if city.Name != "Moscow" {
		t.Errorf("Expected name 'Moscow', got %s", city.Name)
	}

	if city.Population != 12000000 {
		t.Errorf("Expected population 12000000, got %d", city.Population)
	}

	if city.Latitude != 55.7558 {
		t.Errorf("Expected latitude 55.7558, got %f", city.Latitude)
	}

	if city.Longitude != 37.6173 {
		t.Errorf("Expected longitude 37.6173, got %f", city.Longitude)
	}
}

func TestCityRelations(t *testing.T) {
	cityType := Type{
		ID:   1,
		Name: "город",
	}

	district := District{
		ID:   1,
		Name: "Moscow Region",
	}

	city := City{
		ID:         1,
		Name:       "Moscow",
		TypeID:     1,
		DistrictID: 1,
		Type:       cityType,
		District:   district,
		Population: 12000000,
		Childrens:  2000000,
		Latitude:   55.7558,
		Longitude:  37.6173,
	}

	if city.Type.Name != "город" {
		t.Errorf("Expected type name 'город', got %s", city.Type.Name)
	}

	if city.District.Name != "Moscow Region" {
		t.Errorf("Expected district name 'Moscow Region', got %s", city.District.Name)
	}
}

func TestTypeModel(t *testing.T) {
	typ := Type{
		ID:   1,
		Name: "город",
	}

	if typ.ID != 1 {
		t.Errorf("Expected ID 1, got %d", typ.ID)
	}

	if typ.Name != "город" {
		t.Errorf("Expected name 'город', got %s", typ.Name)
	}

	if len(typ.Citys) != 0 {
		t.Errorf("Expected empty Citys slice, got %d items", len(typ.Citys))
	}
}

func TestTypeOneToMany(t *testing.T) {
	city1 := City{ID: 1, Name: "Moscow"}
	city2 := City{ID: 2, Name: "SPB"}

	typ := Type{
		ID:    1,
		Name:  "город",
		Citys: []City{city1, city2},
	}

	if len(typ.Citys) != 2 {
		t.Errorf("Expected 2 cities, got %d", len(typ.Citys))
	}

	if typ.Citys[0].Name != "Moscow" {
		t.Errorf("Expected first city 'Moscow', got %s", typ.Citys[0].Name)
	}

	if typ.Citys[1].Name != "SPB" {
		t.Errorf("Expected second city 'SPB', got %s", typ.Citys[1].Name)
	}
}

func TestDistrictModel(t *testing.T) {
	district := District{
		ID:   1,
		Name: "Moscow Region",
	}

	if district.ID != 1 {
		t.Errorf("Expected ID 1, got %d", district.ID)
	}

	if district.Name != "Moscow Region" {
		t.Errorf("Expected name 'Moscow Region', got %s", district.Name)
	}

	if len(district.Citys) != 0 {
		t.Errorf("Expected empty Citys slice, got %d items", len(district.Citys))
	}
}

func TestDistrictOneToMany(t *testing.T) {
	city1 := City{ID: 1, Name: "Moscow"}
	city2 := City{ID: 2, Name: "Tver"}

	district := District{
		ID:    1,
		Name:  "Moscow Region",
		Citys: []City{city1, city2},
	}

	if len(district.Citys) != 2 {
		t.Errorf("Expected 2 cities, got %d", len(district.Citys))
	}

	if district.Citys[0].Name != "Moscow" {
		t.Errorf("Expected first city 'Moscow', got %s", district.Citys[0].Name)
	}

	if district.Citys[1].Name != "Tver" {
		t.Errorf("Expected second city 'Tver', got %s", district.Citys[1].Name)
	}
}

func TestCityCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{"Moscow", 55.7558, 37.6173},
		{"SPB", 59.9501, 30.3594},
		{"Yekaterinburg", 56.8389, 60.6057},
		{"Novosibirsk", 55.0415, 82.9346},
	}

	for _, test := range tests {
		city := City{
			ID:        1,
			Name:      test.name,
			Latitude:  test.latitude,
			Longitude: test.longitude,
		}

		if city.Latitude != test.latitude {
			t.Errorf("%s: expected latitude %f, got %f", test.name, test.latitude, city.Latitude)
		}

		if city.Longitude != test.longitude {
			t.Errorf("%s: expected longitude %f, got %f", test.name, test.longitude, city.Longitude)
		}
	}
}

func TestCityPopulationEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		population int
		childrens  int
	}{
		{"Large city", 12000000, 2000000},
		{"Small town", 50000, 10000},
		{"Village", 1000, 200},
		{"Empty settlement", 0, 0},
	}

	for _, test := range tests {
		city := City{
			ID:         1,
			Name:       test.name,
			Population: test.population,
			Childrens:  test.childrens,
		}

		if city.Population != test.population {
			t.Errorf("%s: expected population %d, got %d", test.name, test.population, city.Population)
		}

		if city.Childrens != test.childrens {
			t.Errorf("%s: expected childrens %d, got %d", test.name, test.childrens, city.Childrens)
		}
	}
}
