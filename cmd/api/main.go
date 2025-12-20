package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danmaciel/api/config"
	"github.com/danmaciel/api/internal/controller"
	"github.com/danmaciel/api/internal/repository"
	"github.com/danmaciel/api/internal/service"

	_ "github.com/danmaciel/api/docs" // Import for Swagger docs
)

// @title Cliente API
// @version 1.0
// @description REST API for managing clientes (customers) with CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@exemplo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Servidor iniciado em %s", cfg.GetServerAddress())

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Falha ao inicializar o banco de dados: %v", err)
	}

	// Initialize layers (Dependency Injection)
	// Repositories
	clienteRepo := repository.NewClienteRepositorySQLite(db)
	produtoRepo := repository.NewProdutoRepositorySQLite(db)
	pedidoRepo := repository.NewPedidoRepositorySQLite(db)

	// Services
	clienteService := service.NewClienteService(clienteRepo)
	produtoService := service.NewProdutoService(produtoRepo)
	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo)

	// Controllers
	clienteController := controller.NewClienteController(clienteService)
	produtoController := controller.NewProdutoController(produtoService)
	pedidoController := controller.NewPedidoController(pedidoService)

	// Setup router
	router := controller.SetupRouter(clienteController, produtoController, pedidoController)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Servidor iniciado em %s", cfg.GetServerAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Falha ao iniciar o servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Servidor sendo encerrado...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Servidor forçado a encerrar: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
