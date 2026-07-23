# Split Payment Simulator (B2B)

Simulador educacional do mecanismo de **split payment** previsto na Reforma
Tributária brasileira (EC 132/2023 e LC 214/2025), focado em contexto **B2B**
(empresa vendendo para empresa) com crédito tributário não-cumulativo de
IBS e CBS.

## Como rodar

```bash
# Verificar que o setup está ok
go run ./cmd/simulador
```

## Estrutura

```
cmd/simulador/   — CLI principal
cmd/demo/        — Cenário de demonstração
internal/split/  — Cálculo do split (IBS, CBS, líquido)
internal/company — Struct de empresa e lógica de crédito
internal/store/  — Armazenamento em memória
```

## Crédito tributário em B2B

Diferente de vendas a consumidor final, em transações B2B o comprador acumula
crédito tributário equivalente ao IBS+CBS pagos na compra. Esse crédito pode
ser abatido do imposto devido em vendas futuras — evitando o efeito cascata
(imposto sobre imposto) ao longo da cadeia produtiva.

## Licença

MIT
