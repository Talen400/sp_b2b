package domain

func ApplyCredit(company *Company, split SplitResult) (creditoGerado int64) {
	creditoGerado = split.ValorIBS + split.ValorCBS
	company.SaldoCredito += creditoGerado
	return creditoGerado
}

func UseCredit(company *Company, split SplitResult) (creditoUsado int64) {
	impostoTotal := split.ValorIBS + split.ValorCBS
	creditoUsado = company.SaldoCredito
	if creditoUsado > impostoTotal {
		creditoUsado = impostoTotal
	}
	company.SaldoCredito -= creditoUsado
	return creditoUsado
}
