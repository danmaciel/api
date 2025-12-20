package repository

import (
	"context"

	"github.com/danmaciel/api/internal/model"
)

// ProdutoRepository define a interface para operações de dados de Produto
type ProdutoRepository interface {
	Create(ctx context.Context, produto *model.Produto) error
	FindAll(ctx context.Context) ([]model.Produto, error)
	FindByID(ctx context.Context, id uint) (*model.Produto, error)
	FindByName(ctx context.Context, nome string) ([]model.Produto, error)
	FindBySKU(ctx context.Context, sku string) (*model.Produto, error)
	FindByCategoria(ctx context.Context, categoria string) ([]model.Produto, error)
	Update(ctx context.Context, produto *model.Produto) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
