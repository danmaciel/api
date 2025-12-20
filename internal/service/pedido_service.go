package service

import (
	"context"

	"github.com/danmaciel/api/internal/dto"
)

// PedidoService define a interface para operações de negócio de Pedido
type PedidoService interface {
	Create(ctx context.Context, req *dto.CreatePedidoRequest) (*dto.PedidoResponse, error)
	FindAll(ctx context.Context) ([]dto.PedidoResponse, error)
	FindByID(ctx context.Context, id uint) (*dto.PedidoResponse, error)
	FindByClienteID(ctx context.Context, clienteID uint) ([]dto.PedidoResponse, error)
	FindByStatus(ctx context.Context, status string) ([]dto.PedidoResponse, error)
	UpdateStatus(ctx context.Context, id uint, req *dto.UpdatePedidoRequest) (*dto.PedidoResponse, error)
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
