# TASKS.md

Plano de execução. Muda entre sessões — edite aqui, não em `AGENT.md`. Consulte `PROGRESS.md` antes de
começar qualquer fase.

---

## Fase 0: Pivô de escopo — ✅ CONCLUÍDA (decisão registrada)

Projeto migrou de "CLI simples em memória" para "API REST + SQLite persistente", mantendo a mesma lógica
de domínio (cálculo de split, crédito tributário simplificado) já validada na versão anterior. Motivo:
tornar o projeto um artefato de portfólio mais realista — API é o formato que recrutadores/avaliadores
esperam ver, e "persistência de verdade" é a diferença entre brinquedo e sistema.

---

## Fase 1: Setup do projeto

- `go mod init` do módulo.
- Estrutura de pastas conforme `DIR.md`.
- `Makefile` com todos os targets listados em `DIR.md` (podem começar como stubs que evoluem nas fases
  seguintes).
- Escolher e justificar o driver SQLite em `REFERENCES.md` (preferir driver puro-Go, sem cgo, pra manter
  `make build` simples e portátil — decisão a registrar com a razão específica escolhida).

## Fase 2: Camada de domínio (portada do simulador CLI)

- `internal/domain/split.go` — `CalculateSplit`, igual à versão CLI (centavos, `int64`, erro se
  alíquotas somarem mais de 100%).
- `internal/domain/company.go` e `transaction.go` — structs, sem métodos de persistência (isso é
  repository).
- Lógica de crédito tributário simplificada (empresa compradora acumula crédito = IBS+CBS pago; ao
  vender, abate do que deve).
- Testes cobrindo: split normal, valor zero, alíquotas zero, alíquotas > 100%, crédito
  suficiente/insuficiente. Meta: essa é a camada com testes mais exaustivos do projeto (ver `AGENT.md`
  regra 6).

## Fase 3: Persistência

- `migrations/0001_init.sql` — tabelas `companies` (cnpj PK, nome, saldo_credito INTEGER) e
  `transactions` (id PK, vendedor_cnpj, comprador_cnpj, valor_bruto, aliquota_ibs, aliquota_cbs,
  valor_liquido, valor_ibs, valor_cbs, credito_usado, timestamp).
- `internal/repository/repository.go` — interfaces que o domínio/handler dependem (não a implementação
  concreta).
- `internal/repository/sqlite/` — implementação real com `database/sql`, queries parametrizadas (nunca
  concatenar SQL — injeção de SQL não é aceitável nem em projeto de portfólio).
- Testes de repository rodando contra um banco SQLite temporário (arquivo em `t.TempDir()` ou `:memory:`).

## Fase 4: HTTP

- `internal/handler/http/router.go` — `http.NewServeMux()` com os 7 endpoints de `AGENT.md`.
- `internal/handler/http/errors.go` — formato de erro único:
  ```json
  { "error": { "code": "VALIDATION_ERROR", "message": "descrição legível" } }
  ```
  Códigos HTTP: 400 (validação), 404 (não encontrado), 409 (conflito, ex: CNPJ duplicado), 500 (erro
  interno — mensagem genérica pro cliente, log detalhado no servidor).
- Handlers finos: decodificam JSON, chamam domínio + repository, codificam resposta. Nenhuma regra de
  negócio dentro de um handler.
- Testes de integração com `net/http/httptest` cobrindo happy path de cada endpoint + 1-2 erros
  relevantes (CNPJ inválido, empresa inexistente, valor negativo).

## Fase 5: Seed e demonstração

- `internal/seed/seed.go` — reproduz o cenário fictício já validado no simulador CLI: Fazenda → Fábrica
  → Mercado, com os mesmos valores (R$1.000 e R$3.000, IBS 12%/CBS 3%), agora persistido no banco.
- `make seed` roda isso contra um banco já migrado.
- Documentar no README um roteiro de `curl` (ou arquivo `.http`) que reproduz a demonstração via API,
  não só via seed direto no banco.

## Fase 6: Polish e documentação

- `gofmt`, `go vet` limpos.
- README com: o que é o projeto, como rodar (`make build && make migrate && make seed && make run`),
  lista de endpoints com exemplo de request/response, link pro vault de contexto tributário
  (`_agent-vault`/`vault/` do repo, se publicado junto).
- Conferir que nenhuma dependência externa ficou sem justificativa em `REFERENCES.md`.

---

## Definição de "Pronto"

- Todos os 7 endpoints de `AGENT.md` funcionando, com formato de erro padronizado.
- `make build`, `make test`, `make migrate`, `make seed`, `make run` funcionam do zero num checkout
  limpo.
- Camada de domínio com testes exaustivos; handlers com testes de integração básicos.
- Nenhuma regra de negócio vazando pra dentro de um handler HTTP ou de uma query SQL.
- README permite que alguém sem contexto rode a demonstração completa em menos de 5 comandos.

**GO.**
