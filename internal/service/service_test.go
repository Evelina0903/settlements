package service

import (
	"testing"
)

func TestGraphDataStructure(t *testing.T) {
	data := GraphData{
		X: "test",
		Y: 100,
	}

	if data.X != "test" {
		t.Errorf("Expected X 'test', got %v", data.X)
	}

	if data.Y != 100 {
		t.Errorf("Expected Y 100, got %d", data.Y)
	}
}

func TestGraphDataWithFloat(t *testing.T) {
	data := GraphData{
		X: 37.6173,
		Y: 5000000,
	}

	xVal, ok := data.X.(float64)
	if !ok {
		t.Errorf("Expected X to be float64, got %T", data.X)
	}

	if xVal != 37.6173 {
		t.Errorf("Expected X 37.6173, got %f", xVal)
	}

	if data.Y != 5000000 {
		t.Errorf("Expected Y 5000000, got %d", data.Y)
	}
}

func TestSettlementTypeDataStructure(t *testing.T) {
	data := SettlementTypeData{
		Type:          "город",
		AvgPopulation: 1000000.5,
		AvgChildrens:  200000.2,
		MinPopulation: 100000,
		MaxPopulation: 5000000,
	}

	if data.Type != "город" {
		t.Errorf("Expected type 'город', got %s", data.Type)
	}

	if data.AvgPopulation != 1000000.5 {
		t.Errorf("Expected avg population 1000000.5, got %f", data.AvgPopulation)
	}

	if data.AvgChildrens != 200000.2 {
		t.Errorf("Expected avg childrens 200000.2, got %f", data.AvgChildrens)
	}

	if data.MinPopulation != 100000 {
		t.Errorf("Expected min population 100000, got %d", data.MinPopulation)
	}

	if data.MaxPopulation != 5000000 {
		t.Errorf("Expected max population 5000000, got %d", data.MaxPopulation)
	}
}

func TestSettlementTypeDataZeroValues(t *testing.T) {
	data := SettlementTypeData{}

	if data.Type != "" {
		t.Errorf("Expected empty type, got %s", data.Type)
	}

	if data.AvgPopulation != 0 {
		t.Errorf("Expected zero avg population, got %f", data.AvgPopulation)
	}

	if data.MinPopulation != 0 {
		t.Errorf("Expected zero min population, got %d", data.MinPopulation)
	}

	if data.MaxPopulation != 0 {
		t.Errorf("Expected zero max population, got %d", data.MaxPopulation)
	}
}

func TestSettlementTypeDataMultipleTypes(t *testing.T) {
	types := []SettlementTypeData{
		{Type: "город", AvgPopulation: 5000000, MinPopulation: 100000, MaxPopulation: 12000000},
		{Type: "деревня", AvgPopulation: 500, MinPopulation: 100, MaxPopulation: 2000},
		{Type: "село", AvgPopulation: 5000, MinPopulation: 500, MaxPopulation: 50000},
	}

	if len(types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(types))
	}

	if types[0].Type != "город" {
		t.Errorf("Expected first type 'город', got %s", types[0].Type)
	}

	if types[0].AvgPopulation < types[1].AvgPopulation {
		t.Errorf("Cities should have higher avg population than villages")
	}
}
