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

// MockClienteRepository is a mock implementation of ClienteRepository
type MockClienteRepository struct {
	mock.Mock
}

func (m *MockClienteRepository) Create(ctx context.Context, cliente *model.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}

func (m *MockClienteRepository) FindAll(ctx context.Context) ([]model.Cliente, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Cliente), args.Error(1)
}

func (m *MockClienteRepository) FindByID(ctx context.Context, id uint) (*model.Cliente, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cliente), args.Error(1)
}

func (m *MockClienteRepository) FindByName(ctx context.Context, nome string) ([]model.Cliente, error) {
	args := m.Called(ctx, nome)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Cliente), args.Error(1)
}

func (m *MockClienteRepository) Update(ctx context.Context, cliente *model.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}

func (m *MockClienteRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClienteRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Test cases
func TestClienteService_Create_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	req := &dto.CreateClienteRequest{
		Nome:     "João Silva",
		Email:    "joao@example.com",
		CPF:      "12345678901",
		Telefone: "11999999999",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Cliente")).Return(nil)

	result, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Nome, result.Nome)
	assert.Equal(t, req.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Create_ValidationError(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	req := &dto.CreateClienteRequest{
		Nome:  "", // Invalid: empty name
		Email: "joao@example.com",
		CPF:   "12345678901",
	}

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "validation error")
}

func TestClienteService_FindByID_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	expectedCliente := &model.Cliente{
		ID:       1,
		Nome:     "João Silva",
		Email:    "joao@example.com",
		CPF:      "12345678901",
		Telefone: "11999999999",
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedCliente, nil)

	result, err := svc.FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCliente.Nome, result.Nome)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_FindByID_NotFound(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Cliente)(nil), assert.AnError)

	result, err := svc.FindByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Count_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("Count", mock.Anything).Return(int64(10), nil)

	count, err := svc.Count(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_FindAll_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	expectedClientes := []model.Cliente{
		{ID: 1, Nome: "Cliente 1", Email: "c1@example.com"},
		{ID: 2, Nome: "Cliente 2", Email: "c2@example.com"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expectedClientes, nil)

	result, err := svc.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_FindAll_Error(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("FindAll", mock.Anything).Return([]model.Cliente(nil), assert.AnError)

	result, err := svc.FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_FindByName_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	expectedClientes := []model.Cliente{
		{ID: 1, Nome: "João Silva", Email: "joao@example.com"},
	}

	mockRepo.On("FindByName", mock.Anything, "João").Return(expectedClientes, nil)

	result, err := svc.FindByName(context.Background(), "João")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Update_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	existingCliente := &model.Cliente{
		ID:    1,
		Nome:  "João Silva",
		Email: "joao@example.com",
		CPF:   "12345678901",
	}

	mockRepo.On("FindByID", mock.Anything, uint(1)).Return(existingCliente, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Cliente")).Return(nil)

	req := &dto.UpdateClienteRequest{
		Nome: "João Silva Updated",
	}

	result, err := svc.Update(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João Silva Updated", result.Nome)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, uint(999)).Return((*model.Cliente)(nil), assert.AnError)

	req := &dto.UpdateClienteRequest{
		Nome: "Updated Name",
	}

	result, err := svc.Update(context.Background(), 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Delete_Success(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Delete_Error(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("Delete", mock.Anything, uint(999)).Return(assert.AnError)

	err := svc.Delete(context.Background(), 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Count_Error(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	mockRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)

	count, err := svc.Count(context.Background())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	mockRepo.AssertExpectations(t)
}

func TestClienteService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(MockClienteRepository)
	svc := service.NewClienteService(mockRepo)

	req := &dto.CreateClienteRequest{
		Nome:     "João Silva",
		Email:    "joao@example.com",
		CPF:      "12345678901",
		Telefone: "11999999999",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Cliente")).Return(assert.AnError)

	result, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
