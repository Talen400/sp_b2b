package domain

type Company struct {
	CNPJ         string `json:"cnpj"`
	Nome         string `json:"nome"`
	SaldoCredito int64  `json:"saldo_credito"`
}
