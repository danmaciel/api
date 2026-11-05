# ADR 0001 — Consistência de estoque ao criar pedidos

- **Status:** proposta (não implementada)
- **Contexto:** API de demonstração em Go/GORM/SQLite
- **Escopo:** criação de pedido e disponibilidade de produtos

## Problema

Atualmente a camada de serviço lê o estoque e valida as quantidades antes de persistir o pedido. Não há reserva/baixa atômica de estoque. Duas solicitações simultâneas podem validar a mesma disponibilidade e criar pedidos sem uma garantia de consistência. Também não há regra de reposição na exclusão/cancelamento.

## Decisão proposta (não aplicada)

1. Definir a semântica de negócio: ao criar pedido pendente, o estoque é reservado, vendido ou permanece disponível? Estabelecer política explícita para cancelamento, expiração e estorno.
2. Criar uma operação transacional coesa que combine a validação do cliente, agregação de quantidades por produto, verificação de ativação, baixa/reserva condicional e gravação do pedido e seus itens.
3. Realizar cada atualização de estoque por comando condicional dentro da mesma transação, por exemplo, `UPDATE produtos SET estoque = estoque - ? WHERE id = ? AND ativo = true AND estoque >= ?`, verificando `RowsAffected == 1`.
4. Se uma atualização falhar, reverter toda a transação e retornar erro de domínio que diferencie indisponibilidade de estoque, produto inexistente e conflito temporário.
5. Testar cenários com itens duplicados no mesmo pedido, falha na persistência após baixa, cancelamento, limite exato do estoque e duas requisições concorrentes.

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

> Este documento descreve uma decisão arquitetural em estudo; não afirma que as garantias acima já estejam implementadas.
