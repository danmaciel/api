package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/service"
	"github.com/go-chi/chi/v5"
)

type ProdutoController struct {
	service service.ProdutoService
}

// NewProdutoController creates a new controller instance
func NewProdutoController(service service.ProdutoService) *ProdutoController {
	return &ProdutoController{service: service}
}

// Create godoc
// @Summary Create a new produto
// @Description Create a new produto with the provided information
// @Tags produtos
// @Accept json
// @Produce json
// @Param produto body dto.CreateProdutoRequest true "Produto data"
// @Success 201 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos [post]
func (c *ProdutoController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.Create(r.Context(), &req)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao criar produto", err.Error())
		return
	}

	c.respondJSON(w, http.StatusCreated, response)
}

// FindAll godoc
// @Summary Get all produtos
// @Description Retrieve all produtos from the database
// @Tags produtos
// @Produce json
// @Success 200 {array} dto.ProdutoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos [get]
func (c *ProdutoController) FindAll(w http.ResponseWriter, r *http.Request) {
	responses, err := c.service.FindAll(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar produtos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// FindByID godoc
// @Summary Get produto by ID
// @Description Retrieve a specific produto by ID
// @Tags produtos
// @Produce json
// @Param id path int true "Produto ID"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/{id} [get]
func (c *ProdutoController) FindByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	response, err := c.service.FindByID(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "produto not found" {
			c.respondError(w, http.StatusNotFound, "Produto nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar produto", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// FindByName godoc
// @Summary Get produtos by name
// @Description Retrieve produtos matching the specified name (partial match)
// @Tags produtos
// @Produce json
// @Param name path string true "Produto name"
// @Success 200 {array} dto.ProdutoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/nome/{name} [get]
func (c *ProdutoController) FindByName(w http.ResponseWriter, r *http.Request) {
	nome := chi.URLParam(r, "name")

	responses, err := c.service.FindByName(r.Context(), nome)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar produtos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// FindByCategoria godoc
// @Summary Get produtos by categoria
// @Description Retrieve produtos by categoria
// @Tags produtos
// @Produce json
// @Param categoria path string true "Produto categoria"
// @Success 200 {array} dto.ProdutoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/categoria/{categoria} [get]
func (c *ProdutoController) FindByCategoria(w http.ResponseWriter, r *http.Request) {
	categoria := chi.URLParam(r, "categoria")

	responses, err := c.service.FindByCategoria(r.Context(), categoria)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar produtos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// Update godoc
// @Summary Update produto
// @Description Update an existing produto
// @Tags produtos
// @Accept json
// @Produce json
// @Param id path int true "Produto ID"
// @Param produto body dto.UpdateProdutoRequest true "Produto data"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/{id} [put]
func (c *ProdutoController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	var req dto.UpdateProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.Update(r.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "produto not found" {
			c.respondError(w, http.StatusNotFound, "Produto nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao atualizar produto", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete produto
// @Description Delete a produto by ID
// @Tags produtos
// @Param id path int true "Produto ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/{id} [delete]
func (c *ProdutoController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	if err := c.service.Delete(r.Context(), uint(id)); err != nil {
		if err.Error() == "produto not found" {
			c.respondError(w, http.StatusNotFound, "Produto nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao deletar produto", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Count godoc
// @Summary Count produtos
// @Description Get the total number of produtos
// @Tags produtos
// @Produce json
// @Success 200 {object} dto.CountResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /produtos/count [get]
func (c *ProdutoController) Count(w http.ResponseWriter, r *http.Request) {
	count, err := c.service.Count(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao contar produtos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, dto.CountResponse{Count: count})
}

// Helper methods for JSON responses
func (c *ProdutoController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c *ProdutoController) respondError(w http.ResponseWriter, status int, error string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   error,
		Message: message,
	})
}
