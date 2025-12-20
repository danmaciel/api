package dto

// CreateClienteRequest represents the request body for creating a cliente
type CreateClienteRequest struct {
	Nome     string `json:"nome" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	CPF      string `json:"cpf" validate:"required,len=11,numeric"`
	Telefone string `json:"telefone" validate:"omitempty,min=10,max=15"`
}

// UpdateClienteRequest represents the request body for updating a cliente
type UpdateClienteRequest struct {
	Nome     string `json:"nome" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	CPF      string `json:"cpf" validate:"omitempty,len=11,numeric"`
	Telefone string `json:"telefone" validate:"omitempty,min=10,max=15"`
}

// ClienteResponse represents the response for cliente operations
type ClienteResponse struct {
	ID        uint   `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	CPF       string `json:"cpf"`
	Telefone  string `json:"telefone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CountResponse represents the count response
type CountResponse struct {
	Count int64 `json:"count"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
