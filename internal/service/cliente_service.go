package service

import (
	"context"

	"github.com/danmaciel/api/internal/dto"
)

// ClienteService defines the business logic interface for cliente operations
type ClienteService interface {
	Create(ctx context.Context, req *dto.CreateClienteRequest) (*dto.ClienteResponse, error)
	FindAll(ctx context.Context) ([]dto.ClienteResponse, error)
	FindByID(ctx context.Context, id uint) (*dto.ClienteResponse, error)
	FindByName(ctx context.Context, nome string) ([]dto.ClienteResponse, error)
	Update(ctx context.Context, id uint, req *dto.UpdateClienteRequest) (*dto.ClienteResponse, error)
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
