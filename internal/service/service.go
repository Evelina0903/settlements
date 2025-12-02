package service

import (
	"settlements/internal/repo"
	"sort"
)

type Service struct {
	cityRepo *repo.CityRepo
}

type SettlementTypeData struct {
	Type          string  `json:"type"`
	AvgPopulation float32 `json:"avgPopulation"`
	AvgChildrens  float32 `json:"avgChildrens"`
	MinPopulation int     `json:"minPopulation"`
	MaxPopulation int     `json:"maxPopulation"`
}

type GraphData struct {
	X any `json:"x"`
	Y int `json:"y"`
}

func New(cityRepo *repo.CityRepo) *Service {
	return &Service{cityRepo: cityRepo}
}

func (s *Service) GetAllSettelmetTypeData() *[]SettlementTypeData {
	data := s.cityRepo.All()

	populationAcc := map[string]int{}
	childrenAcc := map[string]int{}
	minPopulation := map[string]int{}
	maxPopulation := map[string]int{}
	citiesCounter := map[string]int{}

	for _, d := range *data {
		_, ok := citiesCounter[d.Type]
		if ok {
			citiesCounter[d.Type]++
		} else {
			citiesCounter[d.Type] = 1
		}

		_, ok = populationAcc[d.Type]
		if ok {
			populationAcc[d.Type] += d.Population
		} else {
			populationAcc[d.Type] = d.Population
		}

		_, ok = childrenAcc[d.Type]
		if ok {
			childrenAcc[d.Type] += d.Childrens
		} else {
			childrenAcc[d.Type] = d.Childrens
		}

		_, ok = minPopulation[d.Type]
		if ok {
			if d.Population < minPopulation[d.Type] {
				minPopulation[d.Type] = d.Population
			}
		} else {
			minPopulation[d.Type] = d.Population
		}

		_, ok = maxPopulation[d.Type]
		if ok {
			if d.Population > maxPopulation[d.Type] {
				maxPopulation[d.Type] = d.Population
			}
		} else {
			maxPopulation[d.Type] = d.Population
		}
	}

	res := []SettlementTypeData{}

	for k, v := range citiesCounter {
		typeData := SettlementTypeData{
			Type:          k,
			AvgPopulation: float32(populationAcc[k]) / float32(v),
			AvgChildrens:  float32(childrenAcc[k]) / float32(v),
			MinPopulation: minPopulation[k],
			MaxPopulation: maxPopulation[k],
		}
		res = append(res, typeData)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].AvgPopulation > res[j].AvgPopulation
	})

	return &res
}

func (s *Service) GetLongitudePopulationData() *[]GraphData {
	min := s.cityRepo.MinLongitude()
	max := s.cityRepo.MaxLongitude()
	step := (max - min) / 100

	res := []GraphData{}

	for i := min; i <= max-step; i += step {
		data := s.cityRepo.GetCitiesInLongitudeGap(i, i+step)

		sum := 0
		for _, d := range *data {
			sum += d.Population
		}

		res = append(res, GraphData{
			X: i,
			Y: sum,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].X.(float64) < res[j].X.(float64)
	})

	return &res
}

func (s *Service) GetDistrictPopulationData() *[]GraphData {
	data := s.cityRepo.All()

	populationAcc := map[string]int{}

	for _, d := range *data {
		_, ok := populationAcc[d.District]
		if ok {
			populationAcc[d.District] += d.Population
		} else {
			populationAcc[d.District] = d.Population
		}
	}

	res := []GraphData{}

	for k, v := range populationAcc {
		res = append(res, GraphData{
			X: k,
			Y: v,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Y > res[j].Y
	})

	return &res
}
