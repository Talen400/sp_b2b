# DIR.md

Estrutura do repositório do Split Payment API. Camadas separadas fisicamente em pastas, não só por
convenção de nome — ver `AGENT.md` regra 4.

```
split-payment-api/
├── AGENT.md, DIR.md, REFERENCES.md, TASKS.md, PROGRESS.md   ← controle (este conjunto)
├── Makefile                        ← interface oficial (build, run, test, seed, migrate, clean)
├── go.mod / go.sum
├── README.md                       ← como rodar, o que é split payment, link pro vault de contexto
├── cmd/
│   └── api/
│       └── main.go                 ← entrypoint: lê config, conecta DB, sobe HTTP server
├── internal/
│   ├── domain/
│   │   ├── company.go              ← struct Company + regras puras de crédito
│   │   ├── transaction.go          ← struct Transaction
│   │   ├── split.go                ← CalculateSplit (herdado do simulador CLI, sem mudança de lógica)
│   │   ├── split_test.go
│   │   └── credit_test.go
│   ├── repository/
│   │   ├── repository.go           ← interfaces (CompanyRepo, TransactionRepo) — domínio depende só disso
│   │   └── sqlite/
│   │       ├── company_repo.go
│   │       ├── transaction_repo.go
│   │       └── sqlite_test.go
│   ├── handler/
│   │   └── http/
│   │       ├── router.go           ← monta o ServeMux com todas as rotas
│   │       ├── company_handler.go
│   │       ├── transaction_handler.go
│   │       ├── healthz_handler.go
│   │       ├── errors.go           ← formato de erro JSON padronizado (AGENT.md regra 5)
│   │       └── http_test.go        ← testes de integração leves (httptest)
│   └── seed/
│       └── seed.go                 ← popula o cenário fictício de demonstração (3 empresas, 2 vendas)
├── migrations/
│   └── 0001_init.sql               ← schema inicial (companies, transactions)
└── data/
    └── .gitkeep                    ← onde o arquivo .db do SQLite é criado em runtime (git-ignorado)
```

## Targets do Makefile (contrato mínimo)
- `make build` — compila o binário em `bin/api`.
- `make run` — sobe a API localmente (porta default documentada no README).
- `make test` — roda `go test ./...`.
- `make migrate` — aplica migrations pendentes no arquivo `.db`.
- `make seed` — roda o cenário de demonstração contra o banco já migrado.
- `make clean` — remove `bin/` e o arquivo `.db` local.
- `make re` — `clean` + `build` (nome inspirado no `make re` clássico da 42).

Se uma nova pasta for necessária, proponha aqui antes de criar — mesma regra do vault: sem estrutura
ad-hoc não registrada.
