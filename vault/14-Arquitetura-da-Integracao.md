# Arquitetura da Integração — Split Payment + Plataforma Pública

## TL;DR
Este documento descreve como os atores do ecossistema de pagamentos se relacionam com o split payment, onde cada parte se conecta com a Plataforma Pública (PP), e como nosso simulador se posiciona nesse ecossistema.

## Diagrama de atores

```
                   +------------------+
                   |   Recebedor      |  (quem vende, cria o Documento Fiscal)
                   |  (loja/ERP)      |
                   +--------+---------+
                            |
                     Informa IBS/CBS
                     no Doc. Fiscal
                            |
                            v
                   +------------------+
                   | PSP Recebedor    |  ← Responsável regulatório pela PP
                   |  Direto          |
                   +--------+---------+
                            |
                  +---------+---------+
                  |                   |
          +-------v-------+   +-------v-------+
          | PSP Recebedor |   |  Plataforma   |
          |  Indireto     |   |  Pública (PP) |  ← API REST do governo
          | (opcional)    |   +-------+-------+
          +-------+-------+           |
                  |                   |
                  |         +---------v---------+
                  |         | Super Inteligente |  ← Corrige valores
                  |         | (RFB/CGIBS)       |
                  |         +-------------------+
                  |
          +-------v-------+
          |   CIP / STR   |  ← Liquidação interbancária
          +-------+-------+
                  |
          +-------v-------+
          | PSP Pagador   |
          |  Direto       |
          +-------+-------+
                  |
          +-------v-------+
          |   Pagador     |  (quem compra/paga)
          |  (PJ)         |
          +---------------+
```

## Onde o split acontece

O split **não** acontece na PP. O split acontece dentro do PSP no momento da **liquidação financeira interbancária** (via CIP/STR). O PSP:

1. Recebe o valor bruto do Pagador
2. Calcula o tributo devido (IBS + CBS) com base no Documento Fiscal
3. Segrega: envia o tributo ao governo e o líquido ao Recebedor

A PP é um sistema de **informação e controle**: os PSPs comunicam as transações para a PP, e o governo pode corrigir valores. O split em si é executado pelo PSP.

## Obrigação regulatória de comunicação

Quem responde perante a PP é o **PSP Recebedor Direto** (detentor da conta do Recebedor). Mesmo que exista um PSP Indireto na relação comercial:
- O PSP Recebedor Direto é responsável por todos os informes
- O PSP Recebedor Indireto pode auxiliar na comunicação, mas a responsabilidade final é do Direto
- O `Tenant-Id` nos headers deve identificar o PSP Recebedor Direto

## Categorias de valor no ciclo de vida

```
Documento Fiscal       Pagamento         Liquidação         Repasse
     |                    |                  |                  |
vlInf ──► vlCorr ──► vlAberto ──► vlSeg ──► vlApl
     (criação)    (correção    (parcial)   (efetivo)   (exibido ao
                   Super                               Pagador)
                   Inteligente)

     Informado ──► Corrigido ──► Em Aberto ──► Segregado ──► Aplicado
```

## Modelo Inteligente vs Super Inteligente

### Inteligente (Pix Estático, TED, TEF)
- Iniciado pelo Pagador
- Pagador informa o valor do tributo (via sistema do PSP Pagador)
- Sem correção do governo
- Fluxo: Informe Preliminar de Pagamento → Informe de Segregação
- `vlInf = vlCorr = vlAberto = vlSeg`

### Super Inteligente (Boleto, Pix Dinâmico, Pix Automático)
- Iniciado pelo Recebedor
- Recebedor informa o valor no Documento Fiscal
- Governo pode corrigir via Retorno Super Inteligente (long polling)
- Fluxo: Informe Transação Iniciada → (Retorno SI) → Informe Preliminar de Pagamento → Informe de Segregação
- `vlInf ≠ vlCorr` (se houve correção)

## Fluxo temporal completo

```
T0  ── Recebedor cria cobrança (Boleto/Pix)
    ── PSP Recebedor envia Informe de Transação Iniciada → PP
    ── (se Super Inteligente) PP retorna valores corrigidos via long polling

T1  ── Pagador paga (Pix, TED, TEF, Boleto) —> PSP Pagador

T2  ── PSP Pagador liquida via CIP/STR
    ── PSP Recebedor recebe notificação de liquidação

T3  ── PSP Recebedor envia Informe Preliminar de Pagamento → PP
    ── (se aplicável) PSP Recebedor envia Informe de Baixa → PP

T4  ── Fim da janela de consolidação
    ── PSP Recebedor envia Informe de Segregação (lote) → PP
    ── Informe de Segregação gera obrigação de Repasse Financeiro

T5  ── PSP Recebedor realiza Repasse Financeiro em D+N
    ── CBS → TES, IBS → STR
```

## Como nosso simulador se posiciona

Nosso simulador é uma **API didática** que:
- Implementa a lógica de cálculo de split (CalculateSplit)
- Mantém crédito tributário simplificado (não-cumulatividade B2B)
- Expõe endpoints REST para criar empresas e transações
- **NÃO** integra com a PP real (não temos credencial de PSP)

A **Fase 7** adiciona um cliente HTTP que **pode** falar com o mock Prism da PP, validando que nossos payloads seguem o contrato real. Mas isso é opcional e demonstrativo.

## O que o simulador simplifica

| Aspecto real | No simulador |
|---|---|
| mTLS + OAuth 2.0 | Não implementado (mock sem autenticação) |
| CNPJ alfanumérico | String livre |
| Múltiplos PSPs (Direto/Indireto) | Um único "PSP" implícito |
| Documento Fiscal completo | Só valor bruto + alíquotas |
| Arranjo específico (6 tipos) | Fixo em "boleto" para o hook |
| Retorno Super Inteligente | Mock Prism valida schema, sem lógica |
| Crédito tributário real | Simplificação didática |
