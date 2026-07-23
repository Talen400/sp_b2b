# Fluxo de Informes e Endpoints

## TL;DR
A comunicação com a Plataforma Pública se dá através de **informes** — mensagens HTTP que notificam a PP sobre cada etapa da vida de uma transação. Cada arranjo tem seu próprio conjunto de informes obrigatórios e seu endpoint específico. Este documento detalha todos eles.

## Tabela completa: arranjos × informes

Legenda: M = Obrigatório, O = Opcional, N/A = Não se aplica

| Informe | Boleto | Pix Din. | Pix Auto. | Pix Est. | TED | TEF |
|---|---|---|---|---|---|---|
| Transação Iniciada | M | M | M | N/A | N/A | N/A |
| Transação Atualizada | M | M | N/A | N/A | N/A | N/A |
| Retorno SI — Correção | O | O | O | N/A | N/A | N/A |
| Retorno SI — Em Aberto | O | O | O | N/A | N/A | N/A |
| Informe de Baixa | M | M | M | N/A | N/A | N/A |
| Informe Preliminar de Pagamento | M | M | M | M | M | O |
| Informe de Segregação | M | M | M | M | M | M |

## Endpoints por arranjo

### Informe de Transação Iniciada
Enviado pelo PSP Recebedor ao criar a cobrança.

| Arranjo | Método | Path |
|---|---|---|
| Boleto | `POST` | `/api/v1/boleto` |
| Pix Dinâmico | `POST` | `/api/v1/pix-dinamico` |
| Pix Automático | `POST` | `/api/v1/pix-automatico` |

**Payload (Boleto — exemplo simplificado):**
```json
{
  "nsuId": "123456",
  "cnpjRec": "11.111.111/0001-11",
  "vlInf": "1000.00",
  "vlIbs": "120.00",
  "vlCbs": "30.00",
  "dtHrCriacao": "2026-07-23T10:00:00Z"
}
```

### Informe de Transação Atualizada
Enviado quando há correção no documento fiscal antes do pagamento.

| Arranjo | Método | Path |
|---|---|---|
| Boleto | `PUT` | `/api/v1/boleto/{nsuId}` |
| Pix Dinâmico | `PUT` | `/api/v1/pix-dinamico/{nsuId}` |

### Informe Preliminar de Pagamento
Enviado logo após o pagamento. Informativo, não vinculante.

| Arranjo | Método | Path |
|---|---|---|
| Boleto | `POST` | `/api/v1/boleto` (mesmo endpoint do Iniciada, campo `indPgto` indicação) |
| Pix Dinâmico | `POST` | `/api/v1/pix-dinamico` |
| Pix Automático | `POST` | `/api/v1/pix-automatico` |
| Pix Estático | `POST` | `/api/v1/pix-estatico` |
| TED | `POST` | `/api/v1/ted` |
| TEF | `POST` | `/api/v1/tef` |

**Payload (Pix Estático — exemplo):**
```json
{
  "nsuId": "789012",
  "cnpjRec": "11.111.111/0001-11",
  "cnpjCpfPagEfet": "22.222.222/0001-22",
  "vlInf": "500.00",
  "indPgtoIntegral": true,
  "dtHrPgto": "2026-07-23T14:30:00Z"
}
```

### Informe de Segregação (lote)
Enviado 2x por dia útil. Definitivo, vinculante. Gera obrigação de Repasse Financeiro.

| Método | Path |
|---|---|
| `POST` | `/api/v1/segregacao` |

**Payload (exemplo):**
```json
{
  "cnpjRaizPspRecDir": "12345678000199",
  "dtHrMsg": "2026-07-23T16:00:00Z",
  "idInfSegr": "uuid-do-informe",
  "transacoes": [
    {
      "nsuId": "123456",
      "vlIbsSeg": "120.00",
      "vlCbsSeg": "30.00"
    }
  ],
  "totalTrans": 1,
  "valorTotalCbs": "30.00",
  "valorTotalIbs": "120.00"
}
```

### Retorno Super Inteligente — Long Polling

Fluxo de consulta:

```
PSP                              PP
 │                                │
 │  GET /start                    │
 │───────────────────────────────>│
 │                                │── timeout ou mensagem
 │<───────────────────────────────│
 │  msgs + header: proximoToken   │
 │                                │
 │  GET /{token}                  │
 │───────────────────────────────>│
 │                                │── timeout ou mensagem
 │<───────────────────────────────│
 │  msgs + header: proximoToken   │
 │                                │
 │  ... (repetir)                 │
 │                                │
 │  DELETE /{token}               │
 │───────────────────────────────>│
 │<───────────────────────────────│
 │  204 No Content                │
```

**Start:** `GET /api/v1/{arrj}/{idPsp}/tributos/stream/start`
**Continue:** `GET /api/v1/{arrj}/{idPsp}/tributos/stream/{token}`
**End:** `DELETE /api/v1/{arrj}/{idPsp}/tributos/stream/{token}`

`arrj` = `boleto`, `pix-dinamico`, `pix-automatico`
`idPsp` = identificador do PSP na PP

**Mensagem de retorno (exemplo):**
```json
{
  "tributos": [
    {
      "nsuId": "123456",
      "codMsg": "CORRECAO",
      "idDda": "dda-123",
      "numCtrlOrig": "ctrl-456",
      "vlInf": "1000.00",
      "vlCbsCorr": "28.50",
      "vlIbsCorr": "114.00",
      "docFiscal": "NF-e 12345"
    }
  ]
}
```

Se não houver mensagens até o timeout, PP responde `204 No Content` com `proximoToken` no header.

### Consulta Retroativa Super Inteligente

Permite consultar mensagens já entregues anteriormente.

**Start com intervalo de NSU:**
`GET /api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start?fromNsu=501&toNsu=550`

**Start com stream:**
`GET /api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start?fromNsu=501&toNsu=550&streamId=A`

**Start com stream encerrada:**
`GET /api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start?fromNsu=501&streamId=A`

O fluxo é o mesmo (start → continue → delete), mas os tokens podem referenciar mensagens antigas.

## Headers de todas as requisições

```http
Message-Id: a1b2c3d4-e5f6-7890-abcd-ef1234567890
Correlation-Id: fedcba09-8765-4321-abcd-ef1234567890
Tenant-Id: PSP-RFB-123456
Timestamp: 2026-07-23T12:00:00Z
Content-Type: application/json
```

- `Message-Id`: único por requisição (idempotência). Gerar novo UUID4 para cada tentativa.
- `Correlation-Id`: mesmo valor para todas as requisições de uma mesma transação.
- `Tenant-Id`: identificar do PSP Recebedor Direto.
- `Timestamp`: RFC 3339, momento de geração da requisição.

## Política de retentativa

| Status code | Retentável | Estratégia |
|---|---|---|
| 200/201/204 | — | Sucesso |
| 400/422/409 | Não | Erro do PSP — corrigir payload |
| 404 | Não | NSU ou recurso inexistente |
| 500 | Sim | Backoff exponencial (1s, 2s, 4s, 8s... max 30s) |
| 503 | Sim | Backoff exponencial + jitter |

## Diferença Informe Preliminar vs Informe de Segregação

| Característica | Preliminar | Segregação |
|---|---|---|
| Caráter | Informativo, não vinculante | Definitivo, vinculante |
| Quando enviar | Logo após pagamento | 2x/dia útil, em lote |
| Granularidade | 1 transação por informe | N transações por informe |
| Gera repasse? | Não | Sim |
| Desencadeia penalidades? | Não | Sim |
| Retentativa | Não necessário (pode reenviar) | Crítico (janela de 3h) |

## Como o simulador usa isso

Nosso `POST /transactions` cria uma transação e, se `PP_BASE_URL` estiver configurada:

1. Mapeia a transação para um `InformeTransacaoIniciada` (arranjo fixo "boleto")
2. Envia para `POST {PP_BASE_URL}/api/v1/boleto`
3. Inclui no response da API o campo `pp_notification` com o resultado

O mock Prism valida o schema mas não executa lógica de negócio (não corrige valores, não mantém estado entre chamadas). Isso é suficiente para provar que o formato dos nossos dados é compatível com o contrato real da PP.

## Fontes
- ✅ Manual de Integração v1.0, seção 4 (contrato de mensagens) e seção 6 (endpoints)
- ✅ Manual de Operações, seção 4 (informes) e seção 4.4 (gatilhos)
- ✅ OpenAPI v0.0.10 (paths e schemas)
