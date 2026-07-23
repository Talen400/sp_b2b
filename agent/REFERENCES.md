# REFERENCES.md

Fontes de verificação, nesta ordem de prioridade. Toda afirmação técnica nas notas deve rastrear até uma
destas (ver tags de fonte em `AGENT.md`, regra 3).

1. **Base legal da Reforma Tributária** — normas primárias, valem mais que qualquer resumo de terceiros.
   - **Emenda Constitucional nº 132/2023** — cria a base constitucional (IBS, CBS, Imposto Seletivo,
     não-cumulatividade plena, cobrança no destino).
   - **Lei Complementar nº 214/2025** — regulamenta o modelo: alíquotas reduzidas/zero, cesta básica
     (Anexos I, VII, XV), regimes diferenciados/específicos, split payment como modalidade de
     recolhimento.
   - **Lei Complementar nº 227/2026** — institui o Comitê Gestor do IBS (CGIBS), processo administrativo
     tributário do IBS, distribuição da arrecadação, atualizações da LC 214/2025.
   - `⚠️` Nem toda regra já está consolidada — a transição vai até 2033 e normas complementares continuam
     saindo. Verificar data de qualquer fonte secundária antes de citar como definitivo.

2. **Documentação técnica oficial do Split Payment** — anexada ao repositório em `docs-oficiais/`
   (também referenciada pelo simulador Go). Publicada pela RFB/CGIBS/Serpro/Procergs em jun/2026.
   - **Manual de Operações — Split Payment** (v. preliminar) — fluxos de negócio, os 6 arranjos
     (Boleto, Pix Dinâmico, Pix Automático, Pix Estático, TED, TEF), Modelo Inteligente vs Super
     Inteligente, responsabilidades dos PSPs, MOC (Mecanismo de Ocorrências).
   - **Manual de Integração — Plataforma Pública de Split Payment v1.0** — dicionário de campos,
     endpoints REST, política de erros (RFC 7807), headers padrão.
   - **OpenAPI v0.0.10** (`openapi-v0_0_10.json`) — especificação machine-readable dos endpoints.
   - Ato Conjunto RFB/CGIBS nº 02, de 27/05/2026 — aprova os dois manuais acima.

3. **Portal oficial do governo** — para cronograma, regimes diferenciados e o que está em consulta
   pública.
   - **gov.br/fazenda — Reforma Tributária** — https://www.gov.br/fazenda/pt-br/acesso-a-informacao/acoes-e-programas/reforma-tributaria

4. **Imprensa especializada e portais contábeis/jurídicos** — usar apenas pra contexto e cronograma
   quando a fonte 1–3 não detalha; nunca para inventar obrigatoriedade legal. Tag `Fonte: imprensa
   especializada`.
   - Migalhas, Contábeis, Jettax, IOB, Planning, Vanin Contadores — todos consultados em jun/2026; datar
     qualquer citação, pois a reforma muda rápido.
   - `⚠️` Fontes de imprensa às vezes divergem em detalhe fino (ex: alíquota exata de teste em 2026) —
     quando houver conflito entre imprensa e as fontes 1–3, a lei/manual sempre vence.

5. **Nosso simulador (Go)** — não é fonte de verdade tributária. É só um objeto de estudo próprio: serve
   pra explicar *o que decidimos simplificar e por quê*, nunca pra afirmar como o split payment real
   funciona. Tag `Fonte: nosso simulador (simplificação didática)`. Ver `AGENT.md` do simulador
   (`/AGENT.md`, `/TASK.md`) para o que foi implementado.

Se uma afirmação não se sustenta em nenhuma das fontes 1–4, ela é `⚠️ não verificado` — não usar a fonte 5
para justificar uma regra como se fosse real.
