package repository

import (
	"context"

	"github.com/danmaciel/api/internal/model"
)

// PedidoRepository define a interface para operações de dados de Pedido
type PedidoRepository interface {
	Create(ctx context.Context, pedido *model.Pedido) error
	FindAll(ctx context.Context) ([]model.Pedido, error)
	FindByID(ctx context.Context, id uint) (*model.Pedido, error)
	FindByClienteID(ctx context.Context, clienteID uint) ([]model.Pedido, error)
	FindByStatus(ctx context.Context, status string) ([]model.Pedido, error)
	Update(ctx context.Context, pedido *model.Pedido) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
