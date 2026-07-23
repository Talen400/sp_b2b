# Fontes e Links Oficiais

Este arquivo indexa as seções relevantes dos documentos oficiais que temos em `cgibs/`. Use como referência rápida: cada seção listada contém detalhe normativo ou técnico usado nas notas do vault e no código.

---

## Manual de Integração — Plataforma Pública de Split Payment (v1.0)
`cgibs/03153733-manual-de-integracao-plataforma-publica-de-split-payment-v1.pdf`

### Seção 1 — Introdução
- Propósito da Plataforma Pública
- Abrangência (Fase 1: B2B opcional)

### Seção 2 — Visão geral da arquitetura
- Diagrama de blocos: PSP Pagador → CIP/STR → PSP Recebedor → PP
- Papéis: PSP Pagador Direto/Indireto, PSP Recebedor Direto/Indireto

### Seção 3 — Segurança
- 3.1 mTLS (certificado ICP-Brasil)
- 3.2 OAuth 2.0 Client Credentials
- 3.3 Ambientes (produção vs homologação)

### Seção 4 — Contrato de mensagens (payloads, headers, dicionário)
- 4.1 Headers obrigatórios: `Message-Id` (UUID4), `Correlation-Id` (UUID4), `Tenant-Id`, `Timestamp` (RFC 3339)
- 4.2 Dicionário de campos com tipo, formato, obrigatoriedade
  - Decimal(18,2): 18 dígitos, 2 casas, string no JSON
  - Tabela completa de campos por informe
- 4.3 Exemplos de payload por arranjo

### Seção 5 — Política de tratamento de erros
- 5.1 Formato RFC 7807 (`application/problem+json`)
  - Campos: `type` (URI), `title`, `status`, `detail`, `instance`
- 5.2 Tabela de erros: código, HTTP status, retornabilidade
- 5.3 Boas práticas: backoff exponencial, circuit breaker
- 5.4 Exemplos de response de erro

### Seção 6 — Endpoints por arranjo
- Tabela completa: path, método, descrição, request/response schema

---

## Manual de Operações — Split Payment
`cgibs/30145925-minuta-split-payment-manual-de-operacoes.pdf`

### Seção 1 — Introdução e conceitos
- 1.1 O que é split payment
- 1.2 Atores: Pagador, Recebedor, PSP Pagador (Direto/Indireto), PSP Recebedor (Direto/Indireto)
- 1.3 Arranjos de pagamento cobertos (6 arranjos)

### Seção 2 — Modelos de split
- 2.1 **Modelo Inteligente** (arranjos iniciados pelo Pagador: Pix Estático, TED, TEF)
- 2.2 **Modelo Super Inteligente** (arranjos iniciados pelo Recebedor: Boleto, Pix Dinâmico, Pix Automático)
- 2.3 Tabela comparativa: quando cada modelo se aplica

### Seção 3 — Ciclo de vida da transação com split
- 3.1 As 5 categorias de valor: Informado → Corrigido → Em Aberto → Segregado → Aplicado
- 3.2 Fluxo de liquidação financeira
- 3.3 Papel do Documento Fiscal

### Seção 4 — Informes (comunicação com a PP)
- 4.1 Informes em arranjos iniciados pelo Recebedor
  - Informe de Transação Iniciada (obrigatório: Boleto, Pix Dinâmico, Pix Automático)
  - Informe de Transação Atualizada (obrigatório: Boleto, Pix Dinâmico)
  - Retorno Super Inteligente: Correção e Em Aberto
  - Informe de Baixa (exceto por pagamento)
  - Informe Preliminar de Pagamento
  - Informe de Segregação
- 4.2 Informes em arranjos iniciados pelo Pagador
  - Informe Preliminar de Pagamento (opcional na TEF)
  - Informe de Segregação
- 4.3 Diferença entre Informe Preliminar e Informe de Segregação (tabela comparativa)
- 4.4 Gatilhos de envio por arranjo (tabela completa)

### Seção 5 — Prazos e janelas (visão geral; detalhe no Manual de Tempos)
- Janelas de consolidação
- D+N para repasse financeiro

### Seção 6 — Repasse Financeiro
- 6.1 Obrigação de Repasse do PSP Recebedor Direto
- 6.2 CBS → TES (Sistema de Transferência Eletrônica de Segregação)
- 6.3 IBS → STR (Sistema de Transferência de Reservas)
- 6.4 Periodicidade: D+N (dia útil + N dias)

### Seção 7 — MOC (Mecanismo de Ocorrências)
- 7.1 O que é uma ocorrência
- 7.2 Quando registrar
- 7.3 Fluxo de consulta e registro

### Seção 8 — Disposições finais
- Cronograma de implantação
- Penalidades (vinculadas ao Manual de Tempos)

---

## Manual de Tempos — Split Payment
`cgibs/20150141-20260715-fin-split-payment-manual-de-tempos-minuta.docx`
> ⚠️ Minuta. Ainda não processada. Conteúdo: prazos, janelas de consolidação, regras de penalidade.

---

## OpenAPI Specification v0.0.10
`cgibs/openapi-v0_0_10.json` (versionado)

### Paths
- `POST /api/v1/boleto` — Informe de Transação Iniciada (Boleto)
- `PUT /api/v1/boleto/{nsuId}` — Informe de Transação Atualizada (Boleto)
- `POST /api/v1/pix-dinamico` — Informe de Transação Iniciada (Pix Dinâmico)
- `PUT /api/v1/pix-dinamico/{nsuId}` — Informe de Transação Atualizada (Pix Dinâmico)
- `POST /api/v1/pix-automatico` — Informe de Transação Iniciada (Pix Automático)
- `POST /api/v1/pix-estatico` — Informe Preliminar de Pagamento (Pix Estático)
- `POST /api/v1/ted` — Informe Preliminar de Pagamento (TED)
- `POST /api/v1/tef` — Informe Preliminar de Pagamento (TEF)
- `POST /api/v1/segregacao` — Informe de Segregação (todos arranjos)
- `GET /api/v1/{arrj}/{idPsp}/tributos/stream/start` — Long polling start
- `GET /api/v1/{arrj}/{idPsp}/tributos/stream/{token}` — Long polling continue
- `DELETE /api/v1/{arrj}/{idPsp}/tributos/stream/{token}` — Long polling end
- `GET /api/v1/retroativo/{arrj}/{idPsp}/tributos/stream/start` — Consulta retroativa

### Schemas globais
- ErrorResponse (RFC 7807)
- Headers (Message-Id, Correlation-Id, Tenant-Id, Timestamp)

---

## Legislação
- **EC 132/2023:** Emenda Constitucional da Reforma Tributária
- **LC 214/2025:** Regulamentação do IBS/CBS
- **LC 227/2026:** Instituição do CGIBS
- **Ato Conjunto RFB/CGIBS nº 02/2026:** Aprovação do Manual de Integração e Swagger da PP
- **IN RFB 2.229/2024:** CNPJ alfanumérico (vigência jul/2026)

---

## Links externos
> ⚠️ Nenhum ambiente de teste público acessível sem credencial de PSP.
> O sandbox local usa Prism a partir do OpenAPI oficial (`cgibs/openapi-v0_0_10.json`).
