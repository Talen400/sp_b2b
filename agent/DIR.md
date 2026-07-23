# DIR.md

Estrutura do repositório do Split Payment API. Camadas separadas fisicamente em pastas, não só por
convenção de nome — ver `AGENT.md` regra 4.

```
split-payment-api/
├── AGENT.md, DIR.md, REFERENCES.md, TASKS.md, PROGRESS.md   ← controle (este conjunto)
├── .gitignore                      ← ignora cgibs/*.pdf, .docx, .zip; data/*.db; bin/; /.sandbox/
├── Makefile                        ← interface oficial (build, run, test, seed, migrate, clean, sandbox)
├── go.mod / go.sum
├── README.md                       ← como rodar, o que é split payment, link pro vault de contexto
├── cgibs/                          ← documentos oficiais (openapi-v0_0_10.json versionado; PDF/DOCX/ZIP ignorados)
│   └── openapi-v0_0_10.json        ← OpenAPI real da Plataforma Pública, fonte do mock Prism
├── cmd/
│   └── api/
│       └── main.go                 ← entrypoint: flags -port, -db, -migrate, -seed, -pp-url
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
│   │       ├── sqlite.go           ← DB wrapper, Open/Migrate, CompanyRepo + TransactionRepo (embed das migrations)
│   │       ├── sqlite_test.go
│   │       └── migrations/
│   │           └── 0001_init.sql   ← schema embutido no binário via //go:embed
│   ├── handler/
│   │   └── http/
│   │       ├── router.go           ← monta o ServeMux com todas as rotas
│   │       ├── company_handler.go
│   │       ├── transaction_handler.go  ← hook opcional p/ PP client
│   │       ├── healthz_handler.go
│   │       ├── errors.go           ← formato de erro JSON padronizado (AGENT.md regra 5)
│   │       └── http_test.go        ← testes de integração leves (httptest)
│   ├── client/
│   │   └── pp/
│   │       ├── doc.go              ← documentação do pacote
│   │       ├── types.go            ← structs mapeando schemas reais do OpenAPI PP
│   │       ├── client.go           ← HTTP client contra a Plataforma Pública (mock)
│   │       └── client_test.go      ← testes de serialização + integração opcional c/ Prism
│   └── seed/
│       └── seed.go                 ← popula o cenário fictício de demonstração (3 empresas, 2 vendas)
└── data/
    └── .gitkeep                    ← onde o arquivo .db do SQLite é criado em runtime (git-ignorado)
```

## Targets do Makefile (contrato mínimo)

- `make build` — compila o binário em `bin/api`.
- `make run` — sobe a API localmente (porta 8080).
- `make test` — roda `go test ./...`.
- `make migrate` — aplica migrations embutidas.
- `make seed` — roda o cenário de demonstração contra o banco já migrado.
- `make clean` — remove `bin/` e o arquivo `.db` local.
- `make re` — `clean` + `build`.
- `make sandbox` — sobe mock Prism da Plataforma Pública na porta 4010.
  Requer `cgibs/openapi-v0_0_10.json` presente. Falha com mensagem clara se não existir.
- `make help` — lista todos os targets.

Se uma nova pasta for necessária, proponha aqui antes de criar — mesma regra do vault: sem estrutura
ad-hoc não registrada.
