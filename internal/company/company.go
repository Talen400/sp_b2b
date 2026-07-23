package company

type Company struct {
	CNPJ         string
	Nome         string
	SaldoCredito int64
}

func New(cnpj, nome string) *Company {
	return &Company{
		CNPJ:         cnpj,
		Nome:         nome,
		SaldoCredito: 0,
	}
}
