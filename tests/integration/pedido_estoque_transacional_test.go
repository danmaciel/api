package integration

import (
 "context"
 "path/filepath"
 "strings"
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

func TestAtomicOrderRollbackWhenStockInsufficient(t *testing.T) {
 db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "rollback.db")), &gorm.Config{})
 require.NoError(t, err)
 require.NoError(t, db.AutoMigrate(&model.Cliente{}, &model.Produto{}, &model.Pedido{}, &model.PedidoProduto{}))
 customer := &model.Cliente{Nome:"Cliente Teste", Email:"rollback@example.com", CPF:"12345678902"}
 require.NoError(t, db.Create(customer).Error)
 first := &model.Produto{Nome:"Produto Primeiro",SKU:"ROLLBACK-1",Preco:15,Estoque:5,Ativo:true}
 second := &model.Produto{Nome:"Produto Segundo",SKU:"ROLLBACK-2",Preco:15,Estoque:1,Ativo:true}
 require.NoError(t,db.Create(first).Error)
 require.NoError(t,db.Create(second).Error)
 order := &model.Pedido{ClienteID:customer.ID,Status:"pendente",ValorTotal:75,Itens:[]model.PedidoProduto{
  {ProdutoID:first.ID,Quantidade:3,PrecoUnitario:15},
  {ProdutoID:second.ID,Quantidade:2,PrecoUnitario:15},
 }}
 err = repository.NewPedidoRepositorySQLite(db).Create(context.Background(), order)
 require.ErrorContains(t,err,"estoque indisponível")
 var updated model.Produto
 require.NoError(t,db.First(&updated,first.ID).Error)
 require.Equal(t,5,updated.Estoque)
 var count int64
 require.NoError(t,db.Model(&model.Pedido{}).Count(&count).Error)
 require.Zero(t,count)
}

func TestAtomicOrderDuplicateItemsCountTogether(t *testing.T) {
 db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "duplicates.db")), &gorm.Config{})
 require.NoError(t, err)
 require.NoError(t, db.AutoMigrate(&model.Cliente{}, &model.Produto{}, &model.Pedido{}, &model.PedidoProduto{}))
 customer := &model.Cliente{Nome:"Cliente Teste",Email:"duplicates@example.com",CPF:"12345678903"}
 require.NoError(t, db.Create(customer).Error)
 product := &model.Produto{Nome:"Produto Duplicado",SKU:"DUPLICATE-1",Preco:10,Estoque:5,Ativo:true}
 require.NoError(t,db.Create(product).Error)
 order := &model.Pedido{ClienteID:customer.ID,Status:"pendente",ValorTotal:60,Itens:[]model.PedidoProduto{
  {ProdutoID:product.ID,Quantidade:3,PrecoUnitario:10},
  {ProdutoID:product.ID,Quantidade:3,PrecoUnitario:10},
 }}
 err = repository.NewPedidoRepositorySQLite(db).Create(context.Background(),order)
 require.ErrorContains(t,err,"estoque indisponível")
 require.NoError(t,db.First(product,product.ID).Error)
 require.Equal(t,5,product.Estoque)
}

func TestConcurrentOrdersDoNotOversell(t *testing.T) {
 db,err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(),"concurrent.db")+"?_busy_timeout=5000"),&gorm.Config{})
 require.NoError(t,err)
 sqlDB,err := db.DB()
 require.NoError(t,err)
 sqlDB.SetMaxOpenConns(1)
 require.NoError(t,db.AutoMigrate(&model.Cliente{},&model.Produto{},&model.Pedido{},&model.PedidoProduto{}))
 customer:=&model.Cliente{Nome:"Cliente Concorrente",Email:"concurrent@example.com",CPF:"12345678904"}
 require.NoError(t,db.Create(customer).Error)
 product:=&model.Produto{Nome:"Produto Concorrente",SKU:"CONCURRENT-1",Preco:10,Estoque:5,Ativo:true}
 require.NoError(t,db.Create(product).Error)
 repo:=repository.NewPedidoRepositorySQLite(db)
 start:=make(chan struct{})
 results:=make(chan error,2)
 for i:=0;i<2;i++ {
  go func(){
   <-start
   order:=&model.Pedido{ClienteID:customer.ID,Status:"pendente",ValorTotal:40,Itens:[]model.PedidoProduto{{ProdutoID:product.ID,Quantidade:4,PrecoUnitario:10}}}
   results<-repo.Create(context.Background(),order)
  }()
 }
 close(start)
 successes,rejections:=0,0
 for i:=0;i<2;i++ {
  err:=<-results
  if err==nil {successes++} else if strings.Contains(err.Error(),"estoque indisponível") {rejections++} else {t.Fatalf("unexpected error: %v",err)}
 }
 require.Equal(t,1,successes)
 require.Equal(t,1,rejections)
 require.NoError(t,db.First(product,product.ID).Error)
 require.Equal(t,1,product.Estoque)
}
