# Plataforma Pública — Documentação Técnica

## TL;DR
A **Plataforma Pública do Split Payment** (RFB/CGIBS/Serpro) é a API REST que os PSPs consultam para comunicar transações com split de IBS/CBS. Documentada pelo Manual de Integração v1.0 (jun/2026) + especificação OpenAPI v0.0.10. A comunicação é síncrona (informes) e assíncrona (long polling para retorno do Super Inteligente).

## Headers HTTP obrigatórios
Toda requisição para a PP deve incluir:

| Header | Formato | Obrigatório | Descrição |
|---|---|---|---|
| `Message-Id` | UUID4 | Sim | Idempotência — mesma Message-Id em retentativas é ignorada |
| `Correlation-Id` | UUID4 | Sim | Rastreio ponta a ponta (mesmo valor em toda a cadeia de uma transação) |
| `Tenant-Id` | string | Sim | Identificador do PSP na PP |
| `Timestamp` | RFC 3339 | Sim | Momento da geração da requisição |

- `Content-Type: application/json` (implícito, mas deve ser enviado)

## Formato Decimal(18,2)
- 18 dígitos no total, 2 decimais.
- No payload JSON: representado como **string** (ex: `"1234567890123456.78"`).
- Em centavos `int64` no código Go — a conversão para string com 2 casas decimais ocorre na serialização.
- Aplica-se a todos os campos monetários: `vlInf`, `vlIbs`, `vlCbs`, `vlIbsSeg`, `vlCbsSeg`, `vlLiq`, etc.
- Duas casas decimais, arredondamento padrão (não bancário).

## Endpoints por arranjo (OpenAPI v0.0.10)

| Arranjo | Path | Informe |
|---|---|---|
| Boleto | `POST /api/v1/boleto` | Transação Iniciada |
| Boleto | `PUT /api/v1/boleto/{nsuId}` | Transação Atualizada |
| Pix Dinâmico | `POST /api/v1/pix-dinamico` | Transação Iniciada |
| Pix Dinâmico | `PUT /api/v1/pix-dinamico/{nsuId}` | Transação Atualizada |
| Pix Automático | `POST /api/v1/pix-automatico` | Transação Iniciada |
| Pix Estático | `POST /api/v1/pix-estatico` | Informe Preliminar de Pagamento |
| TED | `POST /api/v1/ted` | Informe Preliminar de Pagamento |
| TEF | `POST /api/v1/tef` | Informe Preliminar de Pagamento |
| Todos | `POST /api/v1/segregacao` | Informe de Segregação (lote) |

Além destes, endpoints de long polling para Retorno Super Inteligente:
- `GET /api/v1/{arrj}/{idPsp}/tributos/stream/start` — inicia stream
- `GET /api/v1/{arrj}/{idPsp}/tributos/stream/{token}` — continua stream
- `DELETE /api/v1/{arrj}/{idPsp}/tributos/stream/{token}` — encerra stream
- `GET /api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start?fromNsu=X&toNsu=Y` — consulta retroativa

## As 5 categorias de valor de tributo

1. **Valor Informado** (`vlInf`): o que consta no Documento Fiscal criado pelo Recebedor.
2. **Valor Corrigido** (`vlCbsCorr`, `vlIbsCorr`): correção enviada pelo Super Inteligente (se aplicável), quando o governo recalcula o tributo devido.
3. **Valor Em Aberto** (`vlCbsAberto`, `vlIbsAberto`): parcela do débito que ainda não foi liquidada (ex: boleto pago parcialmente).
4. **Valor Segregado** (`vlIbsSeg`, `vlCbsSeg`): o que efetivamente foi separado no momento da liquidação — **base do Repasse Financeiro**.
5. **Valor Aplicado** (`vlIbsApl`, `vlCbsApl`): o valor exibido ao Pagador no extrato/fatura.

## Política de erros (RFC 7807)
A API usa `application/problem+json` com o schema:

```json
{
  "type": "https://split.rfb.gov.br/erros/campo-invalido",
  "title": "Campo inválido",
  "status": 422,
  "detail": "O campo 'valorBruto' não pode ser negativo",
  "instance": "/api/v1/boleto"
}
```

Campos adicionais no nosso contexto:
- `retornabilidade`: `true` se o PSP pode retentar a requisição (ex: timeout, 503); `false` se o erro é definitivo (ex: schema inválido).

Tabela resumo de erros comuns:

| Tipo | HTTP | Retentável | Cenário |
|---|---|---|---|
| `erro-validacao-schema` | 400 | Não | JSON mal formado, campo obrigatório ausente |
| `erro-campo-invalido` | 422 | Não | Valor fora do domínio permitido |
| `erro-nao-encontrado` | 404 | Não | NSU inexistente |
| `erro-conflito` | 409 | Não | Message-Id duplicada com conteúdo diferente |
| `erro-interno` | 500 | Sim | Falha interna da PP |
| `erro-indisponivel` | 503 | Sim | PP temporariamente indisponível |

## Autenticação e segurança
- **mTLS (Mutual TLS):** certificado digital ICP-Brasil para autenticação mútua.
- **OAuth 2.0 Client Credentials:** token de acesso no header `Authorization: Bearer`.
- **Produção vs homologação:** ambientes distintos, cada um com seu certificado e credenciais.
- Nosso mock Prism não exige autenticação — usamos apenas para validar o contrato dos payloads.

## Diferença entre os manuais
- **Manual de Operações:** descreve o fluxo de negócio (responsabilidades dos PSPs, gatilhos de envio, MOC, repasse financeiro). Lê primeiro para entender *o que* comunicar.
- **Manual de Integração:** descreve *como* chamar a API (headers, payloads, erros, mTLS). Consulta para implementar o cliente HTTP.
- **OpenAPI v0.0.10:** a especificação machine-readable. Fonte do mock Prism e das structs do client Go.

## Fluxo longo polling / Retorno Super Inteligente
Nos arranjos Super Inteligente (Boleto, Pix Dinâmico, Pix Automático), o governo pode corrigir os valores de tributo após o Informe de Transação Iniciada. O PSP consulta essas correções via long polling:

1. PSP envia `GET /api/v1/{arrj}/{idPsp}/tributos/stream/start`
2. PP mantém conexão aberta (timeout configurável) ou responde imediatamente se houver mensagens
3. PP responde com mensagens + header `proximoToken` (URL da próxima consulta)
4. PSP faz `GET /api/v1/{arrj}/{idPsp}/tributos/stream/{token}`
5. Repete até receber `204 No Content` (sem mais mensagens)
6. PSP envia `DELETE /api/v1/{arrj}/{idPsp}/tributos/stream/{token}` para encerrar

Token de posição garante ordem de leitura e permite reprocessamento.

## Fontes
- ✅ Manual de Integração — Plataforma Pública de Split Payment v1.0 (cgibs/ PDF)
- ✅ Manual de Operações — Split Payment (cgibs/ PDF)
- ✅ OpenAPI v0.0.10 (cgibs/openapi-v0_0_10.json)
- ✅ Ato Conjunto RFB/CGIBS nº 02, de 27/05/2026
