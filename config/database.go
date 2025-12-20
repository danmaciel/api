package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/danmaciel/api/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// inicializa o banco de dados com GORM
func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	// Verifica se o diretório do banco de dados existe, se não, cria-o
	dbDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório do banco de dados: %w", err)
	}

	// abre a conexão com o banco de dados
	db, err := gorm.Open(sqlite.Open(cfg.FilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// execução de migrations automáticas
	if err := db.AutoMigrate(
		&model.Cliente{},
		&model.Produto{},
		&model.Pedido{},
		&model.PedidoProduto{},
	); err != nil {
		return nil, fmt.Errorf("falha ao executar a migration: %w", err)
	}

	log.Println("Banco de dados inicializado com sucesso")
	return db, nil
}
