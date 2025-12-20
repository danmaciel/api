package dto

import "time"

// CreatePedidoRequest representa a requisição para criar um pedido
type CreatePedidoRequest struct {
	ClienteID uint                   `json:"cliente_id" validate:"required"`
	Itens     []CreateItemPedidoRequest `json:"itens" validate:"required,min=1,dive"`
	Status    string                 `json:"status" validate:"omitempty,oneof=pendente pago enviado entregue cancelado"`
}

// CreateItemPedidoRequest representa um item no pedido
type CreateItemPedidoRequest struct {
	ProdutoID  uint `json:"produto_id" validate:"required"`
	Quantidade int  `json:"quantidade" validate:"required,gt=0"`
}

// UpdatePedidoRequest representa a requisição para atualizar um pedido
type UpdatePedidoRequest struct {
	Status string `json:"status" validate:"required,oneof=pendente pago enviado entregue cancelado"`
}

// PedidoResponse representa a resposta de um pedido
type PedidoResponse struct {
	ID          uint                 `json:"id"`
	ClienteID   uint                 `json:"cliente_id"`
	Cliente     *ClienteResponse     `json:"cliente,omitempty"`
	Itens       []ItemPedidoResponse `json:"itens"`
	ValorTotal  float64              `json:"valor_total"`
	Status      string               `json:"status"`
	DataPedido  time.Time            `json:"data_pedido"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// ItemPedidoResponse representa um item do pedido na resposta
type ItemPedidoResponse struct {
	ID            uint             `json:"id"`
	ProdutoID     uint             `json:"produto_id"`
	Produto       *ProdutoResponse `json:"produto,omitempty"`
	Quantidade    int              `json:"quantidade"`
	PrecoUnitario float64          `json:"preco_unitario"`
	Subtotal      float64          `json:"subtotal"`
}
