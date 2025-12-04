package service

import (
	"settlements/internal/dto"
	"settlements/internal/repo"
)

// ServiceV2 is the REFACTORED version using Strategy Pattern
// The original Service is kept for backward compatibility

// ServiceV2 encapsulates business logic using flexible aggregation strategies
// This enables easy extension without modifying existing code (Open/Closed Principle)
type ServiceV2 struct {
	aggregator *StrategyAggregator
}

// NewServiceV2 creates a new ServiceV2 with strategy aggregator
func NewServiceV2(cityRepo *repo.CityRepo) *ServiceV2 {
	return &ServiceV2{
		aggregator: NewStrategyAggregator(cityRepo),
	}
}

// GetSettlementTypeData returns aggregated settlement type statistics
// Uses SettlementTypeAggregationStrategy internally
func (s *ServiceV2) GetSettlementTypeData() *[]SettlementTypeData {
	strategy := &SettlementTypeAggregationStrategy{}
	result := s.aggregator.Aggregate(strategy)
	return result.(*[]SettlementTypeData)
}

// GetDistrictPopulationData returns aggregated district population data
// Uses DistrictAggregationStrategy internally
func (s *ServiceV2) GetDistrictPopulationData() *[]GraphData {
	strategy := &DistrictAggregationStrategy{}
	result := s.aggregator.Aggregate(strategy)
	return result.(*[]GraphData)
}

// GetLongitudePopulationData returns aggregated longitude-based population data
// Uses LongitudeAggregationStrategy internally with default 100 buckets
func (s *ServiceV2) GetLongitudePopulationData() *[]GraphData {
	strategy := NewLongitudeAggregationStrategy(100)
	result := s.aggregator.Aggregate(strategy)
	return result.(*[]GraphData)
}

// GetLongitudePopulationDataWithBuckets returns aggregated longitude data with custom bucket count
// Allows customization of aggregation granularity
func (s *ServiceV2) GetLongitudePopulationDataWithBuckets(bucketCount int) *[]GraphData {
	strategy := NewLongitudeAggregationStrategy(bucketCount)
	result := s.aggregator.Aggregate(strategy)
	return result.(*[]GraphData)
}

// ExecuteCustomStrategy allows execution of custom aggregation strategies
// Demonstrates the extensibility of the Strategy pattern
func (s *ServiceV2) ExecuteCustomStrategy(strategy AggregationStrategy) interface{} {
	return s.aggregator.Aggregate(strategy)
}

// ExecuteMultipleStrategies executes multiple strategies in sequence
// Efficient for fetching multiple aggregations in one pass
func (s *ServiceV2) ExecuteMultipleStrategies(strategies ...AggregationStrategy) []interface{} {
	return s.aggregator.AggregateMultiple(strategies...)
}

// Example of extending functionality: New strategy can be added without modifying Service
// This is a demonstration of how the Strategy Pattern achieves Open/Closed Principle

// CustomAggregationStrategy example showing extensibility
type CustomAggregationStrategy struct {
	filterFunc func(*dto.CityDTO) bool // Custom filter logic
}

// Aggregate implements AggregationStrategy interface
func (s *CustomAggregationStrategy) Aggregate(cities *[]dto.CityDTO) interface{} {
	filtered := []dto.CityDTO{}
	for _, city := range *cities {
		if s.filterFunc(&city) {
			filtered = append(filtered, city)
		}
	}
	return &filtered
}

// Name returns the strategy name
func (s *CustomAggregationStrategy) Name() string {
	return "custom_aggregation"
}
