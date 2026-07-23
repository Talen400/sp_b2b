# PROGRESS.md

Estado atual do projeto. Leia antes de gerar/alterar qualquer código; atualize ao final de cada rodada.

---

## Status por Fase

- [x] Fase 0 — Pivô de escopo (CLI → API + banco)
- [x] Fase 1 — Setup do projeto
- [x] Fase 2 — Camada de domínio
- [x] Fase 3 — Persistência
- [x] Fase 4 — HTTP
- [x] Fase 5 — Seed e demonstração
- [x] Fase 6 — Polish e documentação

## Estrutura final

```
├── agent/                      ← AGENT.md, DIR.md, REFERENCES.md, TASKS.md, PROGRESS.md
├── vault/                      ← 14 notas Obsidian (projeto irmão de estudo)
├── cmd/api/main.go             ← entrypoint da API
├── internal/
│   ├── domain/                 ← regras de negócio puras
│   │   ├── split.go + split_test.go      ← CalculateSplit (7 testes)
│   │   ├── company.go                     ← struct Company
│   │   ├── transaction.go                 ← struct Transaction
│   │   ├── credit.go                      ← ApplyCredit / UseCredit
│   │   └── credit_test.go                 ← crédito (5 testes)
│   ├── repository/
│   │   ├── repository.go       ← interfaces CompanyRepo / TransactionRepo
│   │   └── sqlite/
│   │       ├── sqlite.go       ← implementação SQLite (embedded migrations)
│   │       ├── sqlite_test.go  ← testes repo (8 testes)
│   │       └── migrations/0001_init.sql
│   ├── handler/http/
│   │   ├── router.go           ← ServeMux com 7 endpoints
│   │   ├── company_handler.go
│   │   ├── transaction_handler.go
│   │   ├── healthz_handler.go
│   │   ├── errors.go           ← formato JSON de erro padronizado
│   │   └── http_test.go        ← testes de integração (7 testes)
│   └── seed/seed.go            ← cenário de demonstração
├── migrations/                 ← (realocado para internal/repository/sqlite/migrations/)
├── data/.gitkeep               ← SQLite é criado aqui em runtime
├── Makefile
├── README.md
└── LICENSE
```

## Decisões registradas

- **Driver SQLite:** `modernc.org/sqlite` (puro Go, sem cgo) — registrado em `REFERENCES.md`
- **Porta default:** 8080
- **Separate repository types:** `CompanyRepository` e `TransactionRepository` (evita conflito de assinatura `Get`)
- **ID de transação:** baseado em `time.Now().UnixNano()`, prefixo `TXN-`
- **Migrations:** embutidas via `//go:embed` no pacote `sqlite`, aplicadas na inicialização

## Cobertura de testes

| Camada | Testes | Status |
|---|---|---|
| Domain (split) | 7 | ✅ |
| Domain (crédito) | 5 | ✅ |
| Repository (SQLite) | 8 | ✅ |
| Handler (HTTP) | 7 | ✅ |
| **Total** | **27** | ✅ |

## Observações da Última Rodada

- Projeto pivô concluído: CLI em memória → API REST + SQLite
- Código legado (`internal/split/`, `internal/company/`, `internal/store/`, `cmd/simulador/`, `cmd/demo/`) removido
- `go vet` e `gofmt` limpos
- `make migrate && make seed` funcional — cenário Fazenda → Fábrica → Mercado verificado
- Endpoints testados individualmente via testes de integração com `httptest`
