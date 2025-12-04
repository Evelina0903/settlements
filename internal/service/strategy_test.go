package service

import (
	"testing"

	"settlements/internal/dto"
)

func TestSettlementTypeAggregationStrategy(t *testing.T) {
	cities := []dto.CityDTO{
		{Type: "город", Population: 1000000, Childrens: 200000},
		{Type: "город", Population: 5000000, Childrens: 1000000},
		{Type: "деревня", Population: 500, Childrens: 100},
		{Type: "деревня", Population: 1000, Childrens: 200},
	}

	strategy := &SettlementTypeAggregationStrategy{}
	result := strategy.Aggregate(&cities)
	data := result.(*[]SettlementTypeData)

	if len(*data) != 2 {
		t.Errorf("Expected 2 settlement types, got %d", len(*data))
	}

	// First should be город (higher avg population)
	if (*data)[0].Type != "город" {
		t.Errorf("Expected first type 'город', got %s", (*data)[0].Type)
	}

	expectedAvg := float32(6000000) / 2
	if (*data)[0].AvgPopulation != expectedAvg {
		t.Errorf("Expected avg population %f, got %f", expectedAvg, (*data)[0].AvgPopulation)
	}
}

func TestSettlementTypeAggregationStrategyName(t *testing.T) {
	strategy := &SettlementTypeAggregationStrategy{}
	if strategy.Name() != "settlement_type_aggregation" {
		t.Errorf("Expected name 'settlement_type_aggregation', got %s", strategy.Name())
	}
}

func TestDistrictAggregationStrategy(t *testing.T) {
	cities := []dto.CityDTO{
		{District: "Moscow", Population: 5000000},
		{District: "Moscow", Population: 2000000},
		{District: "SPB", Population: 3000000},
	}

	strategy := &DistrictAggregationStrategy{}
	result := strategy.Aggregate(&cities)
	data := result.(*[]GraphData)

	if len(*data) != 2 {
		t.Errorf("Expected 2 districts, got %d", len(*data))
	}

	// First should be Moscow (highest population)
	if (*data)[0].X != "Moscow" {
		t.Errorf("Expected first district 'Moscow', got %s", (*data)[0].X)
	}

	if (*data)[0].Y != 7000000 {
		t.Errorf("Expected Moscow population 7000000, got %d", (*data)[0].Y)
	}
}

func TestDistrictAggregationStrategyName(t *testing.T) {
	strategy := &DistrictAggregationStrategy{}
	if strategy.Name() != "district_aggregation" {
		t.Errorf("Expected name 'district_aggregation', got %s", strategy.Name())
	}
}

func TestLongitudeAggregationStrategy(t *testing.T) {
	cities := []dto.CityDTO{
		{Longitude: 30.0, Population: 1000},
		{Longitude: 31.0, Population: 2000},
		{Longitude: 32.0, Population: 1500},
	}

	strategy := NewLongitudeAggregationStrategy(10)
	result := strategy.Aggregate(&cities)
	data := result.(*[]GraphData)

	if len(*data) != 10 {
		t.Errorf("Expected 10 buckets, got %d", len(*data))
	}

	// Verify sorted by longitude
	for i := 0; i < len(*data)-1; i++ {
		if (*data)[i].X.(float64) > (*data)[i+1].X.(float64) {
			t.Errorf("Data not sorted by longitude at index %d", i)
		}
	}
}

func TestLongitudeAggregationStrategyName(t *testing.T) {
	strategy := NewLongitudeAggregationStrategy(100)
	if strategy.Name() != "longitude_aggregation" {
		t.Errorf("Expected name 'longitude_aggregation', got %s", strategy.Name())
	}
}

func TestLongitudeAggregationStrategyDefaultBuckets(t *testing.T) {
	strategy := NewLongitudeAggregationStrategy(0)
	if strategy.bucketCount != 100 {
		t.Errorf("Expected default bucket count 100, got %d", strategy.bucketCount)
	}
}

func TestLongitudeAggregationStrategyEmpty(t *testing.T) {
	cities := []dto.CityDTO{}
	strategy := NewLongitudeAggregationStrategy(10)
	result := strategy.Aggregate(&cities)
	data := result.(*[]GraphData)

	if len(*data) != 0 {
		t.Errorf("Expected 0 data points for empty cities, got %d", len(*data))
	}
}

func TestCustomAggregationStrategy(t *testing.T) {
	cities := []dto.CityDTO{
		{Name: "City1", Population: 1000000},
		{Name: "City2", Population: 500000},
		{Name: "City3", Population: 100000},
	}

	// Filter for cities with population > 600000
	strategy := &CustomAggregationStrategy{
		filterFunc: func(c *dto.CityDTO) bool {
			return c.Population > 600000
		},
	}

	result := strategy.Aggregate(&cities)
	filtered := result.(*[]dto.CityDTO)

	if len(*filtered) != 1 {
		t.Errorf("Expected 1 filtered city, got %d", len(*filtered))
	}

	if (*filtered)[0].Population != 1000000 {
		t.Errorf("Expected city with population 1000000, got %d", (*filtered)[0].Population)
	}
}

func TestCustomAggregationStrategyName(t *testing.T) {
	strategy := &CustomAggregationStrategy{}
	if strategy.Name() != "custom_aggregation" {
		t.Errorf("Expected name 'custom_aggregation', got %s", strategy.Name())
	}
}

func TestAggregationStrategyInterface(t *testing.T) {
	var _ AggregationStrategy = &SettlementTypeAggregationStrategy{}
	var _ AggregationStrategy = &DistrictAggregationStrategy{}
	var _ AggregationStrategy = &LongitudeAggregationStrategy{}
	var _ AggregationStrategy = &CustomAggregationStrategy{}
	// All implement AggregationStrategy interface
}
