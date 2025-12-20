package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/service"
	"github.com/go-chi/chi/v5"
)

type ClienteController struct {
	service service.ClienteService
}

// NewClienteController creates a new controller instance
func NewClienteController(service service.ClienteService) *ClienteController {
	return &ClienteController{service: service}
}

// Create godoc
// @Summary Create a new cliente
// @Description Create a new cliente with the provided information
// @Tags clientes
// @Accept json
// @Produce json
// @Param cliente body dto.CreateClienteRequest true "Cliente data"
// @Success 201 {object} dto.ClienteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes [post]
func (c *ClienteController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.Create(r.Context(), &req)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao criar cliente", err.Error())
		return
	}

	c.respondJSON(w, http.StatusCreated, response)
}

// FindAll godoc
// @Summary Get all clientes
// @Description Retrieve all clientes from the database
// @Tags clientes
// @Produce json
// @Success 200 {array} dto.ClienteResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes [get]
func (c *ClienteController) FindAll(w http.ResponseWriter, r *http.Request) {
	responses, err := c.service.FindAll(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar clientes", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// FindByID godoc
// @Summary Get cliente by ID
// @Description Retrieve a specific cliente by ID
// @Tags clientes
// @Produce json
// @Param id path int true "Cliente ID"
// @Success 200 {object} dto.ClienteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/{id} [get]
func (c *ClienteController) FindByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	response, err := c.service.FindByID(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "cliente not found" {
			c.respondError(w, http.StatusNotFound, "Cliente nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar cliente", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// FindByNome godoc
// @Summary Get clientes by name
// @Description Retrieve clientes matching the specified name (partial match)
// @Tags clientes
// @Produce json
// @Param name path string true "Cliente name"
// @Success 200 {array} dto.ClienteResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/nome/{name} [get]
func (c *ClienteController) FindByName(w http.ResponseWriter, r *http.Request) {
	nome := chi.URLParam(r, "name")

	responses, err := c.service.FindByName(r.Context(), nome)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao recuperar clientes", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, responses)
}

// Update godoc
// @Summary Update cliente
// @Description Update an existing cliente
// @Tags clientes
// @Accept json
// @Produce json
// @Param id path int true "Cliente ID"
// @Param cliente body dto.UpdateClienteRequest true "Cliente data"
// @Success 200 {object} dto.ClienteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/{id} [put]
func (c *ClienteController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	var req dto.UpdateClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, "Corpo da requisição inválido", err.Error())
		return
	}

	response, err := c.service.Update(r.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "cliente not found" {
			c.respondError(w, http.StatusNotFound, "Cliente nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao atualizar cliente", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, response)
}

// Delete godoc
// @Summary Delete cliente
// @Description Delete a cliente by ID
// @Tags clientes
// @Param id path int true "Cliente ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/{id} [delete]
func (c *ClienteController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "Id Parametro Invalido", err.Error())
		return
	}

	if err := c.service.Delete(r.Context(), uint(id)); err != nil {
		if err.Error() == "cliente not found" {
			c.respondError(w, http.StatusNotFound, "Cliente nao encontrado", "")
			return
		}
		c.respondError(w, http.StatusInternalServerError, "Falha ao deletar cliente", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Count godoc
// @Summary Count clientes
// @Description Get the total number of clientes
// @Tags clientes
// @Produce json
// @Success 200 {object} dto.CountResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clientes/count [get]
func (c *ClienteController) Count(w http.ResponseWriter, r *http.Request) {
	count, err := c.service.Count(r.Context())
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, "Falha ao contar clientes", err.Error())
		return
	}

	c.respondJSON(w, http.StatusOK, dto.CountResponse{Count: count})
}

// Helper methods for JSON responses
func (c *ClienteController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (c *ClienteController) respondError(w http.ResponseWriter, status int, error string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   error,
		Message: message,
	})
}
