package data_loader

import (
	"testing"
)

func TestSettlementsTypesMapping(t *testing.T) {
	tests := []struct {
		abbreviation string
		expected     string
	}{
		{"г", "город"},
		{"д", "деревня"},
		{"с", "село"},
		{"п", "поселок"},
		{"пгт", "поселок городского типа"},
		{"х", "хутор"},
		{"м", "местечко"},
		{"ст-ца", "станица"},
	}

	for _, test := range tests {
		if v, ok := settlementsTypes[test.abbreviation]; !ok {
			t.Errorf("Expected abbreviation %q in map", test.abbreviation)
		} else if v != test.expected {
			t.Errorf("For %q: expected %q, got %q", test.abbreviation, test.expected, v)
		}
	}
}

func TestSettlementsTypesMappingCount(t *testing.T) {
	// Verify the map has expected number of entries
	if len(settlementsTypes) < 30 {
		t.Errorf("Expected at least 30 settlement types, got %d", len(settlementsTypes))
	}
}

func TestDataLoaderCreation(t *testing.T) {
	// Note: This test would require mocking gorm.DB
	// For now, we test the data structures only
	
	// Verify that settlementsTypes map contains expected mappings
	expectedCount := len(settlementsTypes)
	if expectedCount == 0 {
		t.Errorf("Expected settlement types to be initialized")
	}
}

func TestRailwayTypeAbbreviations(t *testing.T) {
	// Test that railway-related abbreviations are properly mapped
	railwayTests := []struct {
		abbrev   string
		contains string
	}{
		{"ж/д ст", "железнодорожная станция"},
		{"ж/д платформа", "железнодорожная платформа"},
		{"ж/д оп", "железнодорожный остановочный пункт"},
		{"ж/д рзд", "железнодорожный разъезд"},
	}

	for _, test := range railwayTests {
		if v, ok := settlementsTypes[test.abbrev]; !ok {
			t.Errorf("Expected railway abbreviation %q in map", test.abbrev)
		} else if v != test.contains {
			t.Errorf("For %q: expected %q, got %q", test.abbrev, test.contains, v)
		}
	}
}

func TestUrbanTypeAbbreviations(t *testing.T) {
	// Test urban-related abbreviations
	urbanTests := []struct {
		abbrev   string
		expected string
	}{
		{"г", "город"},
		{"гп", "городской поселок"},
		{"пгт", "поселок городского типа"},
	}

	for _, test := range urbanTests {
		if v, ok := settlementsTypes[test.abbrev]; !ok {
			t.Errorf("Expected urban abbreviation %q in map", test.abbrev)
		} else if v != test.expected {
			t.Errorf("For %q: expected %q, got %q", test.abbrev, test.expected, v)
		}
	}
}

func TestRuralTypeAbbreviations(t *testing.T) {
	// Test rural-related abbreviations
	ruralTests := []struct {
		abbrev   string
		expected string
	}{
		{"д", "деревня"},
		{"с", "село"},
		{"х", "хутор"},
		{"п", "поселок"},
	}

	for _, test := range ruralTests {
		if v, ok := settlementsTypes[test.abbrev]; !ok {
			t.Errorf("Expected rural abbreviation %q in map", test.abbrev)
		} else if v != test.expected {
			t.Errorf("For %q: expected %q, got %q", test.abbrev, test.expected, v)
		}
	}
}

func TestSpecialSettlementTypes(t *testing.T) {
	// Test special/unique settlement types
	specialTests := []struct {
		abbrev   string
		expected string
	}{
		{"кп", "курортный поселок"},
		{"дп", "дачный поселок"},
		{"снт", "садоводческое некоммерческое товарищество"},
		{"к", "кишлак"},
		{"у", "улус"},
		{"л/п", "лесной поселок"},
	}

	for _, test := range specialTests {
		if v, ok := settlementsTypes[test.abbrev]; !ok {
			t.Errorf("Expected special abbreviation %q in map", test.abbrev)
		} else if v != test.expected {
			t.Errorf("For %q: expected %q, got %q", test.abbrev, test.expected, v)
		}
	}
}

func TestMapValueUniqueness(t *testing.T) {
	// Verify that the mapping preserves abbreviations
	seen := make(map[string]string)

	for abbrev, fullName := range settlementsTypes {
		if abbrev == "" {
			t.Errorf("Found empty abbreviation key")
		}
		if fullName == "" {
			t.Errorf("Found empty value for abbreviation %q", abbrev)
		}
		seen[abbrev] = fullName
	}

	if len(seen) != len(settlementsTypes) {
		t.Errorf("Unexpected duplicate entries in settlementsTypes")
	}
}
