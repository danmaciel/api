package repository

import (
	"context"
	"errors"

	"github.com/danmaciel/api/internal/model"
	"gorm.io/gorm"
)

type pedidoRepositorySQLite struct {
	db *gorm.DB
}

// NewPedidoRepositorySQLite cria uma nova instância do repositório SQLite
func NewPedidoRepositorySQLite(db *gorm.DB) PedidoRepository {
	return &pedidoRepositorySQLite{db: db}
}

func (r *pedidoRepositorySQLite) Create(ctx context.Context, pedido *model.Pedido) error {
	return r.db.WithContext(ctx).Create(pedido).Error
}

func (r *pedidoRepositorySQLite) FindAll(ctx context.Context) ([]model.Pedido, error) {
	var pedidos []model.Pedido
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Itens").
		Preload("Itens.Produto").
		Find(&pedidos).Error
	return pedidos, err
}

func (r *pedidoRepositorySQLite) FindByID(ctx context.Context, id uint) (*model.Pedido, error) {
	var pedido model.Pedido
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Itens").
		Preload("Itens.Produto").
		First(&pedido, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pedido not found")
		}
		return nil, err
	}
	return &pedido, nil
}

func (r *pedidoRepositorySQLite) FindByClienteID(ctx context.Context, clienteID uint) ([]model.Pedido, error) {
	var pedidos []model.Pedido
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Itens").
		Preload("Itens.Produto").
		Where("cliente_id = ?", clienteID).
		Find(&pedidos).Error
	return pedidos, err
}

func (r *pedidoRepositorySQLite) FindByStatus(ctx context.Context, status string) ([]model.Pedido, error) {
	var pedidos []model.Pedido
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Itens").
		Preload("Itens.Produto").
		Where("status = ?", status).
		Find(&pedidos).Error
	return pedidos, err
}

func (r *pedidoRepositorySQLite) Update(ctx context.Context, pedido *model.Pedido) error {
	result := r.db.WithContext(ctx).Save(pedido)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("pedido not found")
	}
	return nil
}

func (r *pedidoRepositorySQLite) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Pedido{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("pedido not found")
	}
	return nil
}

func (r *pedidoRepositorySQLite) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Pedido{}).Count(&count).Error
	return count, err
}
