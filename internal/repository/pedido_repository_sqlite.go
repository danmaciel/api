package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"

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

// Create reserves stock and persists the order in one transaction.
// A conditional UPDATE prevents orders from consuming the same stock concurrently.
func (r *pedidoRepositorySQLite) Create(ctx context.Context, pedido *model.Pedido) error {
	if len(pedido.Itens) == 0 {
		return errors.New("pedido sem itens")
	}

	quantidades := make(map[uint]int, len(pedido.Itens))
	for _, item := range pedido.Itens {
		if item.ProdutoID == 0 || item.Quantidade <= 0 {
			return errors.New("item de pedido inválido")
		}
		// Guard against integer overflow when aggregating the same product.
		if quantidades[item.ProdutoID] > int(^uint(0)>>1)-item.Quantidade {
			return errors.New("quantidade de produto excede limite")
		}
		quantidades[item.ProdutoID] += item.Quantidade
	}

	// Stable ordering reduces contention for engines with row-level locking.
	ids := make([]uint, 0, len(quantidades))
	for id := range quantidades {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			quantity := quantidades[id]
			result := tx.Model(&model.Produto{}).
				Where("id = ? AND ativo = ? AND estoque >= ?", id, true, quantity).
				UpdateColumn("estoque", gorm.Expr("estoque - ?", quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return fmt.Errorf("estoque indisponível para produto %d", id)
			}
		}
		return tx.Create(pedido).Error
	})
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
