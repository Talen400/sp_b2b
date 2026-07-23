# Split Payment — Mecanismo

## TL;DR
Split payment é a segregação automática do IBS e CBS no momento da liquidação
financeira. O PSP (instituição de pagamento) divide o valor da transação em:
**líquido do vendedor** + **IBS** + **CBS**. O tributo vai direto ao fisco sem
passar pelo caixa do vendedor. O comprador (se PJ) acumula crédito tributário.

## Por que isso importa pra quem programa
Para sistemas de pagamento e ERP, o split payment muda o fluxo pós-venda: em vez de
"recebeu o valor bruto e depois recolhe", o sistema precisa enviar metadados fiscais
junto com a instrução de pagamento e lidar com liquidação parcial (split em duas
etapas). O MOC (Mecanismo de Ocorrências) permite tratamento de falhas na segregação.

## Funcionamento

### Fluxo básico (Modelo Inteligente)
1. Vendedor emite Documento Fiscal com os dados da transação.
2. Comprador (ou sistema) inicia o pagamento no PSP, informando o identificador do
   Documento Fiscal.
3. PSP consulta a **Plataforma Pública** para obter os valores segregados (líquido,
   IBS, CBS) — ou recebe esses valores no próprio payload, dependendo do arranjo.
4. PSP liquida: líquido vai para o vendedor; IBS+CBS vai para a conta do fisco.
5. MOC registra o status da operação (sucesso, falha, pendência).
✅ `Fonte: Manual de Operações — Split Payment, seção 2.1`

### Os 6 arranjos de pagamento
| Arranjo | Split em 1 etapa? | Consulta prévia? |
|---|---|---|
| Pix Dinâmico | Sim | Sim |
| Pix Automático | Sim | Sim |
| Pix Estático | Não (2 etapas) | Sim |
| Boleto | Não (2 etapas) | Sim |
| TED | Depende | Sim |
| TEF (débito) | Sim | Sim |

✅ `Fonte: Manual de Operações, seção 5`

### Modelo Inteligente vs Super Inteligente
- **Inteligente:** PSP consulta a Plataforma Pública para obter os valores do split.
  É o modelo padrão.
- **Super Inteligente:** O PSP já recebe os valores do split no próprio payload da
  transação (ex: Pix Automático com campos estendidos) e não precisa consultar a
  plataforma.
✅ `Fonte: Manual de Operações, seção 2.2`

> [!note]- 🔍 Aprofundando: Split em duas etapas (Pix Estático e Boleto)
> Em arranjos sem consulta prévia ao Documento Fiscal no momento da autorização, o
> split ocorre em duas etapas:
> 1. **Liquidação financeira:** o valor integral vai para o PSP do vendedor.
> 2. **Split propriamente dito:** o PSP do vendedor consulta a Plataforma Pública
>    com os dados da transação e repassa o tributo ao fisco.
> O risco de crédito fica com o PSP do vendedor até a conclusão da etapa 2.
> ✅ `Fonte: Manual de Operações, seção 5.3`

## Perguntas de autoavaliação
1. Qual a diferença entre o fluxo do split no Modelo Inteligente vs Super Inteligente?
2. Por que Pix Estático exige split em duas etapas?
3. O que é o MOC e para que serve?

## Fontes
- ✅ Manual de Operações — Split Payment, seções 2.1, 2.2, 5, 5.3
