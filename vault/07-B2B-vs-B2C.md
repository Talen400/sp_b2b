# B2B vs B2C

## TL;DR
No split payment, o fluxo muda radicalmente entre **B2B** (empresa comprando de
empresa) e **B2C** (consumidor final comprando). No B2B, o comprador tem
visibilidade ampla dos campos fiscais e acumula **crédito tributário**. No B2C,
o consumidor final não vê os campos de split — o tributo é segregado
"por dentro" do valor.

## Por que isso importa pra quem programa
Seu checkout precisa saber se o comprador é **PJ ou PF** para decidir:
- Quais campos fiscais exibir no momento do pagamento.
- Se o split será "informado" (B2B, comprador vê os valores) ou "não informado"
  (B2C, comprador só vê o total).
- Se a transação gera crédito tributário para o comprador (B2B sim, B2C não).

## Funcionamento

### Princípios do Manual de Operações
- **Princípio 1 (B2C):** O Pagador Efetivo Pessoa Física **não tem visibilidade**
  dos campos fiscais do split. O valor do tributo fica embutido no total.
  ✅ `Fonte: Manual de Operações, seção 8`
- **Princípio 2 (B2B):** O Pagador Efetivo Pessoa Jurídica **tem visibilidade**
  ampla dos campos fiscais e pode consultar os valores segregados.
  ✅ `Fonte: Manual de Operações, seção 8`

### Impacto no split
| Característica | B2B | B2C |
|---|---|---|
| Visibilidade dos campos fiscais | Sim | Não |
| Geração de crédito tributário | Sim (comprador PJ) | Não |
| Valor do split visível no checkout | Sim | Não (embutido) |
| Documento Fiscal obrigatório | Sim (NF-e/NFS-e) | Sim (pode ser cupom) |

### Implicação para crédito
No B2B, o comprador PJ acumula crédito do IBS+CBS pagos. No B2C, o consumidor
final **não acumula crédito** — ele é o elo final da cadeia e efetivamente
arca com o tributo. É essa diferença que torna a simulação B2B (nosso projeto
Go, TASK-03) diferente de uma simulação B2C. ✅ `Fonte: LC 214/2025, art. 28`

> [!note]- 🔍 Aprofundando: Como o sistema distingue B2B de B2C?
> A distinção é feita pelo **Documento Fiscal** vinculado ao pagamento. Se o
> Documento Fiscal é uma NF-e (destinatário PJ), a transação é classificada
> como B2B. Se é uma NFC-e (consumidor final), é B2C. O PSP não decide — ele
> apenas executa o split com base nos dados que recebe do Documento Fiscal.
> ✅ `Fonte: Manual de Operações, seção 4`
>
> ⚠️ Isso tem implicações de arquitetura: o sistema de checkout precisa enviar
> o identificador do Documento Fiscal correto para o PSP, e o tipo de documento
> (NF-e vs NFC-e) determina o fluxo de split.

## Perguntas de autoavaliação
1. Quem decide se uma transação é B2B ou B2C para fins de split?
2. Um consumidor PF que compra de uma empresa gera crédito tributário?
3. O que muda na interface de checkout entre B2B e B2C?

## Fontes
- ✅ Manual de Operações — Split Payment, seção 8 (Princípios 1 e 2)
- ✅ Manual de Operações, seção 4
- ✅ LC 214/2025, art. 28
