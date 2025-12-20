package repository

import (
	"context"

	"github.com/danmaciel/api/internal/model"
)

// ClienteRepository defines the interface for cliente data access
type ClienteRepository interface {
	Create(ctx context.Context, cliente *model.Cliente) error
	FindAll(ctx context.Context) ([]model.Cliente, error)
	FindByID(ctx context.Context, id uint) (*model.Cliente, error)
	FindByName(ctx context.Context, nome string) ([]model.Cliente, error)
	Update(ctx context.Context, cliente *model.Cliente) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
