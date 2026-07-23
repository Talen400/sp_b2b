# Armadilhas Comuns para Pessoas Desenvolvedoras

## TL;DR
Modelar split payment em código tem armadilhas específicas: usar float para
valor monetário, confundir "valor informado" com "valor segregado", assumir
que CNPJ é só numérico, e esquecer que split não é apuração. Esta nota lista
os erros mais comuns e como evitá-los.

## Por que isso importa pra quem programa
Erros de modelagem tributária podem gerar inconsistências contábeis,
recolhimento incorreto e risco fiscal para as empresas integradas. Conhecer
as armadilhas antes de programar evita retrabalho.

## Armadilhas

### 1. Usar float/double para valores monetários
O Manual de Integração especifica que valores monetários são em **centavos**
(campo `valorBruto` como inteiro). Usar ponto flutuante causa erros de
arredondamento — um centavo de diferença em milhões de transações vira um
problema grande.
✅ `Fonte: Manual de Integração, seção 4.2 — campo valorBruto`
⚠️ **No nosso simulador Go:** usamos `int64` (centavos) — solução correta.

### 2. Confundir "Informado" com "Segregado"
- **Valor Informado:** o que consta no Documento Fiscal (quanto de IBS/CBS
  deve ser pago naquela operação).
- **Valor Segregado:** o que efetivamente foi separado no momento do pagamento
  (pode divergir do informado em caso de pagamento parcial, desconto, etc.).
O split opera sobre o **valor Segregado**, não sobre o Informado.
✅ `Fonte: Manual de Operações, seção 3.1`

### 3. Assumir que CNPJ é só numérico
A partir de **julho/2026**, o CNPJ pode conter letras (CNPJ Alfanumérico),
seguindo o padrão RFC 6030. Sistemas que validam CNPJ como "apenas dígitos"
vão quebrar.
⚠️ `Fonte: imprensa especializada — RFB publicou instrução normativa em mai/2026`

### 4. Esquecer que split não é apuração
O split segrega o tributo no pagamento, mas **não calcula o crédito**
tributário — isso é feito na apuração contábil periódica (geralmente mensal).
Uma transação pode ter split correto mas crédito incorreto se o sistema não
escriturar os valores adequadamente.
⚠️ `Fonte: nosso simulador (simplificação didática) — no mundo real, escrituração é separada`

### 5. Ignorar a diferença de alíquota por UF
O IBS é cobrado no **destino**. Um sistema que usa sempre a alíquota do
vendedor (origem) calcula o tributo errado. A alíquota varia por estado e
município do comprador.
✅ `Fonte: EC 132/2023, art. 156-A, §1º`

### 6. Tratar B2B e B2C como iguais
Como visto em [[07-B2B-vs-B2C]], no B2C o consumidor não vê os campos fiscais.
Se o checkout expõe valores de IBS/CBS para um consumidor PF, está fora da
regra.
✅ `Fonte: Manual de Operações, seção 8`

### 7. Subestimar a importância do MOC
O Mecanismo de Ocorrências não é opcional — falhas de split precisam ser
registradas e tratadas. Um sistema que ignora o MOC pode perder o rastreio
de splits mal-sucedidos.
✅ `Fonte: Manual de Operações, seção 7`

## Perguntas de autoavaliação
1. Qual o tipo de dado correto para representar valores monetários na API da
   Plataforma Pública?
2. Qual a diferença entre valor "Informado" e valor "Segregado"?
3. Por que usar a alíquota do vendedor em vez da do comprador dá erro no IBS?

## Fontes
- ✅ Manual de Integração, seção 4.2
- ✅ Manual de Operações, seções 3.1, 7, 8
- ✅ EC 132/2023, art. 156-A, §1º
- ⚠️ imprensa especializada (CNPJ alfanumérico)
- ⚠️ nosso simulador (simplificação didática)
