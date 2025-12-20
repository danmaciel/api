package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/service"
	"github.com/go-chi/chi/v5"
)

type PedidoController struct {
	service service.PedidoService
}

// NewPedidoController creates a new controller instance
func NewPedidoController(service service.PedidoService) *PedidoController {
	return &PedidoController{service: service}
}

// Create godoc
// @Summary Create a new pedido
// @Description Create a new pedido with the provided information
// @Tags pedidos
// @Accept json
// @Produce json
// @Param pedido body dto.CreatePedidoRequest true "Pedido data"
// @Success 201 {object} dto.PedidoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos [post]
func (c *PedidoController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.Create(r.Context(), &req)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao criar pedido", err.Error())
		return
	}

	c.respondJSON(w, http.StatusCreated, response)
}

// FindAll godoc
// @Summary Get all pedidos
// @Description Retrieve all pedidos from the database
// @Tags pedidos
// @Produce json
// @Success 200 {array} dto.PedidoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos [get]
func (c *PedidoController) FindAll(w http.ResponseWriter, r *http.Request) {
	responses, err := c.service.FindAll(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar pedidos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// FindByID godoc
// @Summary Get pedido by ID
// @Description Retrieve a specific pedido by ID
// @Tags pedidos
// @Produce json
// @Param id path int true "Pedido ID"
// @Success 200 {object} dto.PedidoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/{id} [get]
func (c *PedidoController) FindByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	response, err := c.service.FindByID(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "pedido not found" {
			c.respondError(w, http.StatusNotFound, "Pedido nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar pedido", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// FindByClienteID godoc
// @Summary Get pedidos by cliente ID
// @Description Retrieve all pedidos for a specific cliente
// @Tags pedidos
// @Produce json
// @Param cliente_id path int true "Cliente ID"
// @Success 200 {array} dto.PedidoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/cliente/{cliente_id} [get]
func (c *PedidoController) FindByClienteID(w http.ResponseWriter, r *http.Request) {
	clienteIDStr := chi.URLParam(r, "cliente_id")
	clienteID, err := strconv.ParseUint(clienteIDStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Cliente ID Parametro Invalido", err.Error())
		return
	}

	responses, err := c.service.FindByClienteID(r.Context(), uint(clienteID))
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar pedidos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// FindByStatus godoc
// @Summary Get pedidos by status
// @Description Retrieve pedidos by status (pendente, pago, enviado, entregue, cancelado)
// @Tags pedidos
// @Produce json
// @Param status path string true "Pedido status"
// @Success 200 {array} dto.PedidoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/status/{status} [get]
func (c *PedidoController) FindByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")

	responses, err := c.service.FindByStatus(r.Context(), status)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar pedidos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// UpdateStatus godoc
// @Summary Update pedido status
// @Description Update the status of an existing pedido
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path int true "Pedido ID"
// @Param pedido body dto.UpdatePedidoRequest true "Pedido status update"
// @Success 200 {object} dto.PedidoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/{id} [put]
func (c *PedidoController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	var req dto.UpdatePedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.UpdateStatus(r.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "pedido not found" {
			c.respondError(w, http.StatusNotFound, "Pedido nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao atualizar pedido", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete pedido
// @Description Delete a pedido by ID
// @Tags pedidos
// @Param id path int true "Pedido ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/{id} [delete]
func (c *PedidoController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	if err := c.service.Delete(r.Context(), uint(id)); err != nil {
		if err.Error() == "pedido not found" {
			c.respondError(w, http.StatusNotFound, "Pedido nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao deletar pedido", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Count godoc
// @Summary Count pedidos
// @Description Get the total number of pedidos
// @Tags pedidos
// @Produce json
// @Success 200 {object} dto.CountResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /pedidos/count [get]
func (c *PedidoController) Count(w http.ResponseWriter, r *http.Request) {
	count, err := c.service.Count(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao contar pedidos", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, dto.CountResponse{Count: count})
}

// Helper methods for JSON responses
func (c *PedidoController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c *PedidoController) respondError(w http.ResponseWriter, status int, error string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   error,
		Message: message,
	})
}
