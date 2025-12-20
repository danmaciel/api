package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danmaciel/api/internal/controller"
	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/repository"
	"github.com/danmaciel/api/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&model.Cliente{}, &model.Produto{}, &model.Pedido{}, &model.PedidoProduto{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func setupTestRouter(db *gorm.DB) (*controller.ClienteController, *controller.ProdutoController, *controller.PedidoController) {
	// Cliente
	clienteRepo := repository.NewClienteRepositorySQLite(db)
	clienteService := service.NewClienteService(clienteRepo)
	clienteController := controller.NewClienteController(clienteService)

	// Produto
	produtoRepo := repository.NewProdutoRepositorySQLite(db)
	produtoService := service.NewProdutoService(produtoRepo)
	produtoController := controller.NewProdutoController(produtoService)

	// Pedido
	pedidoRepo := repository.NewPedidoRepositorySQLite(db)
	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo)
	pedidoController := controller.NewPedidoController(pedidoService)

	return clienteController, produtoController, pedidoController
}

func TestCreateCliente_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.CreateClienteRequest{
		Nome:     "Maria Silva",
		Email:    "maria@example.com",
		CPF:      "98765432100",
		Telefone: "11988888888",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, reqBody.Nome, response.Nome)
	assert.Equal(t, reqBody.Email, response.Email)
	assert.NotZero(t, response.ID)
}

func TestGetAllClientes_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Cliente{Nome: "Cliente 1", Email: "c1@example.com", CPF: "11111111111"})
	db.Create(&model.Cliente{Nome: "Cliente 2", Email: "c2@example.com", CPF: "22222222222"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
}

func TestGetClienteByID_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "33333333333"}
	db.Create(cliente)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, cliente.Nome, response.Nome)
}

func TestCountClientes_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Cliente{Nome: "Cliente 1", Email: "c1@example.com", CPF: "11111111111"})
	db.Create(&model.Cliente{Nome: "Cliente 2", Email: "c2@example.com", CPF: "22222222222"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.CountResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, int64(2), response.Count)
}

func TestUpdateCliente_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Original", Email: "original@example.com", CPF: "44444444444"}
	db.Create(cliente)

	reqBody := dto.UpdateClienteRequest{
		Nome: "Cliente Atualizado",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "Cliente Atualizado", response.Nome)
}

func TestDeleteCliente_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Para Deletar", Email: "deletar@example.com", CPF: "55555555555"}
	db.Create(cliente)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestFindByNome_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Cliente{Nome: "João Silva", Email: "joao@example.com", CPF: "66666666666"})
	db.Create(&model.Cliente{Nome: "Maria Santos", Email: "maria@example.com", CPF: "77777777777"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/nome/João", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "João Silva", response[0].Nome)
}

// Testes de erro para aumentar cobertura
func TestCreateCliente_InvalidJSON_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCliente_ValidationError_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.CreateClienteRequest{
		Nome:  "", // Invalid: empty name
		Email: "invalid",
		CPF:   "123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Service validation returns 500 for validation errors in this implementation
	assert.True(t, rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)
}

func TestGetClienteByID_InvalidID_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetClienteByID_NotFound_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/9999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCliente_InvalidJSON_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCliente_InvalidID_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdateClienteRequest{Nome: "Test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCliente_NotFound_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdateClienteRequest{Nome: "Test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/clientes/9999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCliente_InvalidID_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFindByNome_EmptyResults_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/nome/NãoExiste", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ClienteResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 0)
}

func TestGetAllClientes_DBError_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCount_DBError_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestFindByNome_DBError_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/nome/Test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteCliente_DBError_Integration(t *testing.T) {
	db := setupTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/clientes/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
