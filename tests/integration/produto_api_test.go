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

func setupProdutoTestDB(t *testing.T) *gorm.DB {
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

func setupProdutoTestRouter(db *gorm.DB) (*controller.ClienteController, *controller.ProdutoController, *controller.PedidoController) {
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

func TestCreateProduto_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	ativo := true
	reqBody := dto.CreateProdutoRequest{
		Nome:      "Notebook Dell",
		Descricao: "Notebook Dell Inspiron 15",
		Preco:     2999.99,
		Estoque:   10,
		SKU:       "NB-DELL-001",
		Categoria: "Eletrônicos",
		Ativo:     &ativo,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, reqBody.Nome, response.Nome)
	assert.Equal(t, reqBody.SKU, response.SKU)
	assert.NotZero(t, response.ID)
}

func TestGetAllProdutos_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Produto{Nome: "Produto 1", SKU: "PROD-001", Preco: 100.00})
	db.Create(&model.Produto{Nome: "Produto 2", SKU: "PROD-002", Preco: 200.00})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
}

func TestGetProdutoByID_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	produto := &model.Produto{Nome: "Produto Teste", SKU: "PROD-TEST", Preco: 150.00}
	db.Create(produto)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, produto.Nome, response.Nome)
}

func TestCountProdutos_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Produto{Nome: "Produto 1", SKU: "PROD-001", Preco: 100.00})
	db.Create(&model.Produto{Nome: "Produto 2", SKU: "PROD-002", Preco: 200.00})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.CountResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, int64(2), response.Count)
}

func TestUpdateProduto_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	produto := &model.Produto{Nome: "Produto Original", SKU: "PROD-ORIG", Preco: 100.00}
	db.Create(produto)

	reqBody := dto.UpdateProdutoRequest{
		Nome:  "Produto Atualizado",
		Preco: 150.00,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/produtos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, "Produto Atualizado", response.Nome)
	assert.Equal(t, 150.00, response.Preco)
}

func TestDeleteProduto_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	produto := &model.Produto{Nome: "Produto Para Deletar", SKU: "PROD-DEL", Preco: 100.00}
	db.Create(produto)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/produtos/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestFindByNomeProduto_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Produto{Nome: "Notebook Dell", SKU: "NB-DELL", Preco: 2999.99})
	db.Create(&model.Produto{Nome: "Mouse Logitech", SKU: "MS-LOG", Preco: 99.99})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/nome/Notebook", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 1)
	assert.Equal(t, "Notebook Dell", response[0].Nome)
}

func TestFindByCategoriaProduto_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Create test data
	db.Create(&model.Produto{Nome: "Notebook Dell", SKU: "NB-DELL", Preco: 2999.99, Categoria: "Eletrônicos"})
	db.Create(&model.Produto{Nome: "Mouse Logitech", SKU: "MS-LOG", Preco: 99.99, Categoria: "Eletrônicos"})
	db.Create(&model.Produto{Nome: "Mesa", SKU: "MESA-001", Preco: 500.00, Categoria: "Móveis"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/categoria/Eletrônicos", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []dto.ProdutoResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Len(t, response, 2)
}

// Testes de erro para aumentar cobertura
func TestCreateProduto_InvalidJSON_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateProduto_ValidationError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.CreateProdutoRequest{
		Nome:  "", // Invalid
		Preco: -10, // Invalid
		SKU:   "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)
}

func TestGetProdutoByID_InvalidID_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetProdutoByID_NotFound_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/9999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateProduto_InvalidJSON_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/produtos/1", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProduto_InvalidID_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdateProdutoRequest{Nome: "Test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/produtos/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProduto_NotFound_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	reqBody := dto.UpdateProdutoRequest{Nome: "Test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/produtos/9999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteProduto_InvalidID_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/produtos/invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAllProdutos_DBError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCountProdutos_DBError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestFindByNomeProduto_DBError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/nome/Test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestFindByCategoriaProduto_DBError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos/categoria/Test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteProduto_DBError_Integration(t *testing.T) {
	db := setupProdutoTestDB(t)
	clienteCtrl, produtoCtrl, pedidoCtrl := setupProdutoTestRouter(db)
	router := controller.SetupRouter(clienteCtrl, produtoCtrl, pedidoCtrl)

	// Close database to trigger error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/produtos/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
