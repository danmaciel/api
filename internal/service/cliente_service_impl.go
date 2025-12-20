package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type clienteServiceImpl struct {
	repo     repository.ClienteRepository
	validate *validator.Validate
}

// NewPedidoService cria uma nova instância do serviço
func NewClienteService(repo repository.ClienteRepository) ClienteService {
	return &clienteServiceImpl{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *clienteServiceImpl) Create(ctx context.Context, req *dto.CreateClienteRequest) (*dto.ClienteResponse, error) {
	// Validar request
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Dto para model
	cliente := &model.Cliente{
		Nome:     req.Nome,
		Email:    req.Email,
		CPF:      req.CPF,
		Telefone: req.Telefone,
	}

	// cria no banco
	if err := s.repo.Create(ctx, cliente); err != nil {
		return nil, fmt.Errorf("falha ao criar cliente: %w", err)
	}

	// model para dto
	return s.toResponse(cliente), nil
}

func (s *clienteServiceImpl) FindAll(ctx context.Context) ([]dto.ClienteResponse, error) {
	clientes, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar os clientes: %w", err)
	}

	responses := make([]dto.ClienteResponse, len(clientes))
	for i, cliente := range clientes {
		responses[i] = *s.toResponse(&cliente)
	}
	return responses, nil
}

func (s *clienteServiceImpl) FindByID(ctx context.Context, id uint) (*dto.ClienteResponse, error) {
	cliente, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar cliente por ID: %w", err)
	}
	if cliente == nil {
		return nil, errors.New("cliente not found")
	}
	return s.toResponse(cliente), nil
}

func (s *clienteServiceImpl) FindByName(ctx context.Context, nome string) ([]dto.ClienteResponse, error) {
	clientes, err := s.repo.FindByName(ctx, nome)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar clientes por nome: %w", err)
	}

	responses := make([]dto.ClienteResponse, len(clientes))
	for i, cliente := range clientes {
		responses[i] = *s.toResponse(&cliente)
	}
	return responses, nil
}

func (s *clienteServiceImpl) Update(ctx context.Context, id uint, req *dto.UpdateClienteRequest) (*dto.ClienteResponse, error) {
	// valida request
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("erro de validação: %w", err)
	}

	// procurar cliente existente
	cliente, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar cliente: %w", err)
	}
	if cliente == nil {
		return nil, errors.New("cliente not found")
	}

	// atualizar campos
	if req.Nome != "" {
		cliente.Nome = req.Nome
	}
	if req.Email != "" {
		cliente.Email = req.Email
	}
	if req.CPF != "" {
		cliente.CPF = req.CPF
	}
	if req.Telefone != "" {
		cliente.Telefone = req.Telefone
	}

	// atualizar no banco
	if err := s.repo.Update(ctx, cliente); err != nil {
		return nil, fmt.Errorf("falha ao atualizar cliente: %w", err)
	}

	return s.toResponse(cliente), nil
}

func (s *clienteServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cliente not found")
		}
		return fmt.Errorf("falha ao deletar cliente: %w", err)
	}
	return nil
}

func (s *clienteServiceImpl) Count(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("falha ao contar clientes: %w", err)
	}
	return count, nil
}

// model para response dto
func (s *clienteServiceImpl) toResponse(cliente *model.Cliente) *dto.ClienteResponse {
	return &dto.ClienteResponse{
		ID:        cliente.ID,
		Nome:      cliente.Nome,
		Email:     cliente.Email,
		CPF:       cliente.CPF,
		Telefone:  cliente.Telefone,
		CreatedAt: cliente.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cliente.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
