package service

import (
	"context"
	"errors"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/repository"
	"github.com/go-playground/validator/v10"
)

type produtoServiceImpl struct {
	repo     repository.ProdutoRepository
	validate *validator.Validate
}

// NewProdutoService cria uma nova instância do serviço
func NewProdutoService(repo repository.ProdutoRepository) ProdutoService {
	return &produtoServiceImpl{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *produtoServiceImpl) Create(ctx context.Context, req *dto.CreateProdutoRequest) (*dto.ProdutoResponse, error) {
	// Validar request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Verificar se SKU já existe
	existente, err := s.repo.FindBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, errors.New("SKU já cadastrado")
	}

	// Mapear DTO para Model
	ativo := true
	if req.Ativo != nil {
		ativo = *req.Ativo
	}

	produto := &model.Produto{
		Nome:      req.Nome,
		Descricao: req.Descricao,
		Preco:     req.Preco,
		Estoque:   req.Estoque,
		SKU:       req.SKU,
		Categoria: req.Categoria,
		Ativo:     ativo,
	}

	// Criar no banco
	if err := s.repo.Create(ctx, produto); err != nil {
		return nil, err
	}

	// Mapear Model para Response
	return s.toResponse(produto), nil
}

func (s *produtoServiceImpl) FindAll(ctx context.Context) ([]dto.ProdutoResponse, error) {
	produtos, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProdutoResponse, len(produtos))
	for i, produto := range produtos {
		responses[i] = *s.toResponse(&produto)
	}

	return responses, nil
}

func (s *produtoServiceImpl) FindByID(ctx context.Context, id uint) (*dto.ProdutoResponse, error) {
	produto, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(produto), nil
}

func (s *produtoServiceImpl) FindByName(ctx context.Context, nome string) ([]dto.ProdutoResponse, error) {
	produtos, err := s.repo.FindByName(ctx, nome)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProdutoResponse, len(produtos))
	for i, produto := range produtos {
		responses[i] = *s.toResponse(&produto)
	}

	return responses, nil
}

func (s *produtoServiceImpl) FindByCategoria(ctx context.Context, categoria string) ([]dto.ProdutoResponse, error) {
	produtos, err := s.repo.FindByCategoria(ctx, categoria)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProdutoResponse, len(produtos))
	for i, produto := range produtos {
		responses[i] = *s.toResponse(&produto)
	}

	return responses, nil
}

func (s *produtoServiceImpl) Update(ctx context.Context, id uint, req *dto.UpdateProdutoRequest) (*dto.ProdutoResponse, error) {
	// Validar request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Buscar produto existente
	produto, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar campos se fornecidos
	if req.Nome != "" {
		produto.Nome = req.Nome
	}
	if req.Descricao != "" {
		produto.Descricao = req.Descricao
	}
	if req.Preco > 0 {
		produto.Preco = req.Preco
	}
	if req.Estoque >= 0 {
		produto.Estoque = req.Estoque
	}
	if req.SKU != "" {
		// Verificar se novo SKU já existe em outro produto
		existente, err := s.repo.FindBySKU(ctx, req.SKU)
		if err != nil {
			return nil, err
		}
		if existente != nil && existente.ID != id {
			return nil, errors.New("SKU já cadastrado em outro produto")
		}
		produto.SKU = req.SKU
	}
	if req.Categoria != "" {
		produto.Categoria = req.Categoria
	}
	if req.Ativo != nil {
		produto.Ativo = *req.Ativo
	}

	// Atualizar no banco
	if err := s.repo.Update(ctx, produto); err != nil {
		return nil, err
	}

	return s.toResponse(produto), nil
}

func (s *produtoServiceImpl) Delete(ctx context.Context, id uint) error {
	// Verificar se existe
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *produtoServiceImpl) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// toResponse converte Model para Response DTO
func (s *produtoServiceImpl) toResponse(produto *model.Produto) *dto.ProdutoResponse {
	return &dto.ProdutoResponse{
		ID:        produto.ID,
		Nome:      produto.Nome,
		Descricao: produto.Descricao,
		Preco:     produto.Preco,
		Estoque:   produto.Estoque,
		SKU:       produto.SKU,
		Categoria: produto.Categoria,
		Ativo:     produto.Ativo,
		CreatedAt: produto.CreatedAt,
		UpdatedAt: produto.UpdatedAt,
	}
}
