package domain

// ApplyCredit acumula o crédito tributário gerado por uma compra.
//
// No modelo não-cumulativo do IBS/CBS, a empresa compradora acumula
// crédito no valor total do imposto pago (IBS + CBS segregados).
// Esse crédito poderá ser abatido em vendas futuras (ver UseCredit).
//
// Parâmetros:
//   - company: a empresa compradora (tem o saldo_credito atualizado in-place).
//   - split: resultado do split da transação de compra.
//
// Retorna o valor do crédito gerado nesta transação (em centavos).
//
// Atenção: esta é uma simplificação didática. No mundo real, o crédito
// tributário B2B não faz parte do escopo do Manual de Operações do split
// payment — a apuração e escrituração do crédito é feita na contabilidade
// periódica da empresa, não no momento do pagamento.
// Fonte: simplificação do simulador (não corresponde a regra normativa real).
func ApplyCredit(company *Company, split SplitResult) (creditoGerado int64) {
	creditoGerado = split.ValorIBS + split.ValorCBS
	company.SaldoCredito += creditoGerado
	return creditoGerado
}

// UseCredit abate o crédito acumulado do imposto devido em uma venda.
//
// Quando uma empresa vende, ela deve recolher o IBS+CBS da venda, mas
// pode abater até o valor total do imposto devido usando o crédito que
// acumulou em compras anteriores (não-cumulatividade). O crédito restante
// permanece no saldo para vendas futuras.
//
// Parâmetros:
//   - company: a empresa vendedora (tem o saldo_credito atualizado in-place).
//   - split: resultado do split da transação de venda.
//
// Retorna o valor do crédito efetivamente usado (em centavos).
//
// Regras:
//   - Se saldo_credito >= imposto total da venda: usa exatamente o valor do imposto,
//     a empresa não desembolsa nada do próprio caixa para o tributo.
//   - Se saldo_credito < imposto total: usa todo o crédito disponível, o restante
//     é desembolsado do caixa (simplificação: não estamos modelando esse
//     desembolso separadamente).
//
// Fonte: simplificação do simulador (não corresponde a regra normativa real).
func UseCredit(company *Company, split SplitResult) (creditoUsado int64) {
	impostoTotal := split.ValorIBS + split.ValorCBS
	creditoUsado = company.SaldoCredito
	if creditoUsado > impostoTotal {
		creditoUsado = impostoTotal
	}
	company.SaldoCredito -= creditoUsado
	return creditoUsado
}
