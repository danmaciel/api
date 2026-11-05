package integration

import (
 "context"
 "path/filepath"
 "testing"
 "github.com/danmaciel/api/internal/model"
 "github.com/danmaciel/api/internal/repository"
 "github.com/stretchr/testify/require"
 "gorm.io/driver/sqlite"
 "gorm.io/gorm"
)

func TestAtomicOrderDeductsStock(t *testing.T) {
 dsn := filepath.Join(t.TempDir(), "test.db")
 db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
 require.NoError(t, err)
 require.NoError(t, db.AutoMigrate(&model.Cliente{}, &model.Produto{}, &model.Pedido{}, &model.PedidoProduto{}))
 customer := &model.Cliente{Nome:"Cliente Teste", Email:"test@example.com", CPF:"12345678901"}
 require.NoError(t, db.Create(customer).Error)
 product := &model.Produto{Nome:"Produto Teste", SKU:"TEST-STOCK-1", Preco:20, Estoque:5, Ativo:true}
 require.NoError(t, db.Create(product).Error)
 order := &model.Pedido{ClienteID:customer.ID, Status:"pendente", ValorTotal:80, Itens:[]model.PedidoProduto{{ProdutoID:product.ID, Quantidade:4, PrecoUnitario:20}}}
 require.NoError(t, repository.NewPedidoRepositorySQLite(db).Create(context.Background(), order))
 require.NoError(t, db.First(product,product.ID).Error)
 require.Equal(t, 1, product.Estoque)
}
