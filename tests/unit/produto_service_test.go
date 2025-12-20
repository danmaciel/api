package unit

import (
	"context"
	"testing"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProdutoRepository is a mock implementation of ProdutoRepository
type MockProdutoRepository struct {
	mock.Mock
}

func (m *MockProdutoRepository) Create(ctx context.Context, produto *model.Produto) error {
	args := m.Called(ctx, produto)
	return args.Error(0)
}

func (m *MockProdutoRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) FindByID(ctx context.Context, id uint) (*model.Produto, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) FindByName(ctx context.Context, nome string) ([]model.Produto, error) {
	args := m.Called(ctx, nome)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) FindBySKU(ctx context.Context, sku string) (*model.Produto, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) FindByCategoria(ctx context.Context, categoria string) ([]model.Produto, error) {
	args := m.Called(ctx, categoria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) Update(ctx context.Context, produto *model.Produto) error {
	args := m.Called(ctx, produto)
	return args.Error(0)
}

func (m *MockProdutoRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProdutoRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Test cases
func TestProdutoService_Create_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	ativo := true
	req := &dto.CreateProdutoRequest{
		Nome:      "Notebook Dell",
		Descricao: "Notebook Dell Inspiron 15",
		Preco:     2999.99,
		Estoque:   10,
		SKU:       "NB-DELL-001",
		Categoria: "Eletrônicos",
		Ativo:     &ativo,
	}

	// Mock FindBySKU to verify SKU doesn't exist
	mockRepo.On("FindBySKU", mock.Anything, "NB-DELL-001").Return((*model.Produto)(nil), nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Produto")).Return(nil)

	result, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Nome, result.Nome)
	assert.Equal(t, req.SKU, result.SKU)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Create_ValidationError(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	req := &dto.CreateProdutoRequest{
		Nome:  "", // Invalid: empty name
		Preco: 100.00,
		SKU:   "TEST-001",
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Nome")
}

func TestProdutoService_Create_ValidationError_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	req := &dto.CreateProdutoRequest{
		Nome:  "Produto Teste",
		Preco: -10.00, // Invalid: negative price
		SKU:   "TEST-001",
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Preco")
}

func TestProdutoService_FindAll_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	expectedProdutos := []model.Produto{
		{ID: 1, Nome: "Produto 1", SKU: "PROD-001", Preco: 100.00},
		{ID: 2, Nome: "Produto 2", SKU: "PROD-002", Preco: 200.00},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expectedProdutos, nil)

	result, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByID_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	expectedProduto := &model.Produto{
		ID:        1,
		Nome:      "Notebook Dell",
		SKU:       "NB-DELL-001",
		Preco:     2999.99,
		Categoria: "Eletrônicos",
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedProduto, nil)

	result, err := svc.FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduto.Nome, result.Nome)
	assert.Equal(t, expectedProduto.SKU, result.SKU)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByID_NotFound(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Produto)(nil), assert.AnError)

	result, err := svc.FindByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByName_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	expectedProdutos := []model.Produto{
		{ID: 1, Nome: "Notebook Dell", SKU: "NB-DELL-001", Preco: 2999.99},
	}

	mockRepo.On("FindByName", mock.Anything, "Notebook").Return(expectedProdutos, nil)

	result, err := svc.FindByName(context.Background(), "Notebook")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Notebook Dell", result[0].Nome)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByCategoria_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	expectedProdutos := []model.Produto{
		{ID: 1, Nome: "Notebook Dell", SKU: "NB-DELL-001", Preco: 2999.99, Categoria: "Eletrônicos"},
		{ID: 2, Nome: "Mouse Logitech", SKU: "MS-LOG-001", Preco: 99.99, Categoria: "Eletrônicos"},
	}

	mockRepo.On("FindByCategoria", mock.Anything, "Eletrônicos").Return(expectedProdutos, nil)

	result, err := svc.FindByCategoria(context.Background(), "Eletrônicos")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Update_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	existingProduto := &model.Produto{
		ID:        1,
		Nome:      "Notebook Dell",
		SKU:       "NB-DELL-001",
		Preco:     2999.99,
		Categoria: "Eletrônicos",
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existingProduto, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Produto")).Return(nil)

	req := &dto.UpdateProdutoRequest{
		Nome:  "Notebook Dell Atualizado",
		Preco: 2799.99,
	}

	result, err := svc.Update(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Notebook Dell Atualizado", result.Nome)
	assert.Equal(t, 2799.99, result.Preco)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Produto)(nil), assert.AnError)

	req := &dto.UpdateProdutoRequest{
		Nome: "Produto Atualizado",
	}

	result, err := svc.Update(context.Background(), 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Delete_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	// Mock FindByID to verify produto exists
	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Produto{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Count_Success(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("Count", mock.Anything).Return(int64(25), nil)

	count, err := svc.Count(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(25), count)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindAll_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindAll", mock.Anything).Return([]model.Produto(nil), assert.AnError)

	result, err := svc.FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByName_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindByName", mock.Anything, "Test").Return([]model.Produto(nil), assert.AnError)

	result, err := svc.FindByName(context.Background(), "Test")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_FindByCategoria_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindByCategoria", mock.Anything, "Test").Return([]model.Produto(nil), assert.AnError)

	result, err := svc.FindByCategoria(context.Background(), "Test")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Update_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	existingProduto := &model.Produto{
		ID:    1,
		Nome:  "Produto Teste",
		SKU:   "TEST-001",
		Preco: 100.00,
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existingProduto, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Produto")).Return(assert.AnError)

	req := &dto.UpdateProdutoRequest{
		Nome: "Updated",
	}

	result, err := svc.Update(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Delete_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Produto{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, uint(1)).Return(assert.AnError)

	err := svc.Delete(context.Background(), 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Count_Error(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	mockRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)

	count, err := svc.Count(context.Background())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	ativo := true
	req := &dto.CreateProdutoRequest{
		Nome:    "Produto Teste",
		Preco:   100.00,
		SKU:     "TEST-001",
		Estoque: 10,
		Ativo:   &ativo,
	}

	mockRepo.On("FindBySKU", mock.Anything, "TEST-001").Return((*model.Produto)(nil), nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Produto")).Return(assert.AnError)

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestProdutoService_Create_SKUAlreadyExists(t *testing.T) {
	mockRepo := new(MockProdutoRepository)
	svc := service.NewProdutoService(mockRepo)

	ativo := true
	req := &dto.CreateProdutoRequest{
		Nome:  "Produto Teste",
		Preco: 100.00,
		SKU:   "TEST-001",
		Ativo: &ativo,
	}

	existingProduto := &model.Produto{ID: 1, SKU: "TEST-001"}
	mockRepo.On("FindBySKU", mock.Anything, "TEST-001").Return(existingProduto, nil)

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "já cadastrado")
	mockRepo.AssertExpectations(t)
}
