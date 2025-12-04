package service

import (
	"sort"

	"settlements/internal/dto"
	"settlements/internal/repo"
)

// AggregationStrategy defines the interface for different data aggregation strategies
// This implements the Strategy Pattern to allow flexible data aggregation approaches
type AggregationStrategy interface {
	// Aggregate processes city data according to the strategy and returns results
	// The input is a slice of cities, output format depends on the strategy
	Aggregate(cities *[]dto.CityDTO) interface{}

	// Name returns a descriptive name of the strategy
	Name() string
}

// SettlementTypeAggregationStrategy aggregates cities by settlement type
// Computes statistics: average population, average children, min/max population
type SettlementTypeAggregationStrategy struct{}

// Aggregate groups cities by type and calculates statistics
func (s *SettlementTypeAggregationStrategy) Aggregate(cities *[]dto.CityDTO) interface{} {
	populationAcc := make(map[string]int)
	childrenAcc := make(map[string]int)
	minPopulation := make(map[string]int)
	maxPopulation := make(map[string]int)
	citiesCounter := make(map[string]int)

	for _, d := range *cities {
		// Track count
		citiesCounter[d.Type]++

		// Accumulate population
		populationAcc[d.Type] += d.Population

		// Accumulate children
		childrenAcc[d.Type] += d.Childrens

		// Track min population
		if _, exists := minPopulation[d.Type]; exists {
			if d.Population < minPopulation[d.Type] {
				minPopulation[d.Type] = d.Population
			}
		} else {
			minPopulation[d.Type] = d.Population
		}

		// Track max population
		if _, exists := maxPopulation[d.Type]; exists {
			if d.Population > maxPopulation[d.Type] {
				maxPopulation[d.Type] = d.Population
			}
		} else {
			maxPopulation[d.Type] = d.Population
		}
	}

	// Build result
	result := []SettlementTypeData{}
	for typeKey, count := range citiesCounter {
		typeData := SettlementTypeData{
			Type:          typeKey,
			AvgPopulation: float32(populationAcc[typeKey]) / float32(count),
			AvgChildrens:  float32(childrenAcc[typeKey]) / float32(count),
			MinPopulation: minPopulation[typeKey],
			MaxPopulation: maxPopulation[typeKey],
		}
		result = append(result, typeData)
	}

	// Sort by average population descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].AvgPopulation > result[j].AvgPopulation
	})

	return &result
}

// Name returns the strategy name
func (s *SettlementTypeAggregationStrategy) Name() string {
	return "settlement_type_aggregation"
}

// DistrictAggregationStrategy aggregates cities by district
// Computes total population per district
type DistrictAggregationStrategy struct{}

// Aggregate groups cities by district and calculates total population
func (s *DistrictAggregationStrategy) Aggregate(cities *[]dto.CityDTO) interface{} {
	populationAcc := make(map[string]int)

	for _, d := range *cities {
		populationAcc[d.District] += d.Population
	}

	result := []GraphData{}
	for district, population := range populationAcc {
		result = append(result, GraphData{
			X: district,
			Y: population,
		})
	}

	// Sort by population descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Y > result[j].Y
	})

	return &result
}

// Name returns the strategy name
func (s *DistrictAggregationStrategy) Name() string {
	return "district_aggregation"
}

// LongitudeAggregationStrategy distributes cities into longitude buckets
// Calculates total population per longitude range
type LongitudeAggregationStrategy struct {
	bucketCount int
}

// NewLongitudeAggregationStrategy creates a new longitude strategy with specified bucket count
func NewLongitudeAggregationStrategy(bucketCount int) *LongitudeAggregationStrategy {
	if bucketCount <= 0 {
		bucketCount = 100 // default
	}
	return &LongitudeAggregationStrategy{bucketCount: bucketCount}
}

// Aggregate distributes cities into longitude buckets and sums population
func (s *LongitudeAggregationStrategy) Aggregate(cities *[]dto.CityDTO) interface{} {
	if len(*cities) == 0 {
		return &[]GraphData{}
	}

	// Find min/max longitude
	minLong := (*cities)[0].Longitude
	maxLong := (*cities)[0].Longitude

	for _, c := range *cities {
		if c.Longitude < minLong {
			minLong = c.Longitude
		}
		if c.Longitude > maxLong {
			maxLong = c.Longitude
		}
	}

	step := (maxLong - minLong) / float64(s.bucketCount)

	// Build buckets
	result := []GraphData{}
	for i := 0; i < s.bucketCount; i++ {
		bucketStart := minLong + float64(i)*step
		bucketEnd := bucketStart + step

		sum := 0
		for _, d := range *cities {
			if d.Longitude >= bucketStart && d.Longitude < bucketEnd {
				sum += d.Population
			}
		}

		result = append(result, GraphData{
			X: bucketStart,
			Y: sum,
		})
	}

	// Sort by longitude ascending
	sort.Slice(result, func(i, j int) bool {
		return result[i].X.(float64) < result[j].X.(float64)
	})

	return &result
}

// Name returns the strategy name
func (s *LongitudeAggregationStrategy) Name() string {
	return "longitude_aggregation"
}

// StrategyAggregator is a context class that uses aggregation strategies
// Allows switching between different aggregation approaches at runtime
type StrategyAggregator struct {
	repo *repo.CityRepo
}

// NewStrategyAggregator creates a new aggregator with a repository
func NewStrategyAggregator(repository *repo.CityRepo) *StrategyAggregator {
	return &StrategyAggregator{
		repo: repository,
	}
}

// Aggregate executes the provided strategy with city data from the repository
func (sa *StrategyAggregator) Aggregate(strategy AggregationStrategy) interface{} {
	cities := sa.repo.All()
	return strategy.Aggregate(cities)
}

// AggregateMultiple executes multiple strategies and returns results in order
func (sa *StrategyAggregator) AggregateMultiple(strategies ...AggregationStrategy) []interface{} {
	results := make([]interface{}, len(strategies))
	for i, strategy := range strategies {
		results[i] = sa.Aggregate(strategy)
	}
	return results
}
