# Glossário

## A
- **Alíquota de referência:** percentual do IBS definido anualmente pelo Senado para cada estado (RF) e município (RFM). Varia por localização do comprador. ✅ `LC 214/2025, arts. 7–12`

## B
- **Boleto:** arranjo de pagamento iniciado pelo Recebedor. Modelo Super Inteligente. Possui Informe de Transação Iniciada, Informe de Transação Atualizada, Retorno Super Inteligente, Informe Preliminar de Pagamento e Informe de Segregação. ✅ `Manual de Operações, seção 2`

## C
- **CBS (Contribuição sobre Bens e Serviços):** tributo federal que substitui PIS e Cofins. Administrado pela RFB. Alíquota única nacional. ✅ `EC 132/2023, art. 195`
- **CGIBS (Comitê Gestor do IBS):** órgão responsável por administrar o IBS, definir normas e gerir a Plataforma Pública. Criado pela LC 227/2026. ✅ `LC 227/2026`
- **Correlation-Id:** header UUID4 obrigatório nas requisições à PP. Rastreio ponta a ponta: mesmo valor ao longo de toda a cadeia de uma transação. ✅ `Manual de Integração, seção 4`
- **Crédito tributário:** valor do IBS/CBS pago em uma compra que a empresa compradora pode abater do imposto devido em suas vendas (não-cumulatividade). ✅ `LC 214/2025, arts. 28–34`

## D
- **Decimal(18,2):** formato de campo monetário na PP: 18 dígitos, 2 casas decimais, representado como string no JSON. Ex: `"1234567890123456.78"`. ✅ `Manual de Integração, seção 4.2`
- **Documento Fiscal:** NF-e, NFC-e ou NFS-e que formaliza a operação. Contém os campos de IBS, CBS, IS. ✅ `Manual de Operações, seção 4`

## E
- **EC 132/2023:** Emenda Constitucional que criou a base da Reforma Tributária (IBS, CBS, Imposto Seletivo, não-cumulatividade, cobrança no destino). ✅ `EC 132/2023`

## I
- **IBS (Imposto sobre Bens e Serviços):** tributo estadual/municipal que substitui ICMS e ISS. Administrado pelo CGIBS. Alíquota varia por estado/município do comprador. ✅ `EC 132/2023, art. 156-A`
- **Imposto Seletivo (IS):** imposto federal extrafiscal sobre produtos nocivos. Fora do split payment. ✅ `EC 132/2023, art. 153, §6º`
- **Informe de Transação Iniciada:** primeiro informe do fluxo. Enviado pelo PSP Recebedor ao criar a cobrança (Boleto, Pix Dinâmico, Pix Automático). Contém os valores de IBS/CBS do documento fiscal. ✅ `Manual de Operações, seção 4.1`
- **Informe de Transação Atualizada:** atualização do Informe de Transação Iniciada. Usado quando há correção no documento fiscal antes do pagamento. Aplicável a Boleto e Pix Dinâmico. ✅ `Manual de Operações, seção 4.1`
- **Informe Preliminar de Pagamento:** comunicação informativa e não vinculante de que uma transação foi paga. Enviada logo após o pagamento. Não gera obrigação de repasse. Descartada ao receber o Informe de Segregação correspondente. Obrigatória em todos arranjos exceto TEF (opcional). ✅ `Manual de Operações, seção 4.3`
- **Informe de Segregação:** comunicação definitiva e vinculante. Enviada em lote, 2x por dia útil, contendo todas as transações liquidadas em uma janela. Gera obrigação de Repasse Financeiro. Obrigatória em todos os 6 arranjos. ✅ `Manual de Operações, seção 4.3`

## L
- **LC 214/2025:** Lei Complementar que regulamentou o modelo do IBS/CBS (alíquotas reduzidas, regimes especiais, split payment). ✅
- **LC 227/2026:** Lei Complementar que instituiu o CGIBS e o processo administrativo tributário do IBS. ✅
- **Long Polling:** mecanismo de pull-based event streaming usado pela PP para entregar mensagens do Super Inteligente aos PSPs. PP mantém conexão aberta até surgir mensagem ou atingir timeout. ✅ `Manual de Integração, seção 3.6`

## M
- **Message-Id:** header UUID4 obrigatório nas requisições à PP. Usado para idempotência: mesma Message-Id em retentativas é ignorada pela PP. ✅ `Manual de Integração, seção 4`
- **MOC (Mecanismo de Ocorrências):** sistema da PP para registrar e consultar falhas no split (ex: split não realizado, valor divergente). ✅ `Manual de Operações, seção 7`
- **Modelo Inteligente:** modelo de split para arranjos iniciados pelo Pagador (Pix Estático, TED, TEF). O valor do tributo é o que o Pagador informou, sem correção do governo. ✅ `Manual de Operações, seção 2.1`
- **Modelo Super Inteligente:** modelo de split para arranjos iniciados pelo Recebedor (Boleto, Pix Dinâmico, Pix Automático). O governo pode corrigir os valores via Retorno Super Inteligente. ✅ `Manual de Operações, seção 2.2`

## N
- **Não-cumulatividade:** princípio pelo qual o tributo pago em uma etapa da cadeia é abatido na etapa seguinte, evitando "imposto sobre imposto". ✅ `EC 132/2023, art. 156-A, §1º`
- **NSU (Número Sequencial Único):** identificador único de uma transação dentro da PP. Usado como path param em atualizações e consultas. ✅ `Manual de Integração, seção 3.6.1`

## P
- **Plataforma Pública do Split Payment (PP):** API REST (RFB/CGIBS/Serpro) que os PSPs consultam para comunicar transações, receber retornos do Super Inteligente e enviar informes de segregação. ✅ `Manual de Integração v1.0`
- **PSP (Instituição de Pagamento):** entidade responsável por processar o pagamento e executar o split (ex: bancos, fintechs).
- **PSP Pagador Direto:** PSP que detém a conta do Pagador. ✅ `Manual de Operações, seção 1`
- **PSP Pagador Indireto:** PSP que intermediou a transação mas não detém a conta do Pagador (ex: sub-adquirente). ✅ `Manual de Operações, seção 1`
- **PSP Recebedor Direto:** PSP que detém a conta do Recebedor. **Responsável regulatório** pela comunicação com a PP. ✅ `Manual de Operações, seção 1`
- **PSP Recebedor Indireto:** PSP que intermediou o recebimento mas não detém a conta do Recebedor. ✅ `Manual de Operações, seção 1`

## R
- **Repasse Financeiro:** transferência do valor segregado de CBS e IBS ao governo (RFB e CGIBS). Ocorre em D+N (dias úteis após a liquidação). CBS via TES, IBS via STR. ✅ `Manual de Operações, seção 6`
- **Retorno Super Inteligente:** mensagem enviada pela PP ao PSP com correções de valores de tributo (CBS/IBS Corrigido, Em Aberto). Consultada via long polling. ✅ `Manual de Integração, seção 3.6`
- **RFB (Receita Federal do Brasil):** administradora da CBS e da Plataforma Pública (em conjunto com o CGIBS). ✅ `Manual de Integração`

## S
- **Split payment:** segregação automática do IBS+CBS no momento da liquidação financeira. O PSP divide o valor em líquido do vendedor + imposto para o fisco. ✅ `Manual de Operações, seção 1.1`
- **Stream:** canal de mensagens do Super Inteligente baseado em long polling. Identificado por um `streamId`. Pode ser encerrado via DELETE. ✅ `Manual de Integração, seção 3.6`

## T
- **TED:** arranjo de transferência eletrônica disponível. Iniciado pelo Pagador. Modelo Inteligente. ✅ `Manual de Operações, seção 2`
- **TEF:** arranjo de transferência eletrônica de fundos (débito em conta). Iniciado pelo Pagador. Modelo Inteligente. Informe Preliminar de Pagamento opcional. ✅ `Manual de Operações, seção 2`
- **Tenant-Id:** header obrigatório nas requisições à PP. Identificador do PSP na plataforma. ✅ `Manual de Integração, seção 4`
- **Timestamp:** header obrigatório nas requisições à PP. Momento de geração da mensagem no formato RFC 3339. ✅ `Manual de Integração, seção 4`
- **Token de Posição:** token retornado pela PP no header `proximoToken` durante o long polling. Indica a posição do último lote de mensagens entregue. Usado na URL da próxima consulta. ✅ `Manual de Integração, seção 3.6`

## V
- **Valor Aplicado:** valor de tributo exibido ao Pagador no extrato/fatura. Pode divergir do segregado. ✅ `Manual de Operações, seção 3.1`
- **Valor Corrigido:** valor do tributo após correção do Super Inteligente. ✅ `Manual de Operações, seção 3.1`
- **Valor Em Aberto:** parcela do débito ainda não liquidada (ex: boleto pago em parte). ✅ `Manual de Operações, seção 3.1`
- **Valor Informado:** valor de IBS/CBS constante no Documento Fiscal. ✅ `Manual de Operações, seção 3.1`
- **Valor Segregado:** valor efetivamente separado no momento da liquidação. Base do Repasse Financeiro. ✅ `Manual de Operações, seção 3.1`
