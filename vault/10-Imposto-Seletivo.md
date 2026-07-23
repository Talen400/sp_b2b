# Imposto Seletivo (IS)

## TL;DR
O **Imposto Seletivo** é o "imposto do pecado" — incide sobre produtos nocivos
à saúde (bebidas alcoólicas, cigarro) e ao meio ambiente (combustíveis fósseis,
agrotóxicos). É federal, extrafiscal (desestimular o consumo), e fica **fora
do split payment** — continua sendo recolhido separadamente.

## Por que isso importa pra quem programa
O Imposto Seletivo **não passa pelo split payment** — isso é explícito no
Manual de Operações. Sistemas que lidam com produtos sujeitos ao IS precisam
calcular e recolher esse imposto à parte, sem envolver a Plataforma Pública.
Ignorar essa separação pode levar a erros de integração.

## Funcionamento

### Características
- **Federal:** arrecadação da União.
- **Extrafiscal:** objetivo principal é desestimular o consumo, não arrecadar.
- **Não-cumulativo:** em princípio, mas com direito a crédito apenas nas etapas
  seguintes da cadeia (não para o consumidor final).
✅ `Fonte: EC 132/2023, art. 153, §6º`

### Produtos sujeitos
- Bebidas alcoólicas.
- Cigarros e derivados do tabaco.
- Combustíveis fósseis.
- Agrotóxicos.
- Veículos (em discussão — depende de regulamentação).
⚠️ `Fonte: LC 214/2025, art. 52 — lista pode ser ampliada por lei ordinária`

### Relação com o split payment
O Manual de Operações, seção 1.2, é claro:
> "O split payment se aplica exclusivamente ao IBS e à CBS. Ficam excluídos
> o Imposto Seletivo, as contribuições previdenciárias e demais tributos não
> abrangidos pela reforma."
✅ `Fonte: Manual de Operações, seção 1.2`

> [!note]- 🔍 Aprofundando: Imposto Seletivo na Nota Fiscal
> Apesar de não passar pelo split, o IS deve constar no Documento Fiscal como
> campo separado. Ou seja, a NF-e/NFS-e terá campos para:
> - IBS
> - CBS
> - IS (quando aplicável)
> - Valor líquido
> Isso significa que o sistema de emissão de DF precisa calcular os três
> tributos mesmo que só os dois primeiros sejam segregados no pagamento.
> ✅ `Fonte: Manual de Operações, seção 4.2`

## Perguntas de autoavaliação
1. O Imposto Seletivo passa pelo split payment?
2. Qual o objetivo principal do IS: arrecadação ou desestímulo ao consumo?
3. O IS precisa constar no Documento Fiscal mesmo não passando pelo split?

## Fontes
- ✅ EC 132/2023, art. 153, §6º
- ✅ LC 214/2025, art. 52
- ✅ Manual de Operações — Split Payment, seções 1.2, 4.2
