package domain

import "testing"

func TestApplyCredit(t *testing.T) {
	company := &Company{CNPJ: "11.111.111/0001-11", Nome: "Teste", SaldoCredito: 0}
	split := SplitResult{ValorIBS: 12000, ValorCBS: 3000}

	credito := ApplyCredit(company, split)
	if credito != 15000 {
		t.Errorf("creditoGerado esperado 15000, got %d", credito)
	}
	if company.SaldoCredito != 15000 {
		t.Errorf("saldo esperado 15000, got %d", company.SaldoCredito)
	}
}

func TestUseCredit_Suficiente(t *testing.T) {
	company := &Company{CNPJ: "22.222.222/0001-22", Nome: "Teste", SaldoCredito: 15000}
	split := SplitResult{ValorIBS: 36000, ValorCBS: 9000}

	usado := UseCredit(company, split)
	if usado != 15000 {
		t.Errorf("creditoUsado esperado 15000, got %d", usado)
	}
	if company.SaldoCredito != 0 {
		t.Errorf("saldo esperado 0, got %d", company.SaldoCredito)
	}
}

func TestUseCredit_Insuficiente(t *testing.T) {
	company := &Company{CNPJ: "22.222.222/0001-22", Nome: "Teste", SaldoCredito: 5000}
	split := SplitResult{ValorIBS: 36000, ValorCBS: 9000}

	usado := UseCredit(company, split)
	if usado != 5000 {
		t.Errorf("creditoUsado esperado 5000, got %d", usado)
	}
	if company.SaldoCredito != 0 {
		t.Errorf("saldo esperado 0, got %d", company.SaldoCredito)
	}
}

func TestUseCredit_CreditoZero(t *testing.T) {
	company := &Company{CNPJ: "11.111.111/0001-11", Nome: "Teste", SaldoCredito: 0}
	split := SplitResult{ValorIBS: 12000, ValorCBS: 3000}

	usado := UseCredit(company, split)
	if usado != 0 {
		t.Errorf("creditoUsado esperado 0, got %d", usado)
	}
}

func TestUseCredit_CreditoExato(t *testing.T) {
	company := &Company{CNPJ: "22.222.222/0001-22", Nome: "Teste", SaldoCredito: 45000}
	split := SplitResult{ValorIBS: 36000, ValorCBS: 9000}

	usado := UseCredit(company, split)
	if usado != 45000 {
		t.Errorf("creditoUsado esperado 45000, got %d", usado)
	}
	if company.SaldoCredito != 0 {
		t.Errorf("saldo esperado 0, got %d", company.SaldoCredito)
	}
}
