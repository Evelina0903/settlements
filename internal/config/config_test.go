package config

import (
	"os"
	"testing"
)

func TestLoadWithDefaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Server.Port != "3000" {
		t.Errorf("Expected default port 3000, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "postgres" {
		t.Errorf("Expected default host 'postgres', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default port 5432, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "postgres" {
		t.Errorf("Expected default user 'postgres', got %s", cfg.Database.User)
	}

	if cfg.Database.Password != "postgres" {
		t.Errorf("Expected default password 'postgres', got %s", cfg.Database.Password)
	}

	if cfg.Database.Name != "database" {
		t.Errorf("Expected default name 'database', got %s", cfg.Database.Name)
	}
}

func TestLoadWithCustomValues(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Expected port 5433, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "testuser" {
		t.Errorf("Expected user 'testuser', got %s", cfg.Database.User)
	}

	if cfg.Database.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got %s", cfg.Database.Password)
	}

	if cfg.Database.Name != "testdb" {
		t.Errorf("Expected name 'testdb', got %s", cfg.Database.Name)
	}
}

func TestLoadInvalidDBPort(t *testing.T) {
	os.Setenv("DB_PORT", "not_a_number")
	defer os.Unsetenv("DB_PORT")

	_, err := Load()

	if err == nil {
		t.Fatalf("Expected error for invalid DB_PORT, got nil")
	}
}

func TestDatabaseConfigDSN(t *testing.T) {
	dbCfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		Name:     "testdb",
	}

	dsn := dbCfg.DSN()
	expected := "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable TimeZone=UTC"

	if dsn != expected {
		t.Errorf("Expected DSN %q, got %q", expected, dsn)
	}
}
