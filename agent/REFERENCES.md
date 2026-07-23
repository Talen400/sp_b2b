# REFERENCES.md — Contexto de domínio

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

## O split payment de verdade (fonte: Manual de Operações e Manual de Integração, RFB/CGIBS/Serpro, jun/2026)

Isso aqui resume a documentação técnica oficial que temos em mãos — é mais preciso que qualquer coisa
que eu tenha achado via busca na web antes.

**Fase 1 = "B2B Opcional".** O split só se aplica a operações entre PJs (Pagador Original e Recebedor
precisam ter CNPJ). É facultativo pra quem origina a transação (só ativa se preencher os campos de
CBS/IBS), mas obrigatório pros PSPs disponibilizarem.

**Seis arranjos, dois modelos:**
- **Modelo Inteligente** (arranjos iniciados pelo **Pagador**: Pix Estático, TED, TEF) — o valor do
  tributo é o que o Pagador informou, ponto. Sem correção do governo depois.
- **Modelo Super Inteligente** (arranjos iniciados pelo **Recebedor**: Boleto, Pix Dinâmico, Pix
  Automático) — o Recebedor informa o valor ao criar a cobrança, mas a Receita/CGIBS podem corrigir
  esse valor (ou reduzi-lo, se parte do débito já foi extinta) até a baixa da transação, via retorno da
  Plataforma Pública.

**Papéis de PSP:** existe PSP Pagador Direto/Indireto e PSP Recebedor Direto/Indireto. O PSP Recebedor
**Direto** é sempre quem responde regulatoriamente pelo split perante a Plataforma Pública — mesmo que
exista um PSP Indireto na relação comercial com o cliente final.

**Categorias de valor de tributo** (todas em centavos, decimal 18,2): Informado → Corrigido → Em Aberto
→ Segregado (o que efetivamente vira Repasse Financeiro) → Aplicado (o que é mostrado ao Pagador).

**Fluxo de comunicação com o governo:** Informe de Transação Iniciada → (Retorno Super Inteligente,
quando aplicável) → Informe Preliminar de Pagamento (informativo) → Informe de Segregação (definitivo,
2x/dia útil, gera obrigação de Repasse Financeiro) → Repasse Financeiro em D+N via eventos TES (CBS) e
STR (IBS).

**Simplificações que nosso simulador ainda assume** (ver TASK.md pra escopo real do MVP):
- Não simulamos a Plataforma Pública nem o Super Inteligente — nosso simulador é puramente Modelo
  Inteligente simplificado (valor informado = valor segregado).
- Sem Documento Fiscal, sem CNPJ alfanumérico (o real usa formato novo a partir de 07/2026, IN RFB
  2.229/2024), sem distinção de PSP Direto/Indireto.
- Nosso "crédito tributário" é uma simplificação didática — o mecanismo real de crédito B2B não faz
  parte do escopo do Manual de Operações do split payment (que trata só da segregação e repasse, não da
  apuração/crédito, que é assunto de outro normativo).

## Documentos oficiais que temos como referência primária
- **Manual de Operações — Split Payment** (RFB/Serpro/CGIBS/Procergs, versão preliminar, jun/2026):
  fluxos de negócio, os 6 arranjos, modelos Inteligente/Super Inteligente, responsabilidades dos PSPs,
  MOC (Mecanismo de Ocorrências).
- **Manual de Integração — Plataforma Pública de Split Payment** (v1.0): dicionário de campos, payloads,
  endpoints REST por arranjo, política de tratamento de erros (RFC 7807), headers padrão.
- **OpenAPI v0.0.10**: especificação machine-readable dos endpoints (formato JSON/Swagger).

Esses três documentos já estão nos uploads do projeto — se o agente precisar de detalhe fino sobre um
campo, endpoint ou regra de negócio, a fonte é esses arquivos, não a internet.

## Simplificações que estamos assumindo na simulação

## Glossário rápido
- **IBS**: Imposto sobre Bens e Serviços (substitui ICMS/ISS).
- **CBS**: Contribuição sobre Bens e Serviços (substitui PIS/Cofins).
- **CGIBS**: Comitê Gestor do IBS.
- **RFB**: Receita Federal do Brasil.
- **Split payment**: segregação automática do tributo no momento do pagamento.

## Se quisermos ficar mais realistas depois
A documentação técnica oficial (Manual de Integração + Swagger da Plataforma
Pública do Split Payment) foi publicada pela RFB/CGIBS em jun/2026. Já está anexada
ao projeto (`docs-oficiais/`) e é referência primária pro vault (`_agent-vault/`) —
mas continua fora do escopo do que a API deste repositório implementa (a API é uma
simulação local, não integra com a Plataforma Pública de verdade).

---

## Referências de engenharia (API + banco de dados)

Usadas em conjunto com as referências de domínio tributário acima — estas são sobre *como construir* o
software, não sobre as regras fiscais em si. Ver `AGENT.md`/`TASKS.md` pro escopo travado do projeto.

- **Go standard library** (`net/http`, `database/sql`, `encoding/json`) — https://pkg.go.dev/std —
  fonte primária pra tudo que não for regra de negócio nem SQL específico do driver escolhido.
- **ServeMux com roteamento por método+path** (Go 1.22+) — https://pkg.go.dev/net/http#ServeMux —
  justifica não precisar de router externo (chi, gorilla/mux) pro escopo deste projeto.
- **Driver SQLite**: `modernc.org/sqlite` — driver puro-Go (sem cgo), mantém `make build` simples e
  portátil sem depender de gcc/libsqlite3 no sistema.
  ✅ `Fonte: https://pkg.go.dev/modernc.org/sqlite`
- **Convenções REST/HTTP status codes** — RFC 9110 — https://www.rfc-editor.org/rfc/rfc9110.html
- **Testes com `net/http/httptest`** — https://pkg.go.dev/net/http/httptest
