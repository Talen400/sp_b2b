# Não-Cumulatividade e Crédito Tributário

## TL;DR
IBS e CBS são **não-cumulativos**: quando uma empresa compra insumo de outra, o
IBS+CBS pagos na compra viram **crédito** que ela pode abater do imposto devido
em suas próprias vendas. Isso evita o efeito "imposto sobre imposto" ao longo
da cadeia produtiva — e é a razão de existir do crédito tributário em B2B.

## Por que isso importa pra quem programa
No [[05-Split-Payment-Mecanismo]], o split segrega o tributo no pagamento, mas
**não** calcula o crédito — isso é feito na apuração contábil. Para o simulador
(`TASK-03` do projeto Go), simplificamos: o crédito é gerado automaticamente
na compra e abatido na venda seguinte, mas no mundo real o crédito exige
escrituração fiscal e pode ser parcial (ex: crédito físico sobre insumos, não
sobre toda despesa). Não confundir a simplificação com a regra real.

## Funcionamento

### Fluxo do crédito na cadeia (exemplo)
```
Fazenda (soja) → Indústria (óleo) → Mercado (venda ao consumidor)
                            ↓
                   Crédito da indústria:
                   IBS+CBS pagos na soja
                   abatidos do IBS+CBS
                   devidos na venda do óleo
```
✅ `Fonte: LC 214/2025, arts. 28–34`

### Regras gerais
- Gera crédito: aquisição de **bens, mercadorias e serviços** utilizados como
  insumo na atividade. ✅ `Fonte: LC 214/2025, art. 28, §1º`
- Não gera crédito: aquisição de bens de uso pessoal, despesas com entretenimento,
  brindes. ✅ `Fonte: LC 214/2025, art. 28, §2º`
- O crédito é apropriado na escrita fiscal do comprador — não é automático.
  ✅ `Fonte: Manual de Operações, seção 6`

### Diferença da simplificação do nosso simulador
No simulador Go (`TASK-03`), o comprador **ganha crédito automaticamente** no
momento da compra, e o vendedor **abate automaticamente** seu saldo disponível.
No mundo real:
- O crédito precisa ser **escriturado** pelo comprador.
- O abatimento ocorre na **declaração periódica** (apuração), não no pagamento.
- Pode haver créditos de períodos anteriores, créditos contestados, etc.
⚠️ `Fonte: nosso simulador (simplificação didática)`

> [!note]- 🔍 Aprofundando: Crédito Financeiro vs Crédito Físico
> No modelo de IVA, há dois tipos de crédito:
> - **Crédito físico:** só gera crédito se o bem adquirido for efetivamente usado
>   na produção (ex: matéria-prima).
> - **Crédito financeiro:** gera crédito independentemente do uso (ex: energia
>   elétrica, aluguel).
> O Brasil adota majoritariamente o **crédito financeiro** no IBS/CBS, ampliando
> a base de crédito em relação ao modelo anterior (PIS/Cofins não-cumulativo,
> que era essencialmente físico). ✅ `Fonte: LC 214/2025, art. 28, §3º`
>
> ⚠️ Isso muda a modelagem de sistemas de apuração: insumos indiretos (como
> frete, seguro) agora geram crédito, o que antes não acontecia no PIS/Cofins.

## Perguntas de autoavaliação
1. O que significa "não-cumulatividade" em um IVA?
2. No fluxo real, o crédito é gerado automaticamente no momento da compra ou
   depende de escrituração?
3. Qual a diferença entre a simplificação do nosso simulador e a regra real de
   crédito tributário?

## Fontes
- ✅ LC 214/2025, arts. 28–34
- ✅ Manual de Operações — Split Payment, seção 6
- ⚠️ nosso simulador (simplificação didática) — TASK-03 do projeto Go
