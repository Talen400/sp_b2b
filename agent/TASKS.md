# TASKS.md

Plano de execução. Muda entre sessões — edite aqui, não em `AGENT.md`. Consulte `PROGRESS.md` antes de
começar qualquer fase.

---

## Fase 0: Pivô de escopo — ✅ CONCLUÍDA (decisão registrada)

Projeto migrou de "CLI simples em memória" para "API REST + SQLite persistente", mantendo a mesma lógica
de domínio (cálculo de split, crédito tributário simplificado) já validada na versão anterior.

---

## Fase 1: Setup do projeto — ✅ CONCLUÍDA

## Fase 2: Camada de domínio — ✅ CONCLUÍDA

## Fase 3: Persistência — ✅ CONCLUÍDA

## Fase 4: HTTP — ✅ CONCLUÍDA

## Fase 5: Seed e demonstração — ✅ CONCLUÍDA

## Fase 6: Polish e documentação — ✅ CONCLUÍDA

---

## Fase 7: Sandbox de integração com a Plataforma Pública (mock local)

Não existe um sandbox público hospedado pelo governo acessível sem credencial de PSP — só o Swagger/
OpenAPI real, publicado pela RFB/CGIBS (`cgibs/openapi-v0_0_10.json`, colocado manualmente pelo
usuário; PDF/DOCX/ZIP gitignorados). O objetivo desta fase é usar esse contrato real pra ver
**integração de verdade** localmente: nosso código fazendo chamadas HTTP reais contra um servidor
que responde exatamente como a Plataforma Pública responderia.

- **Ferramenta de mock**: Prism, da Stoplight — `npx @stoplight/prism-cli mock cgibs/openapi-v0_0_10.json`.
  Dependência **dev-only** (via `npx`, não entra no `go.mod`). Documentar limitação: endpoints de long
  polling com token de posição podem não ser plenamente mockados pelo Prism.
- **`internal/client/pp/`** — pacote novo, cliente HTTP que fala com a Plataforma Pública (mockada em dev).
  Implementa:
  - `SendInformeTransacaoIniciada` — POST /api/v1/{arrj} (arranjo mapeado dinamicamente).
  - `SendInformeSegregacao` — POST /api/v1/segregacao.
  - `ConsultarSplitSuperInteligente` — GET long polling (start → continue → delete).
- **Hook no `POST /transactions`**: se `PP_BASE_URL` estiver setada (flag `-pp-url` ou env var),
  notifica o mock em goroutine fire-and-forget. Response inclui `pp_notification` com status.
- **`make sandbox`** sobe o mock (via Prism) na porta 4010. Documentar no README.
- **Critério de sucesso**: um `curl` real contra o mock com payload construído a partir do
  dicionário de campos do Manual de Integração recebe resposta no formato exato do schema OpenAPI.

Esta fase é **opcional/demonstrativa** — não faz parte da Definição de "Pronto" original.

---

## Definição de "Pronto"

- Todos os 7 endpoints funcionando, com formato de erro padronizado.
- `make build`, `make test`, `make migrate`, `make seed`, `make run` funcionam do zero num checkout limpo.
- Camada de domínio com testes exaustivos; handlers com testes de integração básicos.
- Nenhuma regra de negócio vazando pra dentro de um handler HTTP ou de uma query SQL.
- README permite que alguém sem contexto rode a demonstração completa em menos de 5 comandos.

**GO.**
