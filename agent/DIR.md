# DIR.md

> Esta árvore descreve a pasta `vault/`, que fica na **raiz do repositório**, ao lado de `agent/`
> (onde estão este arquivo e os outros 4) e do simulador Go (`cmd/`, `internal/`, etc — projeto irmão).
> Não gerar notas dentro de `agent/` — é só plano de controle.

**Layout geral do repositório:**
```
raiz/
├── agent/                   ← AGENT.md, DIR.md, REFERENCES.md, TASKS.md, PROGRESS.md (controle do vault)
├── vault/                   ← o vault de conhecimento em si, descrito abaixo
├── cmd/, internal/, ...     ← código-fonte do simulador Go (projeto irmão)
└── docs-oficiais/           ← Manual de Operações, Manual de Integração, OpenAPI (fontes primárias)
```

Árvore de `vault/` (nomenclatura numerada, estilo Obsidian):

```
vault/
├── 00-Indice.md                              ← Map of Content, ponto de entrada [x]
├── 01-Visao-Geral-da-Reforma.md              ← o que muda e por quê [x]
├── 02-Cronograma-2026-2033.md                ← linha do tempo ano a ano [x]
├── 03-IBS-e-CBS-Basico.md                    ← os dois novos tributos [x]
├── 04-Nao-Cumulatividade-e-Credito.md        ← crédito tributário, conceito central pro B2B [x]
├── 05-Split-Payment-Mecanismo.md             ← como funciona a segregação automática [x]
├── 06-Plataforma-Publica-Documentacao-Tecnica.md  ← Manual de Integração e Swagger oficiais [x]
├── 07-B2B-vs-B2C.md                          ← por que o fluxo muda dependendo de quem compra [x]
├── 08-Simples-Nacional-e-MEI.md              ← tratamento diferenciado pra pequenas empresas [x]
├── 09-Regimes-Especiais-e-Aliquotas-Reduzidas.md  ← cesta básica, saúde, educação etc. [x]
├── 10-Imposto-Seletivo.md                    ← o "imposto do pecado", separado do IBS/CBS [x]
├── 11-Glossario.md                           ← termos e siglas (mínimo 15 termos) [x]
├── 12-Armadilhas-Comuns-para-Devs.md         ← erros típicos ao modelar isso em código [x]
└── 13-Fontes-e-Links-Oficiais.md             ← onde checar a fonte primária [x]
```

`[x]` = concluída (14/14 — ver `PROGRESS.md` para detalhes).

Se uma nova nota for necessária e não couber em nenhum arquivo existente, proponha o novo arquivo aqui
antes de criar — não crie estrutura ad-hoc sem registrar. Diferente do vault Webserv (que usa subpastas
por tema), este vault é enxuto o bastante pra ficar plano, sem subpastas — revisar essa decisão só se
passar de ~20 notas.
