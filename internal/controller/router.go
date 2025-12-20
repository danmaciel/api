package controller

import (
	"net/http"

	"github.com/danmaciel/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// configura o roteador com todas as rotas e middlewares
func SetupRouter(clienteController *ClienteController, produtoController *ProdutoController, pedidoController *PedidoController) *chi.Mux {
	r := chi.NewRouter()

	// aplicação de middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.ContentType("application/json"))

	// configuração de CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// documentação Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// rotas da API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Rotas de Clientes
		r.Route("/clientes", func(r chi.Router) {
			// IMPORTANT: More specific routes must come before generic ones
			r.Get("/count", clienteController.Count)           // Must be before /{id}
			r.Get("/nome/{name}", clienteController.FindByName) // Must be before /{id}

			r.Post("/", clienteController.Create)
			r.Get("/", clienteController.FindAll)
			r.Get("/{id}", clienteController.FindByID)
			r.Put("/{id}", clienteController.Update)
			r.Delete("/{id}", clienteController.Delete)
		})

		// Rotas de Produtos
		r.Route("/produtos", func(r chi.Router) {
			// IMPORTANT: More specific routes must come before generic ones
			r.Get("/count", produtoController.Count)                        // Must be before /{id}
			r.Get("/nome/{name}", produtoController.FindByName)             // Must be before /{id}
			r.Get("/categoria/{categoria}", produtoController.FindByCategoria) // Must be before /{id}

			r.Post("/", produtoController.Create)
			r.Get("/", produtoController.FindAll)
			r.Get("/{id}", produtoController.FindByID)
			r.Put("/{id}", produtoController.Update)
			r.Delete("/{id}", produtoController.Delete)
		})

		// Rotas de Pedidos
		r.Route("/pedidos", func(r chi.Router) {
			// IMPORTANT: More specific routes must come before generic ones
			r.Get("/count", pedidoController.Count)                          // Must be before /{id}
			r.Get("/cliente/{cliente_id}", pedidoController.FindByClienteID) // Must be before /{id}
			r.Get("/status/{status}", pedidoController.FindByStatus)         // Must be before /{id}

			r.Post("/", pedidoController.Create)
			r.Get("/", pedidoController.FindAll)
			r.Get("/{id}", pedidoController.FindByID)
			r.Put("/{id}", pedidoController.UpdateStatus)
			r.Delete("/{id}", pedidoController.Delete)
		})
	})

	// endpoint de Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
