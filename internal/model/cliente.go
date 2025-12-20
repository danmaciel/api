package model

import (
	"time"

	"gorm.io/gorm"
)

// Cliente represents a customer entity
type Cliente struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Nome      string         `gorm:"type:varchar(100);not null" json:"nome" validate:"required,min=3,max=100"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" validate:"required,email"`
	CPF       string         `gorm:"type:varchar(11);uniqueIndex;not null" json:"cpf" validate:"required,len=11,numeric"`
	Telefone  string         `gorm:"type:varchar(15)" json:"telefone" validate:"omitempty,min=10,max=15"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Cliente
func (Cliente) TableName() string {
	return "clientes"
}
