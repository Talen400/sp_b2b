# PROGRESS.md

Estado atual do vault. O agente deve ler este arquivo **antes** de gerar qualquer nota (para não
duplicar/reprocessar) e **atualizá-lo** ao final de cada rodada.

---

## Status por Fase

- [x] Fase 1 — Leitura e Mapeamento Base
- [x] Fase 2 — Notas Base (14/14)
- [x] Fase 3 — MOC e Revisão Cruzada

## Notas Concluídas

- [x] `00-Indice.md`
- [x] `01-Visao-Geral-da-Reforma.md`
- [x] `02-Cronograma-2026-2033.md`
- [x] `03-IBS-e-CBS-Basico.md`
- [x] `04-Nao-Cumulatividade-e-Credito.md`
- [x] `05-Split-Payment-Mecanismo.md`
- [x] `06-Plataforma-Publica-Documentacao-Tecnica.md`
- [x] `07-B2B-vs-B2C.md`
- [x] `08-Simples-Nacional-e-MEI.md`
- [x] `09-Regimes-Especiais-e-Aliquotas-Reduzidas.md`
- [x] `10-Imposto-Seletivo.md`
- [x] `11-Glossario.md`
- [x] `12-Armadilhas-Comuns-para-Devs.md`
- [x] `13-Fontes-e-Links-Oficiais.md`

## Marcadas como `⚠️ não verificado` (revisar depois)

- **04:** crédito financeiro vs físico — confirmação de que o Brasil adota majoritariamente crédito financeiro. Verificar em fonte oficial.
- **08:** Opções A/B/C para SIMPLES no split — decisão pendente até set/2026. Atualizar após a definição.
- **09:** Tabela machine-readable de alíquotas — depende de regulamentação futura do CGIBS/RFB.
- **10:** Lista de produtos sujeitos ao IS pode ser ampliada por lei ordinária.
- **12:** CNPJ alfanumérico — confirmar na instrução normativa da RFB (mai/2026).

## Observações da Última Rodada

- Estrutura reorganizada: `agent/` (controle) + `vault/` (conteúdo). `DIR.md` atualizado.
- 4 notas que existiam na estrutura anterior foram recriadas (00, 01, 05, 06).
- 10 notas novas da Fase 2 criadas (02, 03, 04, 07, 08, 09, 10, 11, 12, 13).
- Todas as notas seguem o formato padrão (TL;DR / contexto dev / funcionamento / callout / autoavaliação / fontes).
- `00-Indice.md` já contém wikilinks para todas as 14 notas.
- Notas 05 e 06 foram escritas já no novo formato — nenhuma revisão adicional necessária.
