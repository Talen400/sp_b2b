# Regimes Especiais e Alíquotas Reduzidas

## TL;DR
Nem todos os produtos/serviços pagam a alíquota cheia de IBS/CBS. A LC 214/2025
prevê **alíquota zero** para a cesta básica (Anexos I e XV), **redução de 60%**
para saúde, educação e transporte coletivo (Anexo VII), e **redução de 30%**
para profissionais liberais. Há regimes específicos para combustíveis, serviços
financeiros e imóveis.

## Por que isso importa pra quem programa
Sistemas de checkout e ERP precisam mapear **cada produto/serviço** a um código
de regime especial para calcular a alíquota correta. A alíquota não é única —
ela varia por NCM/Serviço e por localização do comprador. Uma modelagem
simplificada (ex: "só tem alíquota padrão") vai falhar em setores como saúde,
educação e alimentos.

## Funcionamento

### Principais reduções

| Tipo | Redução | Exemplos | Base Legal |
|---|---|---|---|
| Cesta básica | 100% (alíquota zero) | Arroz, feijão, pão, frutas, legumes | LC 214/2025, Anexo I e XV |
| Saúde, educação, transporte | 60% | Planos de saúde, escolas, ônibus | LC 214/2025, Anexo VII |
| Profissionais liberais | 30% | Advogados, médicos, engenheiros | LC 214/2025, art. 42 |
| Agronegócio | Redução variável | Insumos agropecuários | LC 214/2025, art. 43 |

✅ `Fonte: LC 214/2025, Anexos I, VII, XV e arts. 42–43`

### Regimes específicos (fora do split padrão)
- **Combustíveis:** monofasia (tributo concentrado no produtor/importador).
  ✅ `Fonte: LC 214/2025, art. 44`
- **Serviços financeiros:** margem de lucro como base de cálculo.
  ✅ `Fonte: LC 214/2025, art. 45`
- **Imóveis:** regime de lucro imobiliário com redução de base.
  ✅ `Fonte: LC 214/2025, arts. 46–50`
- **Cooperativas:** ato cooperativo não tributado.
  ✅ `Fonte: LC 214/2025, art. 51`

> [!note]- 🔍 Aprofundando: Como o sistema determina a alíquota?
> A alíquota aplicável depende de:
> 1. **Código do produto/serviço** (NCM/Serviço) — determina se há redução.
> 2. **Localização do comprador** (destino) — define a alíquota de referência
>    do IBS (RF + RFM).
> 3. **Regime tributário do vendedor** — SIMPLES, lucro real, etc.
> 4. **Tipo de operação** — venda, importação, exportação (imune).
> O sistema precisa consultar uma **tabela de alíquotas** que combina esses
> fatores. ⚠️ Esta tabela ainda não está publicada em formato machine-readable
> — depende de regulamentação futura do CGIBS e da RFB.

## Perguntas de autoavaliação
1. Quais produtos têm alíquota zero de IBS/CBS?
2. Qual a redução para serviços de saúde e educação?
3. Por que a alíquota de um produto depende da localização do comprador?

## Fontes
- ✅ LC 214/2025, Anexos I, VII, XV e arts. 42–51
