package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danmaciel/api/config"
	"github.com/stretchr/testify/assert"
)

func TestInitDatabase_Success(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	db, err := config.InitDatabase(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Verify database file was created
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)

	// Close database
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestInitDatabase_CreatesDirectory(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "subdir", "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	db, err := config.InitDatabase(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Verify directory was created
	_, err = os.Stat(filepath.Dir(dbPath))
	assert.NoError(t, err)

	// Verify database file was created
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)

	// Close database
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestInitDatabase_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "level1", "level2", "level3", "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	db, err := config.InitDatabase(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Verify all nested directories were created
	_, err = os.Stat(filepath.Join(tmpDir, "level1", "level2", "level3"))
	assert.NoError(t, err)

	// Close database
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestInitDatabase_InvalidPath(t *testing.T) {
	// Use an invalid path that cannot be created (root directory with no permissions)
	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: "/root/no-permission/test.db",
	}

	db, err := config.InitDatabase(cfg)

	// Should return error when directory cannot be created
	// Note: This might not fail on all systems depending on permissions
	if err == nil {
		// If it succeeded (unlikely), clean up
		if db != nil {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}
		t.Skip("Test skipped: system allows creation in restricted directory")
	}

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestInitDatabase_MultipleConnections(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	// Initialize first connection
	db1, err := config.InitDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db1)

	// Initialize second connection to same database
	db2, err := config.InitDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db2)

	// Close both connections
	sqlDB1, _ := db1.DB()
	sqlDB1.Close()

	sqlDB2, _ := db2.DB()
	sqlDB2.Close()
}

func TestInitDatabase_ExistingDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	// Initialize database first time
	db1, err := config.InitDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db1)

	sqlDB1, _ := db1.DB()
	sqlDB1.Close()

	// Initialize database second time (should work with existing file)
	db2, err := config.InitDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db2)

	sqlDB2, _ := db2.DB()
	sqlDB2.Close()
}

func TestInitDatabase_EmptyPath(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: "",
	}

	db, err := config.InitDatabase(cfg)

	// Should handle empty path gracefully
	if err == nil && db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
	// Either succeeds with default or fails - both are acceptable
	// Main goal is to not panic
}

func TestInitDatabase_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: "./relative/path/test.db",
	}

	db, err := config.InitDatabase(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Verify relative path works
	_, err = os.Stat("./relative/path/test.db")
	assert.NoError(t, err)

	// Close database
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestInitDatabase_TablesCreated(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.DatabaseConfig{
		Driver:   "sqlite",
		FilePath: dbPath,
	}

	db, err := config.InitDatabase(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Check that tables were created by AutoMigrate
	assert.True(t, db.Migrator().HasTable("clientes"))
	assert.True(t, db.Migrator().HasTable("produtos"))
	assert.True(t, db.Migrator().HasTable("pedidos"))
	assert.True(t, db.Migrator().HasTable("pedido_produtos"))

	// Close database
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
