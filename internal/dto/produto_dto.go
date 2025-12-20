package dto

import "time"

// CreateProdutoRequest representa a requisição para criar um produto
type CreateProdutoRequest struct {
	Nome      string  `json:"nome" validate:"required,min=3,max=200"`
	Descricao string  `json:"descricao" validate:"max=1000"`
	Preco     float64 `json:"preco" validate:"required,gt=0"`
	Estoque   int     `json:"estoque" validate:"gte=0"`
	SKU       string  `json:"sku" validate:"required,min=3,max=50"`
	Categoria string  `json:"categoria" validate:"max=100"`
	Ativo     *bool   `json:"ativo"` // pointer para permitir false explícito
}

// UpdateProdutoRequest representa a requisição para atualizar um produto
type UpdateProdutoRequest struct {
	Nome      string  `json:"nome" validate:"omitempty,min=3,max=200"`
	Descricao string  `json:"descricao" validate:"max=1000"`
	Preco     float64 `json:"preco" validate:"omitempty,gt=0"`
	Estoque   int     `json:"estoque" validate:"gte=0"`
	SKU       string  `json:"sku" validate:"omitempty,min=3,max=50"`
	Categoria string  `json:"categoria" validate:"max=100"`
	Ativo     *bool   `json:"ativo"`
}

// ProdutoResponse representa a resposta de um produto
type ProdutoResponse struct {
	ID        uint      `json:"id"`
	Nome      string    `json:"nome"`
	Descricao string    `json:"descricao"`
	Preco     float64   `json:"preco"`
	Estoque   int       `json:"estoque"`
	SKU       string    `json:"sku"`
	Categoria string    `json:"categoria"`
	Ativo     bool      `json:"ativo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
