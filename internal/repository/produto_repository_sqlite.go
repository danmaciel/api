package repository

import (
	"context"
	"errors"

	"github.com/danmaciel/api/internal/model"
	"gorm.io/gorm"
)

type produtoRepositorySQLite struct {
	db *gorm.DB
}

// NewProdutoRepositorySQLite cria uma nova instância do repositório SQLite
func NewProdutoRepositorySQLite(db *gorm.DB) ProdutoRepository {
	return &produtoRepositorySQLite{db: db}
}

func (r *produtoRepositorySQLite) Create(ctx context.Context, produto *model.Produto) error {
	return r.db.WithContext(ctx).Create(produto).Error
}

func (r *produtoRepositorySQLite) FindAll(ctx context.Context) ([]model.Produto, error) {
	var produtos []model.Produto
	err := r.db.WithContext(ctx).Find(&produtos).Error
	return produtos, err
}

func (r *produtoRepositorySQLite) FindByID(ctx context.Context, id uint) (*model.Produto, error) {
	var produto model.Produto
	err := r.db.WithContext(ctx).First(&produto, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto not found")
		}
		return nil, err
	}
	return &produto, nil
}

func (r *produtoRepositorySQLite) FindByName(ctx context.Context, nome string) ([]model.Produto, error) {
	var produtos []model.Produto
	err := r.db.WithContext(ctx).Where("nome LIKE ?", "%"+nome+"%").Find(&produtos).Error
	return produtos, err
}

func (r *produtoRepositorySQLite) FindBySKU(ctx context.Context, sku string) (*model.Produto, error) {
	var produto model.Produto
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&produto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // SKU não encontrado não é erro
		}
		return nil, err
	}
	return &produto, nil
}

func (r *produtoRepositorySQLite) FindByCategoria(ctx context.Context, categoria string) ([]model.Produto, error) {
	var produtos []model.Produto
	err := r.db.WithContext(ctx).Where("categoria = ?", categoria).Find(&produtos).Error
	return produtos, err
}

func (r *produtoRepositorySQLite) Update(ctx context.Context, produto *model.Produto) error {
	result := r.db.WithContext(ctx).Save(produto)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("produto not found")
	}
	return nil
}

func (r *produtoRepositorySQLite) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Produto{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("produto not found")
	}
	return nil
}

func (r *produtoRepositorySQLite) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Produto{}).Count(&count).Error
	return count, err
}
