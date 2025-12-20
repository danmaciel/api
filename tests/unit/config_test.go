package unit

import (
	"os"
	"testing"

	"github.com/danmaciel/api/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("DB_FILE_PATH")

	cfg := config.Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "sqlite", cfg.Database.Driver)
	assert.Equal(t, "./database/api.db", cfg.Database.FilePath)
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_FILE_PATH", "/custom/path/db.db")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("DB_FILE_PATH")
	}()

	cfg := config.Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "postgres", cfg.Database.Driver)
	assert.Equal(t, "/custom/path/db.db", cfg.Database.FilePath)
}

func TestLoad_InvalidPort(t *testing.T) {
	os.Setenv("SERVER_PORT", "invalid")
	defer os.Unsetenv("SERVER_PORT")

	cfg := config.Load()

	// Should use default value when port is invalid
	assert.Equal(t, 8080, cfg.Server.Port)
}

func TestLoad_EmptyEnvironmentVariables(t *testing.T) {
	os.Setenv("SERVER_PORT", "")
	os.Setenv("SERVER_HOST", "")
	os.Setenv("DB_DRIVER", "")
	os.Setenv("DB_FILE_PATH", "")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("DB_FILE_PATH")
	}()

	cfg := config.Load()

	// Should use defaults when env vars are empty strings
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "sqlite", cfg.Database.Driver)
	assert.Equal(t, "./database/api.db", cfg.Database.FilePath)
}

func TestGetServerAddress(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 3000,
		},
	}

	address := cfg.GetServerAddress()
	assert.Equal(t, "localhost:3000", address)
}

func TestGetServerAddress_DefaultValues(t *testing.T) {
	cfg := config.Load()
	address := cfg.GetServerAddress()
	assert.Equal(t, "0.0.0.0:8080", address)
}

func TestGetServerAddress_CustomPort(t *testing.T) {
	os.Setenv("SERVER_PORT", "9000")
	defer os.Unsetenv("SERVER_PORT")

	cfg := config.Load()
	address := cfg.GetServerAddress()
	assert.Equal(t, "0.0.0.0:9000", address)
}

func TestConfig_StructureIntegrity(t *testing.T) {
	cfg := config.Load()

	// Verify struct fields are properly initialized
	assert.NotZero(t, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Server.Host)
	assert.NotEmpty(t, cfg.Database.Driver)
	assert.NotEmpty(t, cfg.Database.FilePath)
}

func TestServerConfig_Types(t *testing.T) {
	cfg := config.Load()

	// Verify field types
	assert.IsType(t, 0, cfg.Server.Port)
	assert.IsType(t, "", cfg.Server.Host)
}

func TestDatabaseConfig_Types(t *testing.T) {
	cfg := config.Load()

	// Verify field types
	assert.IsType(t, "", cfg.Database.Driver)
	assert.IsType(t, "", cfg.Database.FilePath)
}

func TestLoad_PartialEnvironmentVariables(t *testing.T) {
	// Set only some environment variables
	os.Setenv("SERVER_PORT", "5000")
	os.Setenv("DB_DRIVER", "mysql")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_DRIVER")
	}()

	cfg := config.Load()

	// Check that set values are used
	assert.Equal(t, 5000, cfg.Server.Port)
	assert.Equal(t, "mysql", cfg.Database.Driver)

	// Check that unset values use defaults
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "./database/api.db", cfg.Database.FilePath)
}
