# ADR 0001 — Consistência de estoque ao criar pedidos

- **Status:** parcialmente implementada — criação de pedidos com baixa transacional; cancelamento/estorno pendente
- **Contexto:** API de demonstração em Go/GORM/SQLite
- **Escopo:** criação de pedido e disponibilidade de produtos

## Problema

A implementação original lia o estoque antes de persistir o pedido, sem baixa transacional. Esta etapa adiciona uma atualização condicional de estoque e a inserção do pedido na mesma transação de banco, executadas no repositório. Ainda não há reposição automática de estoque para cancelamentos ou exclusões.

## Decisão e implementação por etapas

1. **Implementado nesta etapa:** ao criar um pedido (inclusive pendente), as quantidades são imediatamente baixadas do estoque. Não há reservas com expiração.
2. **Implementado parcialmente:** o serviço valida o cliente e os itens; o repositório agrega quantidades por produto e executa baixa condicional e gravação do pedido e itens dentro de uma transação.
3. **Implementado:** realizar cada atualização de estoque por comando condicional dentro da mesma transação, por exemplo, `UPDATE produtos SET estoque = estoque - ? WHERE id = ? AND ativo = true AND estoque >= ?`, verificando `RowsAffected == 1`.
4. **Implementado parcialmente:** erro de atualização condicional reverte a transação. Diferenciação refinada de erros de domínio e respostas HTTP específicas ainda pendentes.
5. **Coberto nesta etapa:** itens duplicados, falha ao persistir pedido, estoque insuficiente e duas operações concorrentes em SQLite com conexões serializadas. **Pendente:** estorno no cancelamento, bloqueios entre múltiplas conexões e testes de carga.

## SQLite versus PostgreSQL

SQLite é adequado para estudo e execução local, mas tem particularidades de escrita concorrente e bloqueio. O tratamento de `SQLITE_BUSY`, o tempo de espera e a persistência em arquivo devem ser validados separadamente. Para um sistema comercial com maior concorrência, considerar PostgreSQL e transações com controle de concorrência apropriado.

## Alternativas e trade-offs

- **Mutex em memória:** simples, mas não protege processos ou réplicas diferentes; rejeitado como garantia principal.
- **Transação com atualização condicional:** evita overselling ao combinar teste e modificação de saldo no banco, sem depender de leitura anterior não protegida; abordagem preferida.
- **Reserva com prazo de expiração:** permite checkout assíncrono, porém exige modelagem de reservas, idempotência e processo de expiração.

## Critérios de aceite

- Nunca persistir um pedido se não for possível reservar/baixar todas as quantidades necessárias.
- Nenhum estoque negativo após solicitações concorrentes.
- Rollback completo se falhar qualquer item ou a criação do pedido.
- Comportamento de cancelamento e reembolso documentado e testado.
- Teste de concorrência automatizado usando banco compartilhado, sem depender apenas de mocks.

> A baixa transacional é implementada apenas na criação do pedido. Atualização de status e exclusão ainda não devolvem unidades ao estoque; não utilizar esse fluxo como checkout completo antes de definir e testar cancelamentos.
