# Split Payment API (Go + SQLite)

API REST que simula o fluxo de **split payment** B2B (segregação automática de
IBS/CBS + crédito tributário simplificado), com persistência em SQLite.

Inclui um **cliente HTTP para a Plataforma Pública do Split Payment** (RFB/CGIBS)
que pode ser usado com um mock local (Prism) para validar que os payloads seguem
o contrato real do governo.

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

## Sandbox de Integração (Fase 7)

O sandbox usa o **Prism** (Stoplight) para subir um mock local da Plataforma
Pública a partir do OpenAPI oficial (`cgibs/openapi-v0_0_10.json`).

### Pré-requisito

Node.js + npx (Prism é baixado automaticamente via npx).

### Como usar

**Terminal 1 — mock da PP:**

```bash
make sandbox
# Saída: servidor rodando em http://localhost:4010
```

**Terminal 2 — nossa API apontando pro mock:**

```bash
make run -pp-url http://localhost:4010 -pp-tenant PSP-SIMULADOR-001
```

Ou via variável de ambiente (não implementado — use flags por enquanto).

### O que acontece

Ao criar uma transação via `POST /api/v1/transactions` com `-pp-url` ativado,
a API dispara uma notificação síncrona para o mock:

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{"vendedor_cnpj":"11.111.111/0001-11","comprador_cnpj":"22.222.222/0001-22","valor_bruto":100000,"aliquota_ibs":0.12,"aliquota_cbs":0.03}'
```

Response inclui:

```json
{
  "transaction": { ... },
  "pp_notification": {
    "status": "sent",
    "arrangement": "boleto"
  }
}
```

Se `-pp-url` não for passado, `pp_notification.status` é `"skipped"`.

### Limitações do mock

- Prism valida schema e headers, mas **não executa lógica de negócio** (não corrige valores, não mantém estado).
- Endpoints de long polling com token de posição podem não funcionar plenamente no Prism.
- O fluxo de segregação real (2x/dia, lote, repasse financeiro) não é simulado.

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
| `make sandbox` | Sobe mock Prism da PP em http://localhost:4010 |
| `make help` | Lista todos os targets |

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
