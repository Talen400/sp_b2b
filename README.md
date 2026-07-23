# Split Payment API (Go + SQLite)

API REST que simula o fluxo de **split payment** B2B (segregação automática de
IBS/CBS + crédito tributário simplificado), com persistência em SQLite.

Projeto irmão do vault de contexto tributário em `vault/` (notas Obsidian sobre
a Reforma Tributária).

## Stack

- Go 1.25+ (stdlib `net/http`, `database/sql`)
- SQLite (`modernc.org/sqlite` — driver puro-Go, sem cgo)

## Como rodar

```bash
# 1. Build
make build

# 2. Rodar migrations + seed + servidor
make migrate && make seed && make run

# Tudo em um comando (build + migrate já acontecem automaticamente)
make run
```

## Endpoints

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/healthz` | Health check |
| `POST` | `/api/v1/companies` | Criar empresa |
| `GET` | `/api/v1/companies/{cnpj}` | Detalhe + saldo de crédito |
| `GET` | `/api/v1/companies` | Listar empresas |
| `POST` | `/api/v1/transactions` | Simular venda |
| `GET` | `/api/v1/transactions/{id}` | Detalhe de transação |
| `GET` | `/api/v1/transactions?cnpj={cnpj}` | Histórico (com filtro opcional) |

### Exemplos com curl

```bash
# Criar empresas
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"cnpj":"11.111.111/0001-11","nome":"Fazenda Boa Vista"}'

curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"cnpj":"22.222.222/0001-22","nome":"Fábrica de Sucos SA"}'

# Simular venda (IBS 12%, CBS 3%)
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{"vendedor_cnpj":"11.111.111/0001-11","comprador_cnpj":"22.222.222/0001-22","valor_bruto":100000,"aliquota_ibs":0.12,"aliquota_cbs":0.03}'

# Consultar saldo
curl http://localhost:8080/api/v1/companies/22.222.222/0001-22

# Histórico
curl http://localhost:8080/api/v1/transactions
```

## Makefile

| Target | Descrição |
|---|---|
| `make build` | Compila binário em `bin/api` |
| `make run` | Sobe servidor (porta 8080) |
| `make test` | Roda todos os testes |
| `make migrate` | Aplica migrations |
| `make seed` | Popula cenário de demonstração |
| `make clean` | Remove bin/ e banco |
| `make re` | Clean + build |

## Sobre o split payment

Split payment é o mecanismo da Reforma Tributária (EC 132/2023 + LC 214/2025)
que segrega automaticamente o IBS e CBS no momento do pagamento. Em transações
B2B, o comprador acumula crédito tributário do IBS+CBS pagos, que pode ser
abatido em vendas futuras (não-cumulatividade).

Esta API é uma **simulação didática** — não integra com a Plataforma Pública
real da RFB/CGIBS. Para contexto tributário completo, consulte as notas em
`vault/`.

## Licença

MIT
