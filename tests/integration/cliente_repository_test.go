package integration

import (
	"context"
	"testing"

	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClienteRepoTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&model.Cliente{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestClienteRepository_Create(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	cliente := &model.Cliente{
		Nome:  "Test Cliente",
		Email: "test@example.com",
		CPF:   "12345678901",
	}

	err := repo.Create(context.Background(), cliente)
	assert.NoError(t, err)
	assert.NotZero(t, cliente.ID)
}

func TestClienteRepository_FindAll(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	// Create test data
	db.Create(&model.Cliente{Nome: "Cliente 1", Email: "c1@test.com", CPF: "11111111111"})
	db.Create(&model.Cliente{Nome: "Cliente 2", Email: "c2@test.com", CPF: "22222222222"})

	clientes, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, clientes, 2)
}

func TestClienteRepository_FindAll_DBError(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	// Close database to simulate error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	_, err := repo.FindAll(context.Background())
	assert.Error(t, err)
}

func TestClienteRepository_FindByID(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	created := &model.Cliente{Nome: "Test", Email: "test@test.com", CPF: "12345678901"}
	db.Create(created)

	cliente, err := repo.FindByID(context.Background(), created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, cliente)
	assert.Equal(t, created.Nome, cliente.Nome)
}

func TestClienteRepository_FindByID_NotFound(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	cliente, err := repo.FindByID(context.Background(), 9999)
	assert.NoError(t, err)
	assert.Nil(t, cliente)
}

func TestClienteRepository_FindByID_DBError(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	// Close database to simulate error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	_, err := repo.FindByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestClienteRepository_FindByName(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	db.Create(&model.Cliente{Nome: "João Silva", Email: "joao@test.com", CPF: "11111111111"})
	db.Create(&model.Cliente{Nome: "Maria Santos", Email: "maria@test.com", CPF: "22222222222"})

	clientes, err := repo.FindByName(context.Background(), "João")
	assert.NoError(t, err)
	assert.Len(t, clientes, 1)
	assert.Equal(t, "João Silva", clientes[0].Nome)
}

func TestClienteRepository_Update(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	cliente := &model.Cliente{Nome: "Original", Email: "original@test.com", CPF: "12345678901"}
	db.Create(cliente)

	cliente.Nome = "Updated"
	err := repo.Update(context.Background(), cliente)
	assert.NoError(t, err)

	var updated model.Cliente
	db.First(&updated, cliente.ID)
	assert.Equal(t, "Updated", updated.Nome)
}

func TestClienteRepository_Delete(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	cliente := &model.Cliente{Nome: "To Delete", Email: "delete@test.com", CPF: "12345678901"}
	db.Create(cliente)

	err := repo.Delete(context.Background(), cliente.ID)
	assert.NoError(t, err)

	var count int64
	db.Model(&model.Cliente{}).Where("id = ?", cliente.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestClienteRepository_Delete_NotFound(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	err := repo.Delete(context.Background(), 9999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestClienteRepository_Delete_DBError(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	// Close database to simulate error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	err := repo.Delete(context.Background(), 1)
	assert.Error(t, err)
}

func TestClienteRepository_Count(t *testing.T) {
	db := setupClienteRepoTestDB(t)
	repo := repository.NewClienteRepositorySQLite(db)

	db.Create(&model.Cliente{Nome: "Cliente 1", Email: "c1@test.com", CPF: "11111111111"})
	db.Create(&model.Cliente{Nome: "Cliente 2", Email: "c2@test.com", CPF: "22222222222"})

	count, err := repo.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
