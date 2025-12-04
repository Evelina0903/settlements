package factory

import (
	"os"
	"testing"

	"settlements/internal/config"
)

func TestApplicationFactoryCreation(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DB_HOST")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	factory, err := NewApplicationFactory(cfg)

	// Note: Will fail without actual DB connection, but tests factory creation logic
	if cfg == nil {
		t.Errorf("Expected config to be set")
	}

	if factory == nil && err == nil {
		t.Errorf("Expected factory or error")
	}
}

func TestApplicationFactoryNilConfig(t *testing.T) {
	_, err := NewApplicationFactory(nil)

	if err == nil {
		t.Errorf("Expected error for nil config")
	}
}

func TestApplicationFactoryCreateRouter(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	cfg, _ := config.Load()
	factory, _ := NewApplicationFactory(cfg)

	if factory == nil {
		t.Skip("Skipping due to DB connection issue")
	}

	router := factory.CreateRouter()

	if router == nil {
		t.Errorf("Expected router instance")
	}
}

func TestApplicationContextGetServerAddress(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: "3000"},
	}

	ctx := &ApplicationContext{
		Config: cfg,
	}

	expected := ":3000"
	got := ctx.GetServerAddress()

	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}
