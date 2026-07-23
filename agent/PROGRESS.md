# PROGRESS.md

Estado atual do projeto. Leia antes de gerar/alterar qualquer código; atualize ao final de cada rodada.

---

## Status por Fase

- [x] Fase 0 — Pivô de escopo (CLI → API + banco), decisão registrada em `TASKS.md`
- [x] Fase 1 — Setup do projeto
- [x] Fase 2 — Camada de domínio
- [x] Fase 3 — Persistência (SQLite, `modernc.org/sqlite` — puro Go, sem cgo)
- [x] Fase 4 — HTTP (7 endpoints, `net/http` ServeMux nativo)
- [x] Fase 5 — Seed e demonstração (Fazenda → Fábrica → Mercado)
- [x] Fase 6 — Polish e documentação (gofmt, go vet, README, Makefile)
- [ ] Fase 7 — Sandbox de integração com a Plataforma Pública (mock local Prism + client PP)

## Notas de arquitetura já decididas

- **Driver SQLite:** `modernc.org/sqlite v1.54.0` — puro Go, sem cgo, build portável.
- **Porta default:** 8080.
- **Schema:** `companies` (cnpj PK, nome, saldo_credito) e `transactions` (id PK, vendedor_cnpj, comprador_cnpj, valor_bruto, aliquota_ibs, aliquota_cbs, valor_liquido, valor_ibs, valor_cbs, credito_usado, timestamp). Crédito usado e gerado estão na mesma transação.
- **Documentos oficiais:** em `cgibs/` (não `docs-oficiais/`). PDF/DOCX/ZIP gitignorados. `openapi-v0_0_10.json` versionado.
- **Migrations:** embutidas via `//go:embed` em `internal/repository/sqlite/migrations/`, não soltas na raiz.

## Observações da Última Rodada

- Fase 7 iniciada. Cliente PP em `internal/client/pp/`. Hook opcional no `POST /transactions`. Vault reescrito com dados reais dos manuais.
