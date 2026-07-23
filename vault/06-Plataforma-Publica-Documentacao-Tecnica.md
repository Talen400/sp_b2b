# Plataforma Pública — Documentação Técnica

## TL;DR
A **Plataforma Pública do Split Payment** (RFB/CGIBS, jun/2026) é a API REST que
PSPs consultam para obter os valores segregados (split) de uma transação. A
documentação inclui um **Manual de Integração** (dicionário de campos, endpoints,
política de erros) e uma **especificação OpenAPI** (machine-readable).

## Por que isso importa pra quem programa
Se você for integrar um sistema de pagamento ou ERP com o split payment real, a
Plataforma Pública é o ponto de contato. Entender os endpoints, headers e formato
de erros (RFC 7807) é requisito para qualquer implementação.

## Funcionamento

### Endpoints principais (OpenAPI v0.0.10)
- **Consulta de split:** informa os valores segregados para um Documento Fiscal.
- **Registro de liquidação:** PSP notifica a plataforma de que a liquidação ocorreu.
- **MOC (Mecanismo de Ocorrências):** consulta/registro de falhas no split.
✅ `Fonte: Manual de Integração — Plataforma Pública v1.0, seção 4`

### Headers padrão
Toda requisição deve incluir:
- `X-Idempotency-Key` — chave de idempotência (UUID).
- `X-Correlation-Id` — ID de correlação para rastreio.
- `Content-Type: application/json`
✅ `Fonte: Manual de Integração, seção 4.1`

### Política de erros
A API usa o formato RFC 7807 (Problem Details):
```json
{
  "type": "https://split.rfb.gov.br/erros/campo-invalido",
  "title": "Campo inválido",
  "status": 422,
  "detail": "O campo 'valorBruto' não pode ser negativo",
  "instance": "/split/v1/consultar"
}
```
✅ `Fonte: Manual de Integração, seção 5`

### Autenticação
A comunicação usa certificado digital (ICP-Brasil) para autenticação mútua
(mTLS). A plataforma também exige autorização via OAuth 2.0 (client credentials).
✅ `Fonte: Manual de Integração, seção 3`

> [!note]- 🔍 Aprofundando: Diferença entre os manuais
> O **Manual de Operações** descreve o fluxo de negócio (como o split funciona,
> responsabilidades dos PSPs). O **Manual de Integração** descreve como chamar a
> API. Ambos foram aprovados pelo Ato Conjunto RFB/CGIBS nº 02/2026 e são
> complementares.
> ⚠️ A versão atual (v1.0, jun/2026) é preliminar — pode mudar até a obrigatoriedade
> plena em 2027.

## Perguntas de autoavaliação
1. Quais são os dois documentos oficiais que descrevem a Plataforma Pública?
2. Qual formato de erro a API usa?
3. Quais mecanismos de segurança a plataforma exige?

## Fontes
- ✅ Manual de Integração — Plataforma Pública de Split Payment v1.0, seções 3, 4, 4.1, 5
- ✅ Ato Conjunto RFB/CGIBS nº 02, de 27/05/2026
