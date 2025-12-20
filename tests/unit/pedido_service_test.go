package unit

import (
	"context"
	"testing"
	"time"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPedidoRepository is a mock implementation of PedidoRepository
type MockPedidoRepository struct {
	mock.Mock
}

func (m *MockPedidoRepository) Create(ctx context.Context, pedido *model.Pedido) error {
	args := m.Called(ctx, pedido)
	return args.Error(0)
}

func (m *MockPedidoRepository) FindAll(ctx context.Context) ([]model.Pedido, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) FindByID(ctx context.Context, id uint) (*model.Pedido, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) FindByClienteID(ctx context.Context, clienteID uint) ([]model.Pedido, error) {
	args := m.Called(ctx, clienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) FindByStatus(ctx context.Context, status string) ([]model.Pedido, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) Update(ctx context.Context, pedido *model.Pedido) error {
	args := m.Called(ctx, pedido)
	return args.Error(0)
}

func (m *MockPedidoRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPedidoRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Test cases
func TestPedidoService_Create_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	// Mock cliente exists
	mockClienteRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Cliente{
		ID:    1,
		Nome:  "João Silva",
		Email: "joao@example.com",
	}, nil)

	// Mock produto exists
	mockProdutoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Produto{
		ID:      1,
		Nome:    "Notebook",
		Preco:   2999.99,
		Estoque: 10,
		SKU:     "NB-001",
		Ativo:   true,
	}, nil)

	mockPedidoRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Pedido")).Run(func(args mock.Arguments) {
		pedido := args.Get(1).(*model.Pedido)
		pedido.ID = 1
	}).Return(nil)

	// Mock FindByID after create to return the complete pedido
	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Pedido{
		ID:         1,
		ClienteID:  1,
		ValorTotal: 5999.98,
		Status:     "pendente",
		DataPedido: time.Now(),
		Itens: []model.PedidoProduto{
			{
				ProdutoID:     1,
				Quantidade:    2,
				PrecoUnitario: 2999.99,
				Subtotal:      5999.98,
			},
		},
	}, nil)

	req := &dto.CreatePedidoRequest{
		ClienteID: 1,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: 1, Quantidade: 2},
		},
	}

	result, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ClienteID)
	mockPedidoRepo.AssertExpectations(t)
	mockClienteRepo.AssertExpectations(t)
	mockProdutoRepo.AssertExpectations(t)
}

func TestPedidoService_Create_ValidationError(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	req := &dto.CreatePedidoRequest{
		ClienteID: 0, // Invalid: missing cliente_id
		Itens:     []dto.CreateItemPedidoRequest{},
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ClienteID")
}

func TestPedidoService_Create_ClienteNotFound(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	// Mock cliente not found - return error
	mockClienteRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Cliente)(nil), assert.AnError)

	req := &dto.CreatePedidoRequest{
		ClienteID: 999,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: 1, Quantidade: 1},
		},
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "não encontrado")
	mockClienteRepo.AssertExpectations(t)
}

func TestPedidoService_Create_ProdutoNotFound(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	// Mock cliente exists
	mockClienteRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Cliente{
		ID: 1,
	}, nil)

	// Mock produto not found - return error
	mockProdutoRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Produto)(nil), assert.AnError)

	req := &dto.CreatePedidoRequest{
		ClienteID: 1,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: 999, Quantidade: 1},
		},
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "não encontrado")
	mockClienteRepo.AssertExpectations(t)
	mockProdutoRepo.AssertExpectations(t)
}

func TestPedidoService_FindAll_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	expectedPedidos := []model.Pedido{
		{ID: 1, ClienteID: 1, ValorTotal: 299.99, Status: "pendente"},
		{ID: 2, ClienteID: 2, ValorTotal: 499.99, Status: "pago"},
	}

	mockPedidoRepo.On("FindAll", mock.Anything).Return(expectedPedidos, nil)

	result, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByID_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	expectedPedido := &model.Pedido{
		ID:         1,
		ClienteID:  1,
		ValorTotal: 299.99,
		Status:     "pendente",
		DataPedido: time.Now(),
	}

	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedPedido, nil)

	result, err := svc.FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPedido.ID, result.ID)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByID_NotFound(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Pedido)(nil), assert.AnError)

	result, err := svc.FindByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByClienteID_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	expectedPedidos := []model.Pedido{
		{ID: 1, ClienteID: 1, ValorTotal: 299.99, Status: "pendente"},
		{ID: 2, ClienteID: 1, ValorTotal: 199.99, Status: "pago"},
	}

	mockPedidoRepo.On("FindByClienteID", mock.Anything, uint(1)).Return(expectedPedidos, nil)

	result, err := svc.FindByClienteID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByStatus_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	expectedPedidos := []model.Pedido{
		{ID: 1, ClienteID: 1, ValorTotal: 299.99, Status: "pendente"},
		{ID: 2, ClienteID: 2, ValorTotal: 199.99, Status: "pendente"},
	}

	mockPedidoRepo.On("FindByStatus", mock.Anything, "pendente").Return(expectedPedidos, nil)

	result, err := svc.FindByStatus(context.Background(), "pendente")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_UpdateStatus_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	existingPedido := &model.Pedido{
		ID:         1,
		ClienteID:  1,
		ValorTotal: 299.99,
		Status:     "pendente",
	}

	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(existingPedido, nil)
	mockPedidoRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pedido")).Return(nil)

	req := &dto.UpdatePedidoRequest{
		Status: "pago",
	}

	result, err := svc.UpdateStatus(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "pago", result.Status)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_UpdateStatus_NotFound(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Pedido)(nil), assert.AnError)

	req := &dto.UpdatePedidoRequest{
		Status: "pago",
	}

	result, err := svc.UpdateStatus(context.Background(), 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_UpdateStatus_ValidationError(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	req := &dto.UpdatePedidoRequest{
		Status: "invalid_status", // Invalid status
	}

	result, err := svc.UpdateStatus(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Status")
}

func TestPedidoService_Delete_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	// Mock FindByID to verify pedido exists
	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Pedido{ID: 1}, nil)
	mockPedidoRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_Count_Success(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("Count", mock.Anything).Return(int64(50), nil)

	count, err := svc.Count(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(50), count)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindAll_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindAll", mock.Anything).Return([]model.Pedido(nil), assert.AnError)

	result, err := svc.FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByClienteID_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindByClienteID", mock.Anything, uint(1)).Return([]model.Pedido(nil), assert.AnError)

	result, err := svc.FindByClienteID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_FindByStatus_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindByStatus", mock.Anything, "pendente").Return([]model.Pedido(nil), assert.AnError)

	result, err := svc.FindByStatus(context.Background(), "pendente")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_UpdateStatus_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	existingPedido := &model.Pedido{
		ID:     1,
		Status: "pendente",
	}

	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(existingPedido, nil)
	mockPedidoRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pedido")).Return(assert.AnError)

	req := &dto.UpdatePedidoRequest{
		Status: "pago",
	}

	result, err := svc.UpdateStatus(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_Delete_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Pedido{ID: 1}, nil)
	mockPedidoRepo.On("Delete", mock.Anything, uint(1)).Return(assert.AnError)

	err := svc.Delete(context.Background(), 1)

	assert.Error(t, err)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_Count_Error(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockPedidoRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)

	count, err := svc.Count(context.Background())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	mockPedidoRepo.AssertExpectations(t)
}

func TestPedidoService_Create_EstoqueInsuficiente(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockClienteRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Cliente{
		ID: 1,
	}, nil)

	// Produto com estoque insuficiente
	mockProdutoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Produto{
		ID:      1,
		Preco:   100.00,
		Estoque: 5, // Estoque insuficiente
		Ativo:   true,
	}, nil)

	req := &dto.CreatePedidoRequest{
		ClienteID: 1,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: 1, Quantidade: 10}, // Quantidade maior que estoque
		},
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "insuficiente")
	mockClienteRepo.AssertExpectations(t)
	mockProdutoRepo.AssertExpectations(t)
}

func TestPedidoService_Create_ProdutoInativo(t *testing.T) {
	mockPedidoRepo := new(MockPedidoRepository)
	mockClienteRepo := new(MockClienteRepository)
	mockProdutoRepo := new(MockProdutoRepository)
	svc := service.NewPedidoService(mockPedidoRepo, mockClienteRepo, mockProdutoRepo)

	mockClienteRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Cliente{
		ID: 1,
	}, nil)

	// Produto inativo
	mockProdutoRepo.On("FindByID", mock.Anything, uint(1)).Return(&model.Produto{
		ID:      1,
		Preco:   100.00,
		Estoque: 10,
		Ativo:   false, // Produto inativo
	}, nil)

	req := &dto.CreatePedidoRequest{
		ClienteID: 1,
		Itens: []dto.CreateItemPedidoRequest{
			{ProdutoID: 1, Quantidade: 1},
		},
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "inativo")
	mockClienteRepo.AssertExpectations(t)
	mockProdutoRepo.AssertExpectations(t)
}
