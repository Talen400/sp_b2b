# REFERENCES.md — Contexto de domínio + engenharia

## O que é o split payment (resumo)
Mecanismo da Reforma Tributária (Emenda Constitucional nº 132/2023 e Lei
Complementar nº 214/2025) em que, no momento do pagamento de uma compra, o valor
correspondente aos tributos sobre consumo — **IBS** (Imposto sobre Bens e Serviços,
estadual/municipal) e **CBS** (Contribuição sobre Bens e Serviços, federal) — é
segregado automaticamente pelo sistema financeiro. O vendedor recebe só o valor
líquido; o imposto vai direto para o fisco na liquidação, sem depender de
recolhimento posterior pela empresa.

## Linha do tempo real (para contexto, não é o que estamos implementando)
- 2023: EC 132/2023 cria a base constitucional da reforma.
- 2025: LC 214/2025 regulamenta o modelo.
- Mai/2026: Ato Conjunto RFB/CGIBS nº 02 aprova o Manual de Integração e o Swagger
  da Plataforma Pública do Split Payment.
- 2026: ano-teste — CBS com alíquota reduzida e caráter informativo.
- 2027: cobrança real via split payment, começando por Pix e transferências a débito.

## B2B: por que o crédito tributário importa
No IBS/CBS (assim como no modelo de IVA), o imposto é **não-cumulativo**: quando
uma empresa compra insumo/serviço de outra, o IBS+CBS pago naquela compra vira
**crédito** que ela pode abater do imposto devido quando ela mesma vender. É isso
que evita "imposto sobre imposto" ao longo da cadeia produtiva. No B2C (venda pro
consumidor final) esse crédito não existe pro comprador, porque ele é quem
efetivamente arca com o tributo — é aí que a simulação B2B difere da simulação
B2C original.

## Documentos oficiais que temos como referência primária

Todos em `cgibs/`:
- **Manual de Operações — Split Payment** (RFB/Serpro/CGIBS/Procergs, versão preliminar, jun/2026):
  fluxos de negócio, os 6 arranjos, modelos Inteligente/Super Inteligente, responsabilidades dos PSPs,
  MOC, ciclo de vida dos informes, repasse financeiro.
- **Manual de Integração — Plataforma Pública de Split Payment** (v1.0): dicionário de campos, payloads,
  endpoints REST por arranjo, política de tratamento de erros (RFC 7807), headers padrão, mTLS/OAuth.
- **OpenAPI v0.0.10**: especificação machine-readable dos endpoints (formato JSON/Swagger, 201K).

Os PDFs foram extraídos para `/tmp/manual-integracao.txt` (3773 linhas) e `/tmp/manual-operacoes.txt`
(2728 linhas). As seções abaixo resumem os pontos mais relevantes extraídos dos PDFs reais.

### Headers HTTP obrigatórios (Manual de Integração, seção 4)
- `Message-Id`: UUID4, obrigatório, idempotência.
- `Correlation-Id`: UUID4, rastreio ponta a ponta.
- `Tenant-Id`: identificador do PSP na PP.
- `Timestamp`: RFC 3339.

### Formato Decimal(18,2)
- 18 dígitos no total, 2 decimais.
- No payload JSON: string, ex: `"1234567890123456.78"`.
- Em nosso Go: centavos `int64`, convertido na serialização.
- Valores de tributo: `vlIbs`, `vlCbs`, `vlInf`, `vlCbsCorr`, `vlIbsSeg`, etc.
- Duas casas decimais, sem arredondamento bancário.

### Os 6 arranjos
| Arranjo | Modelo | Iniciado por | Informe Transação |
|---|---|---|---|
| Boleto | Super Inteligente | Recebedor | Obrigatório |
| Pix Dinâmico | Super Inteligente | Recebedor | Obrigatório |
| Pix Automático | Super Inteligente | Recebedor | Obrigatório |
| Pix Estático | Inteligente | Pagador | N/A |
| TED | Inteligente | Pagador | N/A |
| TEF | Inteligente | Pagador | N/A |

### 5 categorias de valor de tributo
1. **Valor Informado** (`vlInfo`): o que consta no Documento Fiscal.
2. **Valor Corrigido** (`vlCorr`): correção do Super Inteligente (se aplicável).
3. **Valor Em Aberto** (`vlAberto`): parte do débito ainda não liquidada.
4. **Valor Segregado** (`vlSeg`): o que efetivamente foi separado no pagamento — base do Repasse Financeiro.
5. **Valor Aplicado** (`vlApl`): o que é mostrado ao Pagador no extrato.

### Política de erros (Manual de Integração, seção 5)
- Formato: RFC 7807, media type `application/problem+json`.
- Campos: `type` (URI), `title`, `status`, `detail`, `instance`, `retornabilidade`.
- Erros de validação: 400/422 com `type` específico.
- Erros internos: 500, sem detalhes no response.
- `retornabilidade`: booleano indicando se o PSP pode retentar a requisição.

### Long polling / Retorno Super Inteligente (Manual de Integração, seção 3.6)
- Modelo pull-based: PSP consulta PP por novas mensagens.
- Token de posição trocado em cada response (header `proximoToken`).
- Fluxo: `GET /start` → recebe mensagens + token → `GET /{token}` → ... → `DELETE /{token}`.
- Consulta retroativa: `/api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start?fromNsu=X&toNsu=Y`.
- Timeout configurável na PP (long polling HTTP).

## Simplificações que estamos assumindo na simulação
- Não simulamos a Plataforma Pública nem o Super Inteligente em modo real — o mock Prism valida
  o contrato, mas não executa lógica de negócio.
- Nosso "crédito tributário" é uma simplificação didática — o mecanismo real de crédito B2B não faz
  parte do escopo do Manual de Operações do split payment.
- Arranjo fixo  "boleto" para o hook do client PP (configurável via código ou env).
- Sem mTLS, OAuth, ou certificado digital — o mock Prism roda sem autenticação.
- CNPJ tratado como string livre, sem validação de dígitos ou formato alfanumérico.

## Referências de engenharia (API + banco de dados)

- **Go standard library** (`net/http`, `database/sql`, `encoding/json`) — https://pkg.go.dev/std
- **ServeMux com roteamento por método+path** (Go 1.22+) — https://pkg.go.dev/net/http#ServeMux
  Justifica não precisar de router externo.
- **Driver SQLite**: `modernc.org/sqlite v1.54.0` — puro Go (sem cgo), build portável sem GCC,
  ativamente mantido, implementação pura do protocolo SQLite. Justificativa: manter `make build`
  simples em qualquer plataforma, sem dependência de toolchain C.
- **Convenções REST/HTTP status codes** — RFC 9110 — https://www.rfc-editor.org/rfc/rfc9110.html
- **Testes com `net/http/httptest`** — https://pkg.go.dev/net/http/httptest
- **Prism (Stoplight)** — https://github.com/stoplightio/prism — mock server gerado a partir de OpenAPI.
  Dependência **dev-only** (via `npx`, não entra no `go.mod` nem no binário).

## Glossário rápido
- **IBS**: Imposto sobre Bens e Serviços (substitui ICMS/ISS).
- **CBS**: Contribuição sobre Bens e Serviços (substitui PIS/Cofins).
- **CGIBS**: Comitê Gestor do IBS.
- **RFB**: Receita Federal do Brasil.
- **Split payment**: segregação automática do tributo no momento do pagamento.
- **PP**: Plataforma Pública do Split Payment.
- **PSP**: Instituição de Pagamento (ex: bancos, fintechs).
- **MOC**: Mecanismo de Ocorrências — registro de falhas no split.
