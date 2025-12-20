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

func setupPedidoTestDB(t *testing.T) *gorm.DB {
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

func setupPedidoTestRouter(db *gorm.DB) (*controller.ClienteController, *controller.ProdutoController, *controller.PedidoController) {
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

func TestCreatePedido_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "João Silva", Email: "joao@example.com", CPF: "12345678901"}
	db.Create(cliente)

	produto := &model.Produto{Nome: "Notebook", SKU: "NB-001", Preco: 2999.99, Estoque: 10, Ativo: true}
	db.Create(produto)

	reqBody := dto.CreatePedidoRequest{
		ClienteID: cliente.ID,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: produto.ID, Quantidade: 2},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pedidos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, cliente.ID, response.ClienteID)
	assert.NotZero(t, response.ID)
	assert.NotZero(t, response.ValorTotal)
}

func TestGetAllPedidos_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 100.00, Status: "pendente"})
	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 200.00, Status: "pago"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
}

func TestGetPedidoByID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	pedido := &model.Pedido{ClienteID: cliente.ID, ValorTotal: 150.00, Status: "pendente"}
	db.Create(pedido)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, pedido.ClienteID, response.ClienteID)
}

func TestCountPedidos_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 100.00, Status: "pendente"})
	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 200.00, Status: "pago"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.CountResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, int64(2), response.Count)
}

func TestUpdatePedidoStatus_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	pedido := &model.Pedido{ClienteID: cliente.ID, ValorTotal: 100.00, Status: "pendente"}
	db.Create(pedido)

	reqBody := dto.UpdatePedidoRequest{
		Status: "pago",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pedidos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "pago", response.Status)
}

func TestDeletePedido_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	pedido := &model.Pedido{ClienteID: cliente.ID, ValorTotal: 100.00, Status: "pendente"}
	db.Create(pedido)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pedidos/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestFindByClienteID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente1 := &model.Cliente{Nome: "Cliente 1", Email: "cliente1@example.com", CPF: "11111111111"}
	db.Create(cliente1)

	cliente2 := &model.Cliente{Nome: "Cliente 2", Email: "cliente2@example.com", CPF: "22222222222"}
	db.Create(cliente2)

	db.Create(&model.Pedido{ClienteID: cliente1.ID, ValorTotal: 100.00, Status: "pendente"})
	db.Create(&model.Pedido{ClienteID: cliente1.ID, ValorTotal: 200.00, Status: "pago"})
	db.Create(&model.Pedido{ClienteID: cliente2.ID, ValorTotal: 300.00, Status: "pendente"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/cliente/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
	assert.Equal(t, cliente1.ID, response[0].ClienteID)
}

func TestFindByStatus_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	cliente := &model.Cliente{Nome: "Cliente Teste", Email: "teste@example.com", CPF: "11111111111"}
	db.Create(cliente)

	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 100.00, Status: "pendente"})
	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 200.00, Status: "pendente"})
	db.Create(&model.Pedido{ClienteID: cliente.ID, ValorTotal: 300.00, Status: "pago"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/status/pendente", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.PedidoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
	assert.Equal(t, "pendente", response[0].Status)
}

// Testes de erro para aumentar cobertura
func TestCreatePedido_InvalidJSON_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pedidos", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreatePedido_ValidationError_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.CreatePedidoRequest{
		ClienteID: 0, // Invalid
		Itens:     []dto.CreateItemPedidoRequest{},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pedidos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)
}

func TestGetPedidoByID_InvalidID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetPedidoByID_NotFound_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/9999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdatePedidoStatus_InvalidJSON_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pedidos/1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePedidoStatus_InvalidID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdatePedidoRequest{Status: "pago"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pedidos/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePedidoStatus_NotFound_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdatePedidoRequest{Status: "pago"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/pedidos/9999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeletePedido_InvalidID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pedidos/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFindByClienteID_InvalidID_Integration(t *testing.T) {
	db := setupPedidoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupPedidoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pedidos/cliente/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
