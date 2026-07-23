# Visão Geral da Reforma Tributária

## TL;DR
A Reforma Tributária do consumo (EC 132/2023 + LC 214/2025) substitui cinco tributos
(ICMS, ISS, IPI, PIS, Cofins) por dois IVAs — **IBS** (estadual/municipal) e **CBS**
(federal) — mais o **Imposto Seletivo** (produtos nocivos). A cobrança passa da origem
para o **destino**, e o **split payment** é o mecanismo que segrega o tributo no
momento do pagamento.

## Por que isso importa pra quem programa
Toda a lógica de checkout, faturamento e contabilidade de sistemas B2B vai mudar: em
vez de a empresa calcular, recolher e declarar tributos *depois* da venda, o sistema
financeiro fará a segregação **no ato do pagamento**. Isso exige novos campos em
documentos fiscais, novas integrações com PSPs (instituições de pagamento) e nova
lógica de crédito tributário.

## Funcionamento

### O que muda
- **Antes:** ICMS/ISS (estadual/municipal) + PIS/Cofins (federal) — cada um com regras
  próprias, cumulatividade parcial, guerra fiscal entre estados.
- **Depois:** IBS + CBS — dois IVAs padronizados, não-cumulativos, cobrança no destino.
  O contribuinte lida com dois tributos em vez de cinco. ✅ `Fonte: EC 132/2023, art. 156-A e 195`

### O que é split payment
É a modalidade de recolhimento em que o valor do IBS e da CBS é separado do valor
do pagamento **no momento da liquidação financeira**, pelo PSP (instituição de
pagamento). O vendedor recebe o líquido; o tributo vai direto ao fisco.
✅ `Fonte: Manual de Operações — Split Payment, seção 1.1`

### Impacto em sistemas
- **Checkout:** precisa informar se a transação é B2B ou B2C, pois só empresa
  (B2B) tem visibilidade ampla dos campos fiscais. ✅ `Fonte: Manual de Operações, seção 8`
- **Documento fiscal:** novo leiaute com campos de split. ✅ `Fonte: Manual de Operações, seção 4`
- **Conciliação:** o MOC (Mecanismo de Ocorrências) permite rastrear falhas de split.

> [!note]- 🔍 Aprofundando: Transição 2026–2033
> - **2026:** CBS com alíquota reduzida (0,9% — caráter informativo). Split payment
>   vale para PIX e débito.
> - **2027:** CBS plena (8,8%). Split payment obrigatório para todos os meios.
> - **2029–2032:** Transição do IBS (estadual/municipal). Alíquota de referência
>   estadual (RF) e municipal (RFM) definidas anualmente pelo Senado.
> - **2033:** Regime pleno — tudo consolidado.
> ✅ `Fonte: LC 214/2025, arts. 119-127 e Anexo II`

## Perguntas de autoavaliação
1. Quantos tributos o IBS substitui? E a CBS?
2. Por que uma empresa B2B precisa de campos fiscais adicionais no checkout que uma
   B2C não precisa?
3. Em que ano o split payment se torna obrigatório para todos os meios de pagamento?

## Fontes
- ✅ EC 132/2023, arts. 156-A, 195
- ✅ LC 214/2025, arts. 119-127 e Anexo II
- ✅ Manual de Operações — Split Payment, seções 1.1, 4, 8
