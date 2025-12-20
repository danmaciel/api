package model

import (
	"time"

	"gorm.io/gorm"
)

// Produto representa a entidade de domínio Produto
type Produto struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Nome       string         `gorm:"type:varchar(200);not null" json:"nome" validate:"required,min=3,max=200"`
	Descricao  string         `gorm:"type:text" json:"descricao" validate:"max=1000"`
	Preco      float64        `gorm:"type:decimal(10,2);not null" json:"preco" validate:"required,gt=0"`
	Estoque    int            `gorm:"not null;default:0" json:"estoque" validate:"gte=0"`
	SKU        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku" validate:"required,min=3,max=50"`
	Categoria  string         `gorm:"type:varchar(100)" json:"categoria" validate:"max=100"`
	Ativo      bool           `gorm:"default:true" json:"ativo"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName especifica o nome da tabela para o GORM
func (Produto) TableName() string {
	return "produtos"
}
