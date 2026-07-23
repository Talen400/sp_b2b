# AGENT.md — Split Payment API (Go + banco de dados)

## Localização dos arquivos
`AGENT.md`, `DIR.md`, `REFERENCES.md`, `TASKS.md`, `PROGRESS.md` ficam na raiz do repo, ao lado do
código (`cmd/`, `internal/`). Não confundir com `_agent-vault/` (controle do vault de conhecimento
tributário — projeto irmão, documentação, sem código).

Esses arquivos substituem os antigos `AGENT.md`/`TASK.md` da versão CLI do simulador — o projeto pivotou
de "CLI em memória" para "API REST + banco de dados persistente". `TASK.md` foi renomeado `TASKS.md` pra
acompanhar o padrão do vault (arquivo mutável, ao contrário deste, que é regra permanente).

## Papel
Você é um agente de engenharia trabalhando sozinho em sandbox. Constrói uma **API REST em Go**, com
persistência em banco de dados, que simula o fluxo de split payment B2B (segregação de IBS/CBS +
crédito tributário simplificado). Vai para portfólio — o código e os testes importam tanto quanto a
funcionalidade.

Este arquivo contém as regras **permanentes**. O que fazer *agora* vai em `TASKS.md`. Onde as coisas
ficam, `DIR.md`. Fontes de domínio (tributário) e de engenharia, `REFERENCES.md`. O que já foi feito,
`PROGRESS.md`.

---

## Regras Gerais (sempre válidas)

1. **Leia `TASKS.md` e `PROGRESS.md` antes de tocar em código.** Não regenere o que já está marcado como
   concluído sem motivo explícito. Ao terminar uma etapa, atualize `PROGRESS.md`.

2. **Escopo travado, no estilo "subject".** A seção "Requisitos Obrigatórios" abaixo é fixa — não
   adicione autenticação, múltiplos usuários, UI web, filas assíncronas ou qualquer coisa fora dela sem
   que isso primeiro seja escrito em `TASKS.md`. Preferível fazer pouco e bem testado do que muito e
   frágil.

3. **Dependências externas: mínimas e justificadas.** Cada dependência fora da standard library do Go
   precisa de uma linha em `REFERENCES.md` explicando por que ela é necessária (a stdlib não cobre) —
   nunca adicionar framework HTTP (gin/echo/fiber) ou ORM (gorm) só por conveniência. Usar
   `net/http` (ServeMux nativo, Go 1.22+) e `database/sql` puro. Driver de banco é a única exceção
   aceitável de cara (a stdlib não inclui driver de SQL nenhum).

4. **Separação em camadas estrita, sem vazamento.**
   - `internal/domain/` — regras de negócio puras (cálculo de split, lógica de crédito). **Zero** import
     de `net/http`, `database/sql` ou qualquer coisa de infraestrutura. Se o domínio importa infra, é bug
     de arquitetura, não detalhe.
   - `internal/repository/` — implementa persistência. Domínio não conhece SQL; repository não conhece
     HTTP.
   - `internal/handler/` — HTTP: decodifica request, chama domínio/repository, codifica response. Sem
     regra de negócio aqui.

5. **Todo endpoint tem contrato de erro padronizado.** Nunca vazar erro interno cru (`err.Error()`) pro
   cliente. Formato de erro JSON único, definido em `TASKS.md`, usado em toda a API.

6. **Teste a lógica de domínio antes de expor no HTTP.** Toda função em `internal/domain/` tem
   `_test.go` cobrindo caso normal, valor zero, e caso de erro/limite. Handlers HTTP podem ter testes de
   integração mais leves (happy path + 1-2 erros), não precisam da mesma exaustão do domínio.

7. **Banco de dados: migrations versionadas, nunca `AutoMigrate` mágico.** Arquivos SQL numerados em
   `migrations/`, aplicados na inicialização ou via comando explícito — nunca criação de schema
   implícita dentro do código Go sem arquivo rastreável.

8. **Valores monetários em centavos (`int64`), nunca `float64`.** Isso já valia no simulador CLI e
   continua valendo — inclusive no schema do banco (coluna `INTEGER`, não `REAL`/`DECIMAL` de ponto
   flutuante).

9. **`Makefile` é a interface oficial do projeto.** Não documentar comandos `go run`/`go build` soltos no
   README como forma primária de uso — o Makefile expõe os targets padronizados (ver `DIR.md`). Isso é
   estilo 42: quem chega no projeto não deveria precisar adivinhar o comando certo.

10. **Cite a fonte de qualquer regra de domínio tributário que aparecer em comentário/doc do código**
    (`// Fonte: Manual de Operações — Split Payment, seção X`), do mesmo jeito que as notas do vault
    fazem. Código sem fonte pra uma regra fiscal específica é `⚠️ não verificado` — documente como tal.

---

## Requisitos Obrigatórios (escopo travado — "subject")

### Domínio
- `Company` (CNPJ, nome, saldo de crédito em centavos).
- `Transaction` (id, CNPJ vendedor, CNPJ comprador, valor bruto, alíquotas de IBS/CBS, timestamp,
  resultado do split, crédito usado/gerado).
- `CalculateSplit` — função pura, testada exaustivamente (ver regras do simulador original em
  `PROGRESS.md`/histórico — a lógica de cálculo não muda, só ganha persistência).

### Persistência
- SQLite (arquivo local — sem servidor de banco separado, mantém o projeto rodável com um único
  binário + um arquivo `.db`). Migrations em `migrations/0001_init.sql` etc.

### Endpoints mínimos (todos em `/api/v1`)
- `POST /companies` — cria empresa.
- `GET /companies/{cnpj}` — detalhe + saldo de crédito.
- `GET /companies` — lista.
- `POST /transactions` — simula uma venda entre duas empresas (calcula e persiste o split + atualiza
  crédito).
- `GET /transactions` — histórico, com filtro opcional por CNPJ.
- `GET /transactions/{id}` — detalhe de uma transação.
- `GET /healthz` — liveness check simples (obrigatório em qualquer API "de verdade").

### Qualidade
- `go vet` e `gofmt -l` limpos.
- Testes de domínio com cobertura das funções de cálculo/crédito.
- Cenário fictício da demonstração (3 empresas, 2 vendas — ver histórico do projeto) reproduzido via
  script de seed, não hardcoded no `main.go`.

## Fora de Escopo (não implementar sem atualizar TASKS.md primeiro)
- Autenticação/autorização.
- Múltiplos tenants/usuários.
- Frontend/UI.
- Integração real com a Plataforma Pública (RFB/CGIBS) — continua sendo simulação local.
- Filas, workers assíncronos, cache distribuído.
