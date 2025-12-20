package repository

import (
	"context"
	"errors"

	"github.com/danmaciel/api/internal/model"
	"gorm.io/gorm"
)

type clienteRepositorySQLite struct {
	db *gorm.DB
}

// NewClienteRepositorySQLite creates a new SQLite implementation of ClienteRepository
func NewClienteRepositorySQLite(db *gorm.DB) ClienteRepository {
	return &clienteRepositorySQLite{db: db}
}

func (r *clienteRepositorySQLite) Create(ctx context.Context, cliente *model.Cliente) error {
	result := r.db.WithContext(ctx).Create(cliente)
	return result.Error
}

func (r *clienteRepositorySQLite) FindAll(ctx context.Context) ([]model.Cliente, error) {
	var clientes []model.Cliente
	result := r.db.WithContext(ctx).Find(&clientes)
	if result.Error != nil {
		return nil, result.Error
	}
	return clientes, nil
}

func (r *clienteRepositorySQLite) FindByID(ctx context.Context, id uint) (*model.Cliente, error) {
	var cliente model.Cliente
	result := r.db.WithContext(ctx).First(&cliente, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cliente, nil
}

func (r *clienteRepositorySQLite) FindByName(ctx context.Context, nome string) ([]model.Cliente, error) {
	var clientes []model.Cliente
	result := r.db.WithContext(ctx).Where("nome LIKE ?", "%"+nome+"%").Find(&clientes)
	if result.Error != nil {
		return nil, result.Error
	}
	return clientes, nil
}

func (r *clienteRepositorySQLite) Update(ctx context.Context, cliente *model.Cliente) error {
	result := r.db.WithContext(ctx).Save(cliente)
	return result.Error
}

func (r *clienteRepositorySQLite) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Cliente{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *clienteRepositorySQLite) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Cliente{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
