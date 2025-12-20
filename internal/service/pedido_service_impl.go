package service

import (
	"context"
	"errors"
	"time"

	"github.com/danmaciel/api/internal/dto"
	"github.com/danmaciel/api/internal/model"
	"github.com/danmaciel/api/internal/repository"
	"github.com/go-playground/validator/v10"
)

type pedidoServiceImpl struct {
	pedidoRepo  repository.PedidoRepository
	clienteRepo repository.ClienteRepository
	produtoRepo repository.ProdutoRepository
	validate    *validator.Validate
}

// NewPedidoService cria uma nova instância do serviço
func NewPedidoService(pedidoRepo repository.PedidoRepository, clienteRepo repository.ClienteRepository, produtoRepo repository.ProdutoRepository) PedidoService {
	return &pedidoServiceImpl{
		pedidoRepo:  pedidoRepo,
		clienteRepo: clienteRepo,
		produtoRepo: produtoRepo,
		validate:    validator.New(),
	}
}

func (s *pedidoServiceImpl) Create(ctx context.Context, req *dto.CreatePedidoRequest) (*dto.PedidoResponse, error) {
	// Validar request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Validar se cliente existe
	cliente, err := s.clienteRepo.FindByID(ctx, req.ClienteID)
	if err != nil {
		return nil, errors.New("cliente não encontrado")
	}

	// Preparar itens do pedido
	var itens []model.PedidoProduto
	var valorTotal float64

	for _, itemReq := range req.Itens {
		// Buscar produto
		produto, err := s.produtoRepo.FindByID(ctx, itemReq.ProdutoID)
		if err != nil {
			return nil, errors.New("produto " + string(rune(itemReq.ProdutoID)) + " não encontrado")
		}

		// Verificar estoque
		if produto.Estoque < itemReq.Quantidade {
			return nil, errors.New("estoque insuficiente para produto: " + produto.Nome)
		}

		// Verificar se produto está ativo
		if !produto.Ativo {
			return nil, errors.New("produto inativo: " + produto.Nome)
		}

		// Criar item do pedido
		subtotal := float64(itemReq.Quantidade) * produto.Preco
		item := model.PedidoProduto{
			ProdutoID:     itemReq.ProdutoID,
			Quantidade:    itemReq.Quantidade,
			PrecoUnitario: produto.Preco,
			Subtotal:      subtotal,
		}

		itens = append(itens, item)
		valorTotal += subtotal
	}

	// Definir status padrão se não fornecido
	status := req.Status
	if status == "" {
		status = "pendente"
	}

	// Criar pedido
	pedido := &model.Pedido{
		ClienteID:  req.ClienteID,
		Itens:      itens,
		ValorTotal: valorTotal,
		Status:     status,
		DataPedido: time.Now(),
	}

	// Salvar no banco (com cascade para itens)
	if err := s.pedidoRepo.Create(ctx, pedido); err != nil {
		return nil, err
	}

	// Buscar pedido completo com relacionamentos
	pedidoCompleto, err := s.pedidoRepo.FindByID(ctx, pedido.ID)
	if err != nil {
		return nil, err
	}

	// Atribuir cliente ao pedido completo
	pedidoCompleto.Cliente = *cliente

	return s.toResponse(pedidoCompleto), nil
}

func (s *pedidoServiceImpl) FindAll(ctx context.Context) ([]dto.PedidoResponse, error) {
	pedidos, err := s.pedidoRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		responses[i] = *s.toResponse(&pedido)
	}

	return responses, nil
}

func (s *pedidoServiceImpl) FindByID(ctx context.Context, id uint) (*dto.PedidoResponse, error) {
	pedido, err := s.pedidoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(pedido), nil
}

func (s *pedidoServiceImpl) FindByClienteID(ctx context.Context, clienteID uint) ([]dto.PedidoResponse, error) {
	pedidos, err := s.pedidoRepo.FindByClienteID(ctx, clienteID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		responses[i] = *s.toResponse(&pedido)
	}

	return responses, nil
}

func (s *pedidoServiceImpl) FindByStatus(ctx context.Context, status string) ([]dto.PedidoResponse, error) {
	pedidos, err := s.pedidoRepo.FindByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PedidoResponse, len(pedidos))
	for i, pedido := range pedidos {
		responses[i] = *s.toResponse(&pedido)
	}

	return responses, nil
}

func (s *pedidoServiceImpl) UpdateStatus(ctx context.Context, id uint, req *dto.UpdatePedidoRequest) (*dto.PedidoResponse, error) {
	// Validar request
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Buscar pedido existente
	pedido, err := s.pedidoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Atualizar status
	pedido.Status = req.Status

	// Atualizar no banco
	if err := s.pedidoRepo.Update(ctx, pedido); err != nil {
		return nil, err
	}

	return s.toResponse(pedido), nil
}

func (s *pedidoServiceImpl) Delete(ctx context.Context, id uint) error {
	// Verificar se existe
	_, err := s.pedidoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.pedidoRepo.Delete(ctx, id)
}

func (s *pedidoServiceImpl) Count(ctx context.Context) (int64, error) {
	return s.pedidoRepo.Count(ctx)
}

// toResponse converte Model para Response DTO
func (s *pedidoServiceImpl) toResponse(pedido *model.Pedido) *dto.PedidoResponse {
	// Converter cliente
	var clienteResp *dto.ClienteResponse
	if pedido.Cliente.ID != 0 {
		clienteResp = &dto.ClienteResponse{
			ID:        pedido.Cliente.ID,
			Nome:      pedido.Cliente.Nome,
			Email:     pedido.Cliente.Email,
			CPF:       pedido.Cliente.CPF,
			Telefone:  pedido.Cliente.Telefone,
			CreatedAt: pedido.Cliente.CreatedAt.Format(time.RFC3339),
			UpdatedAt: pedido.Cliente.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Converter itens
	itens := make([]dto.ItemPedidoResponse, len(pedido.Itens))
	for i, item := range pedido.Itens {
		var produtoResp *dto.ProdutoResponse
		if item.Produto.ID != 0 {
			produtoResp = &dto.ProdutoResponse{
				ID:        item.Produto.ID,
				Nome:      item.Produto.Nome,
				Descricao: item.Produto.Descricao,
				Preco:     item.Produto.Preco,
				Estoque:   item.Produto.Estoque,
				SKU:       item.Produto.SKU,
				Categoria: item.Produto.Categoria,
				Ativo:     item.Produto.Ativo,
				CreatedAt: item.Produto.CreatedAt,
				UpdatedAt: item.Produto.UpdatedAt,
			}
		}

		itens[i] = dto.ItemPedidoResponse{
			ID:            item.ID,
			ProdutoID:     item.ProdutoID,
			Produto:       produtoResp,
			Quantidade:    item.Quantidade,
			PrecoUnitario: item.PrecoUnitario,
			Subtotal:      item.Subtotal,
		}
	}

	return &dto.PedidoResponse{
		ID:         pedido.ID,
		ClienteID:  pedido.ClienteID,
		Cliente:    clienteResp,
		Itens:      itens,
		ValorTotal: pedido.ValorTotal,
		Status:     pedido.Status,
		DataPedido: pedido.DataPedido,
		CreatedAt:  pedido.CreatedAt,
		UpdatedAt:  pedido.UpdatedAt,
	}
}
