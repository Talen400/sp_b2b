# TASKS.md

Plano de execução. Este arquivo é o que muda entre sessões — ao concluir uma fase ou trocar o foco, edite
aqui, não em `AGENT.md`.

Antes de começar qualquer fase, consulte `PROGRESS.md` para não regenerar o que já existe.

---

## Fase 1: Leitura e Mapeamento Base — ✅ CONCLUÍDA

LC 214/2025, EC 132/2023 e os manuais oficiais do split payment já foram consultados. Estrutura completa
registrada em `REFERENCES.md`. Notas iniciais (`00-Indice`, `01-Visao-Geral-da-Reforma`,
`05-Split-Payment-Mecanismo`, `06-Plataforma-Publica-Documentacao-Tecnica`) já existem — ver
`PROGRESS.md`.

---

## Fase 2: Completar as Notas Base

Gerar as notas que faltam, cada uma seguindo o formato padrão de `AGENT.md`. Caminhos em `DIR.md`.

- `02-Cronograma-2026-2033.md` — linha do tempo ano a ano (2026 teste → 2027 CBS plena → 2029-2032
  transição do IBS → 2033 regime pleno). Fonte: EC 132/2023 + LC 214/2025 + gov.br/fazenda.
- `03-IBS-e-CBS-Basico.md` — o que cada tributo substitui (ICMS/ISS → IBS; PIS/Cofins → CBS), quem
  administra (CGIBS vs RFB), por que dois tributos e não um.
- `04-Nao-Cumulatividade-e-Credito.md` — o conceito mais importante pro nosso simulador B2B: como o
  crédito tributário flui na cadeia de produção. Ligar explicitamente com o que o simulador faz
  (`TASK.md` do projeto Go, TASK-03).
- `07-B2B-vs-B2C.md` — por que o Pagador Efetivo PF não tem visibilidade dos campos fiscais (ver Manual
  de Operações, seção 8, Princípio 1) mas a PJ tem visibilidade ampla (Princípio 2).
- `08-Simples-Nacional-e-MEI.md` — sem mudança em 2026; decisão de permanência/migração até set/2026;
  IBS/CBS "por dentro ou por fora" do DAS.
- `09-Regimes-Especiais-e-Aliquotas-Reduzidas.md` — cesta básica (alíquota zero, Anexo I e XV), redução
  de 60% (saúde, educação, Anexo VII), redução de 30% (profissionais liberais).
- `10-Imposto-Seletivo.md` — o "imposto do pecado" (produtos nocivos à saúde/meio ambiente), separado do
  IBS/CBS, fora do escopo da Fase 1 do split payment (ver Manual de Operações, seção 1.2).
- `11-Glossario.md` — no mínimo 15 termos (IBS, CBS, EC 132, LC 214, CGIBS, RFB, split payment, PSP
  Recebedor Direto/Indireto, Documento Fiscal, crédito tributário, não-cumulatividade, Modelo
  Inteligente/Super Inteligente, Plataforma Pública, alíquota de referência).
- `12-Armadilhas-Comuns-para-Devs.md` — erros típicos: tratar valor monetário como float, confundir
  Informado com Segregado, assumir que todo CNPJ é só numérico (muda em 07/2026), esquecer que
  split ≠ apuração/crédito.
- `13-Fontes-e-Links-Oficiais.md` — espelha `REFERENCES.md`, mas em formato de nota consultável dentro
  do próprio vault (com links clicáveis), não é duplicação burra — é o "cole aqui pra conferir a fonte
  primária" pra quem está lendo as notas, não o `_agent-vault/`.

---

## Fase 3: MOC e Revisão Cruzada

- Atualizar `00-Indice.md` com links para todas as 14 notas, confirmando que nenhum `[[wikilink]]` está
  quebrado.
- Revisar `05-Split-Payment-Mecanismo.md` e `06-Plataforma-Publica-Documentacao-Tecnica.md` (já
  existentes) contra o formato padrão desta versão do `AGENT.md` — foram escritas antes dessa
  padronização, podem não ter todas as 6 seções (TL;DR, contexto, funcionamento, aprofundamento,
  autoavaliação, fontes).

## Definição de "Pronto"

A pasta `vault/` contém:

- `00-Indice.md` navegável, com todos os links funcionando.
- **14 notas-conceito** (ver `DIR.md`).
- `11-Glossario.md` com pelo menos 15 termos.
- Cada nota segue o formato padrão (TL;DR, contexto pra dev, funcionamento, aprofundamento,
  autoavaliação, fontes).
- Toda afirmação técnica tem fonte explicitamente citada (lei, manual oficial, ou marcada como
  simplificação/não verificado).
- Nenhuma regra do nosso simulador Go foi apresentada como se fosse a regra tributária real.

**GO.**
