package service

import (
	"context"

	"github.com/danmaciel/api/internal/dto"
)

// ProdutoService define a interface para operações de negócio de Produto
type ProdutoService interface {
	Create(ctx context.Context, req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error)
	FindAll(ctx context.Context) ([]dto.ProdutoResponse, error)
	FindByID(ctx context.Context, id uint) (*dto.ProdutoResponse, error)
	FindByName(ctx context.Context, nome string) ([]dto.ProdutoResponse, error)
	FindByCategoria(ctx context.Context, categoria string) ([]dto.ProdutoResponse, error)
	Update(ctx context.Context, id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error)
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
