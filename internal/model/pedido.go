package model

import (
	"time"

	"gorm.io/gorm"
)

// Pedido representa a entidade de domínio Pedido
type Pedido struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	ClienteID   uint            `gorm:"not null" json:"cliente_id" validate:"required"`
	Cliente     Cliente         `gorm:"foreignKey:ClienteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"cliente,omitempty"`
	Itens       []PedidoProduto `gorm:"foreignKey:PedidoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"itens,omitempty"`
	ValorTotal  float64         `gorm:"type:decimal(10,2);not null;default:0" json:"valor_total"`
	Status      string          `gorm:"type:varchar(20);not null;default:'pendente'" json:"status" validate:"required,oneof=pendente pago enviado entregue cancelado"`
	DataPedido  time.Time       `gorm:"not null" json:"data_pedido"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// PedidoProduto representa a tabela de junção entre Pedido e Produto (itens do pedido)
type PedidoProduto struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PedidoID      uint           `gorm:"not null" json:"pedido_id"`
	Pedido        Pedido         `gorm:"foreignKey:PedidoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ProdutoID     uint           `gorm:"not null" json:"produto_id"`
	Produto       Produto        `gorm:"foreignKey:ProdutoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"produto,omitempty"`
	Quantidade    int            `gorm:"not null" json:"quantidade" validate:"required,gt=0"`
	PrecoUnitario float64        `gorm:"type:decimal(10,2);not null" json:"preco_unitario" validate:"required,gt=0"`
	Subtotal      float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica o nome da tabela para o GORM
func (Pedido) TableName() string {
	return "pedidos"
}

// TableName especifica o nome da tabela para o GORM
func (PedidoProduto) TableName() string {
	return "pedido_produtos"
}

// BeforeSave hook para calcular o subtotal antes de salvar
func (pp *PedidoProduto) BeforeSave(tx *gorm.DB) error {
	pp.Subtotal = float64(pp.Quantidade) * pp.PrecoUnitario
	return nil
}
