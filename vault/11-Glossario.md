# Glossário

## A
- **Alíquota de referência:** percentual do IBS definido anualmente pelo Senado
  para cada estado (RF) e município (RFM). Varia por localização do comprador.
  ✅ `Fonte: LC 214/2025, arts. 7–12`

## C
- **CBS (Contribuição sobre Bens e Serviços):** tributo federal que substitui
  PIS e Cofins. Administrado pela RFB. Alíquota única nacional.
  ✅ `Fonte: EC 132/2023, art. 195`
- **CGIBS (Comitê Gestor do IBS):** órgão responsável por administrar o IBS,
  definir normas e gerir a Plataforma Pública. Criado pela LC 227/2026.
  ✅ `Fonte: LC 227/2026`
- **Crédito tributário:** valor do IBS/CBS pago em uma compra que a empresa
  compradora pode abater do imposto devido em suas vendas (não-cumulatividade).
  ✅ `Fonte: LC 214/2025, arts. 28–34`

## D
- **Documento Fiscal:** NF-e, NFC-e ou NFS-e que formaliza a operação. Contém
  os campos de IBS, CBS, IS e demais tributos. Orientação do split.
  ✅ `Fonte: Manual de Operações, seção 4`

## E
- **EC 132/2023:** Emenda Constitucional que criou a base da Reforma Tributária
  (IBS, CBS, Imposto Seletivo, não-cumulatividade, cobrança no destino).
  ✅ `Fonte: EC 132/2023`

## I
- **IBS (Imposto sobre Bens e Serviços):** tributo estadual/municipal que
  substitui ICMS e ISS. Administrado pelo CGIBS. Alíquota varia por
  estado/município do comprador. ✅ `Fonte: EC 132/2023, art. 156-A`
- **Imposto Seletivo (IS):** imposto federal extrafiscal sobre produtos nocivos
  (cigarro, álcool, combustíveis fósseis). Fora do split payment.
  ✅ `Fonte: EC 132/2023, art. 153, §6º`

## L
- **LC 214/2025:** Lei Complementar que regulamentou o modelo do IBS/CBS
  (alíquotas reduzidas, regimes especiais, split payment).
  ✅ `Fonte: LC 214/2025`
- **LC 227/2026:** Lei Complementar que instituiu o CGIBS e o processo
  administrativo tributário do IBS.
  ✅ `Fonte: LC 227/2026`

## M
- **MOC (Mecanismo de Ocorrências):** sistema da Plataforma Pública para
  registrar e consultar falhas no split payment (ex: split não realizado,
  valor divergente). ✅ `Fonte: Manual de Operações, seção 7`

## N
- **Não-cumulatividade:** princípio pelo qual o tributo pago em uma etapa da
  cadeia é abatido na etapa seguinte, evitando "imposto sobre imposto".
  ✅ `Fonte: EC 132/2023, art. 156-A, §1º`

## P
- **Plataforma Pública do Split Payment:** API REST (RFB/CGIBS) que os PSPs
  consultam para obter valores segregados. Documentada no Manual de Integração
  + OpenAPI. ✅ `Fonte: Manual de Integração v1.0`
- **PSP (Instituição de Pagamento):** entidade responsável por processar o
  pagamento e executar o split (ex: bancos, fintechs).
  ✅ `Fonte: Manual de Operações`

## R
- **RFB (Receita Federal do Brasil):** administradora da CBS e da Plataforma
  Pública (em conjunto com o CGIBS). ✅ `Fonte: Manual de Integração`

## S
- **Split payment:** segregação automática do IBS+CBS no momento da liquidação
  financeira. O PSP divide o valor em líquido do vendedor + imposto para o
  fisco. ✅ `Fonte: Manual de Operações, seção 1.1`

## ⚠️ Marcados como não verificado
- **Modelo Inteligente vs Super Inteligente:** conceitos documentados no Manual
  de Operações, mas sem confirmação de que ambos estarão operacionais já em 2027.
  ✅ (ver seções 2.1–2.2)
