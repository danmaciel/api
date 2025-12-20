# API MVC(Golang) - Sistema de Gerenciamento

> Uma API REST moderna em Go para gerenciar clientes, produtos e pedidos com arquitetura MVC

## Sobre o projeto?

Imagine que você precisa criar um sistema para uma loja online. Esse sistema precisa cadastrar clientes, gerenciar produtos e processar pedidos. É exatamente isso que esta API faz!

## O que você consegue fazer com esta API?

### Gerenciar Clientes
- Cadastrar novos clientes com nome, email, CPF e telefone
- Buscar clientes por nome ou ID
- Atualizar dados de clientes
- Remover clientes (soft delete - não apaga de verdade, só marca como inativo)

### Controlar Produtos
- Adicionar produtos com preço, estoque e categoria
- Atualizar informações e quantidade em estoque
- Buscar produtos por categoria ou nome
- Marcar produtos como ativos ou inativos

### Processar Pedidos
- Criar pedidos associando clientes e produtos
- Validar estoque automaticamente
- Calcular valor total do pedido
- Acompanhar status (pendente → confirmado → enviado → entregue)
- Cancelar pedidos quando necessário

## Como o projeto está organizado?

A API segue uma arquitetura em camadas bem definida:

![Diagrama de Contexto](docs/diagramas/1-diagrama_C4_Contexto.png)

*Visão geral: como o sistema interage com usuários e o banco de dados*

![Diagrama de Container](docs/diagramas/2-diagrama_C4_Container.png)

*Dentro da aplicação: principais componentes e suas responsabilidades*

### Estrutura de Pastas

```
api/
├── cmd/api/                 # Ponto de entrada da aplicação
├── internal/
│   ├── controller/          # Recebe as requisições HTTP (a porta de entrada)
│   ├── service/             # Processa a lógica de negócio (o cérebro)
│   ├── repository/          # Conversa com o banco de dados (o arquivo)
│   ├── model/               # Define as entidades (Cliente, Produto, Pedido)
│   ├── dto/                 # Formatos de dados para entrada/saída
│   └── middleware/          # Filtros e validações
├── config/                  # Configurações e conexão com o banco
├── docs/                    # Documentação automática (Swagger), collection do postaman para testes e 
                             # diagramas do projeto(c4 model, sequencia e er)
└── tests/                   # Testes automatizados
```

![Diagrama de Componentes](docs/diagramas/3-diagrama_C4_Componentes.png)
*Detalhamento: como os componentes internos se conectam*

## Padrões de Design (Design Patterns)

Este projeto implementa diversos padrões de design reconhecidos pela indústria, tornando o código mais organizado, testável e manutenível.

### Padrões Arquiteturais

**Arquitetura em Camadas**
- **Onde**: Estrutura completa do `internal/`
- **Por quê**: Separa responsabilidades em camadas (Controller → Service → Repository → Model)
- **Benefício**: Mudanças em uma camada não afetam as outras

**MVC (Model-View-Controller)**
- **Onde**: `internal/controller/` (Controllers), `internal/model/` (Models), `internal/dto/` (Views/DTOs)
- **Por quê**: Organiza a aplicação separando lógica de apresentação, negócio e dados
- **Benefício**: Código mais organizado e fácil de manter

### Padrões de Criação

**Factory Method Pattern**
- **Onde**: Funções `New*` como `NewClienteController()`, `NewClienteService()`, `NewClienteRepositorySQLite()`
  - `internal/controller/cliente_controller.go:17-21`
  - `internal/service/cliente_service_impl.go:20-26`
  - `internal/repository/cliente_repository_sqlite.go:15-18`
- **Por quê**: Centraliza a criação de objetos complexos
- **Benefício**: Garante que objetos sejam criados corretamente com todas as dependências

### Padrões Estruturais

**Repository Pattern**
- **Onde**: Todo o `internal/repository/`
  - Interfaces: `cliente_repository.go`, `produto_repository.go`, `pedido_repository.go`
  - Implementações: `*_repository_sqlite.go`
- **Por quê**: Abstrai o acesso aos dados, isolando a lógica de persistência
- **Benefício**: Fácil trocar SQLite por PostgreSQL/MySQL sem alterar Services

**Adapter Pattern**
- **Onde**: Implementações de Repository (ex: `clienteRepositorySQLite`)
- **Por quê**: Adapta a interface genérica do repositório para o GORM
- **Benefício**: Desacopla a aplicação de detalhes de implementação do banco

**Facade Pattern**
- **Onde**: `internal/controller/router.go:SetupRouter()`
- **Por quê**: Simplifica a configuração complexa de rotas, middlewares e CORS
- **Benefício**: Ponto único de configuração do servidor HTTP

**Decorator Pattern (Middleware)**
- **Onde**: `internal/middleware/middleware.go`
  - Logger, Recovery, ContentType
- **Por quê**: Adiciona funcionalidades (logging, recuperação de erros) sem modificar handlers
- **Benefício**: Funcionalidades transversais aplicadas de forma limpa

### Padrões Comportamentais

**Dependency Injection (DI)**
- **Onde**: `cmd/api/main.go:44-61`
- **Por quê**: Injeta dependências via construtor em vez de criá-las internamente
- **Benefício**: Código altamente testável e com baixo acoplamento
```go
// Exemplo: Repositórios → Services → Controllers
clienteRepo := repository.NewClienteRepositorySQLite(db)
clienteService := service.NewClienteService(clienteRepo)
clienteController := controller.NewClienteController(clienteService)
```

**Strategy Pattern**
- **Onde**: Interfaces de Service (`ClienteService`, `ProdutoService`, `PedidoService`)
- **Por quê**: Define contratos de negócio que podem ter múltiplas implementações
- **Benefício**: Fácil criar implementações alternativas (ex: mock para testes)

**Chain of Responsibility**
- **Onde**: Cadeia de middlewares em `router.go:17-21`
- **Por quê**: Requisições passam por vários middlewares sequencialmente
- **Benefício**: Cada middleware processa e passa adiante
```go
r.Use(chimiddleware.RequestID)  // 1. Gera ID único
r.Use(middleware.Logger)         // 2. Registra log
r.Use(middleware.Recovery)       // 3. Recupera de panics
r.Use(middleware.ContentType)    // 4. Define Content-Type
```

**Mapper/Transformer Pattern**
- **Onde**: Métodos `toResponse()` nos Services
  - `internal/service/cliente_service_impl.go:143-154`
  - `internal/service/pedido_service_impl.go:197-253`
- **Por quê**: Converte entidades de domínio (`model`) para DTOs de resposta
- **Benefício**: Controla exatamente quais dados são expostos pela API

### Padrões de Dados

**DTO (Data Transfer Object)**
- **Onde**: `internal/dto/`
  - `CreateClienteRequest`, `UpdateClienteRequest`, `ClienteResponse`
- **Por quê**: Separa modelos de domínio dos dados de entrada/saída da API
- **Benefício**: Validação e controle independentes, protege modelo interno

**Active Record Hooks (GORM)**
- **Onde**: `internal/model/pedido.go:48-52` - `BeforeSave()`
- **Por quê**: Executa lógica antes de salvar (ex: calcular subtotal)
- **Benefício**: Lógica de domínio próxima ao modelo

### Resumo dos Padrões

| Categoria | Padrão | Localização Principal | Vantagem Chave |
|-----------|--------|----------------------|----------------|
| Arquitetural | Clean Architecture | `internal/*` | Separação de responsabilidades |
| Arquitetural | MVC | Controller/Service/Model | Organização clara |
| Criação | Factory Method | Funções `New*()` | Criação consistente |
| Estrutural | Repository | `internal/repository/` | Abstração de dados |
| Estrutural | Decorator | `internal/middleware/` | Funcionalidades modulares |
| Comportamental | Dependency Injection | `main.go` | Testabilidade |
| Comportamental | Strategy | Interfaces de Service | Flexibilidade |
| Dados | DTO | `internal/dto/` | Controle de API |

## Como funciona na prática?

### Exemplo: Criar um Pedido

Quando você cria um pedido, a API faz todo este processo automaticamente:

![Fluxo de Criação de Pedido](docs/diagramas/6-diagrama_de_sequencia_criar_pedido.png)

1. **Você envia**: Cliente ID + lista de produtos e quantidades
2. **API valida**: O cliente existe? Os produtos existem? Tem estoque suficiente?
3. **API calcula**: Valor total do pedido (quantidade × preço)
4. **API salva**: Pedido no banco de dados com todos os detalhes
5. **Você recebe**: Confirmação com número do pedido e total a pagar

### Outros Fluxos

**Operações com Cliente:**
![Operações de Cliente](docs/diagramas/4-diagrama_de_sequencia_operacoes_de_cliente.png)
*Criar, buscar, atualizar e deletar clientes*

**Operações com Produto:**
![Operações de Produto](docs/diagramas/5-diagrama_de_sequencia_operacoes_de_produto.png)
*Gerenciar produtos e estoque*

**Outras Operações de Pedido:**
![Outras Operações de Pedido](docs/diagramas/7-diagrama_de_sequencia_outras_operacoes_pedido.png)
*Consultar, atualizar status, cancelar e gerar relatórios*

## Banco de Dados

A API usa SQLite (um banco de dados simples em arquivo) com a seguinte estrutura:

![Diagrama ER](docs/diagramas/8-diagrama_er.png)
*Estrutura das tabelas e relacionamentos*

**Principais relacionamentos:**
- Um **Cliente** pode ter vários **Pedidos**
- Um **Pedido** pode ter vários **Produtos** (e vice-versa)
- Cada item do pedido guarda o preço no momento da compra (histórico)

## Como executar o projeto?

### Pré-requisitos
- Go 1.21 ou superior instalado
- Um terminal/console

### Passos

1️⃣ **Clone o repositório**
```bash
git clone <url-do-repositorio>
cd api
```

2️⃣ **Instale as dependências**
```bash
make install
# Ou: go mod download
```

3️⃣ **Configure as variáveis de ambiente**

Crie um arquivo `.env` na raiz do projeto:
```env
SERVER_PORT=8080
DB_FILE_PATH=./database/api.db
```

4️⃣ **Execute a aplicação**
```bash
make run
# Ou: go run cmd/api/main.go
```

Pronto! A API estará rodando em `http://localhost:8080`

## Documentação Interativa

Depois de iniciar a aplicação, acesse:

**http://localhost:8080/swagger/index.html**

Lá você encontra:
- Todos os endpoints disponíveis
- Como usar cada um deles
- Exemplos de requisições e respostas
- Possibilidade de testar diretamente pelo navegador

## Exemplos de Uso

### Criar um Cliente
```bash
curl -X POST http://localhost:8080/api/v1/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Maria Silva",
    "email": "maria@example.com",
    "cpf": "12345678901",
    "telefone": "11999999999"
  }'
```

**Resposta:**
```json
{
  "id": 1,
  "nome": "Maria Silva",
  "email": "maria@example.com",
  "cpf": "12345678901",
  "telefone": "11999999999",
  "created_at": "2025-12-17T15:30:00Z"
}
```

### Criar um Produto
```bash
curl -X POST http://localhost:8080/api/v1/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Notebook Dell",
    "preco": 3500.00,
    "estoque": 10,
    "sku": "DELL-NB-001",
    "categoria": "Eletrônicos"
  }'
```

### Criar um Pedido
```bash
curl -X POST http://localhost:8080/api/v1/pedidos \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_id": 1,
    "itens": [
      {
        "produto_id": 1,
        "quantidade": 2
      }
    ]
  }'
```

**Resposta:**
```json
{
  "id": 1,
  "cliente_id": 1,
  "cliente": {
    "id": 1,
    "nome": "Maria Silva",
    "email": "maria@example.com"
  },
  "itens": [
    {
      "id": 1,
      "produto_id": 1,
      "produto": {
        "id": 1,
        "nome": "Notebook Dell",
        "preco": 3500.00
      },
      "quantidade": 2,
      "preco_unitario": 3500.00,
      "subtotal": 7000.00
    }
  ],
  "valor_total": 7000.00,
  "status": "pendente",
  "data_pedido": "2025-12-17T15:35:00Z"
}
```

## Comandos Úteis

```bash
make help              # Ver todos os comandos disponíveis
make run               # Executar a aplicação
make test              # Rodar testes
make test-coverage     # Rodar testes com relatório de cobertura
make build             # Compilar a aplicação
make swagger           # Gerar documentação Swagger
make clean             # Limpar arquivos de build
```

## Tecnologias Utilizadas

| Tecnologia | Para que serve? |
|------------|-----------------|
| **Go** | Linguagem de programação rápida e eficiente |
| **Chi Router** | Roteamento HTTP leve e flexível |
| **GORM** | Facilita operações com banco de dados |
| **SQLite** | Banco de dados simples em arquivo |
| **Swagger** | Documentação interativa da API |
| **Validator** | Validação automática de dados |
| **Testify** | Framework de testes |

## Principais Funcionalidades

- **Validação automática**: Dados são validados antes de serem salvos
- **Soft Delete**: Registros "deletados" ficam marcados, não são removidos
- **Relacionamentos**: Pedidos conectam clientes e produtos automaticamente
- **Status HTTP corretos**: 200 OK, 201 Created, 404 Not Found, etc.
- **CORS habilitado**: Pode ser acessada de qualquer frontend
- **Logs completos**: Todas as requisições são registradas
- **Graceful Shutdown**: Encerra conexões de forma limpa
- **Testes automatizados**: Garantem que tudo funciona

## Endpoints Disponíveis

### Clientes (7 endpoints)
- `POST /api/v1/clientes` - Criar cliente
- `GET /api/v1/clientes` - Listar todos
- `GET /api/v1/clientes/{id}` - Buscar por ID
- `GET /api/v1/clientes/nome/{nome}` - Buscar por nome
- `GET /api/v1/clientes/count` - Contar total
- `PUT /api/v1/clientes/{id}` - Atualizar
- `DELETE /api/v1/clientes/{id}` - Deletar

### Produtos (8 endpoints)
- `POST /api/v1/produtos` - Criar produto
- `GET /api/v1/produtos` - Listar todos
- `GET /api/v1/produtos/{id}` - Buscar por ID
- `GET /api/v1/produtos/categoria/{cat}` - Buscar por categoria
- `GET /api/v1/produtos/count` - Contar total
- `PATCH /api/v1/produtos/{id}/estoque` - Atualizar estoque
- `PUT /api/v1/produtos/{id}` - Atualizar
- `DELETE /api/v1/produtos/{id}` - Deletar

### Pedidos (12 endpoints)
- `POST /api/v1/pedidos` - Criar pedido
- `GET /api/v1/pedidos` - Listar todos
- `GET /api/v1/pedidos/{id}` - Buscar por ID
- `GET /api/v1/pedidos/cliente/{id}` - Pedidos de um cliente
- `GET /api/v1/pedidos/status/{status}` - Filtrar por status
- `PATCH /api/v1/pedidos/{id}/status` - Atualizar status
- `DELETE /api/v1/pedidos/{id}` - Cancelar pedido
- E mais...

### Utilitários
- `GET /health` - Verificar se a API está funcionando
- `GET /swagger/*` - Documentação interativa


## Autor

**Danilo Maciel(dan maciel)**

---

**Versão:** 1.0.0
**Data:** 17/12/2025
**Projeto Acadêmico** - API REST em Go com padrão MVC
