# Simples Nacional e MEI

## TL;DR
Empresas optantes pelo **Simples Nacional** e **MEI** **não mudam nada em 2026**.
O governo decidirá até set/2026 se esses regimes continuam com IBS/CBS "por
dentro" do DAS (documento único de arrecadação) ou se migram para o split
payment padrão. Até lá, seguem no regime atual.

## Por que isso importa pra quem programa
Se seu sistema atende pequenas empresas (SIMPLES/MEI), você **não precisa**
implementar split payment para elas em 2026–2027 — mas precisa monitorar a
decisão legal. Uma eventual migração pode exigir novas regras de cálculo e
recolhimento específicas para esse segmento.

## Funcionamento

### Situação atual (2026)
- O Simples Nacional e o MEI **não estão sujeitos ao split payment** em 2026.
- O IBS/CBS para essas empresas, se aplicável, será recolhido **dentro do DAS**
  (Documento de Arrecadação do Simples), como ocorre hoje com ICMS, ISS,
  PIS/Cofins.
✅ `Fonte: LC 214/2025, art. 131`

### Decisão pendente (até set/2026)
O Anexo III da LC 214/2025 define que até **setembro de 2026** o governo deve
decidir:
- **Opção A (permanência no DAS):** IBS/CBS continuam incluídos no DAS, sem
  split separado. A empresa não precisa se integrar com a Plataforma Pública.
- **Opção B (migração parcial):** parte do IBS/CBS é segregada via split, parte
  fica no DAS.
- **Opção C (migração total):** SIMPLES/MEI passam a usar o split payment
  padrão, com tratamento diferenciado nas alíquotas.
⚠️ `Fonte: LC 214/2025, Anexo III — situação não definida até jul/2026`

### Efeito prático
- **Sistemas de pequenas empresas:** não precisam de integração com split payment
  por enquanto.
- **Sistemas de PSP:** precisam identificar se o vendedor é optante pelo
  SIMPLES para decidir se aplica split ou não — o campo de regime tributário
  no Documento Fiscal é essencial.
✅ `Fonte: Manual de Operações, seção 4.2`

> [!note]- 🔍 Aprofundando: Implicações de modelagem
> Se a Opção C for escolhida, o split payment para SIMPLES/MEI pode ter:
> - **Alíquotas reduzidas** em relação ao regime padrão.
> - **Faixas de isenção** (ex: transações abaixo de R$ 1.000 não passam por split).
> - **Crédito tributário limitado** para o comprador (já que o SIMPLES não
>   acumula crédito da mesma forma que o lucro real).
> Tudo isso depende de regulamentação futura. ⚠️ não verificado

## Perguntas de autoavaliação
1. O Simples Nacional precisa implementar split payment em 2026?
2. Até quando o governo precisa decidir o futuro do SIMPLES no split payment?
3. Qual campo do Documento Fiscal é essencial para o PSP decidir se aplica
   split?

## Fontes
- ✅ LC 214/2025, art. 131 e Anexo III
- ✅ Manual de Operações, seção 4.2
