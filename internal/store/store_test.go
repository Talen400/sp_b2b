package store

import (
	"testing"
)

func setupStore(t *testing.T) *Store {
	t.Helper()
	s := NewStore()
	if err := s.AddCompany("11.111.111/0001-11", "Fazenda Boa Vista"); err != nil {
		t.Fatalf("erro ao cadastrar empresa: %v", err)
	}
	if err := s.AddCompany("22.222.222/0001-22", "Fábrica de Sucos SA"); err != nil {
		t.Fatalf("erro ao cadastrar empresa: %v", err)
	}
	if err := s.AddCompany("33.333.333/0001-33", "Mercado Central"); err != nil {
		t.Fatalf("erro ao cadastrar empresa: %v", err)
	}
	return s
}

func TestAddCompany_Duplicate(t *testing.T) {
	s := NewStore()
	if err := s.AddCompany("11.111.111/0001-11", "Empresa A"); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	err := s.AddCompany("11.111.111/0001-11", "Empresa A duplicada")
	if err == nil {
		t.Fatal("esperava erro para CNPJ duplicado")
	}
}

func TestGetCompany_NotFound(t *testing.T) {
	s := NewStore()
	_, err := s.GetCompany("00.000.000/0000-00")
	if err == nil {
		t.Fatal("esperava erro para empresa inexistente")
	}
}

func TestGetBalance(t *testing.T) {
	s := setupStore(t)
	balance, err := s.GetBalance("11.111.111/0001-11")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if balance != 0 {
		t.Errorf("saldo inicial esperado 0, got %d", balance)
	}
}

func TestProcessTransaction_WithoutCredit(t *testing.T) {
	s := setupStore(t)

	result, err := s.ProcessTransaction(
		"11.111.111/0001-11", // vendedor: Fazenda
		"22.222.222/0001-22", // comprador: Fábrica
		100000,               // R$ 1.000,00 em centavos
		0.12,                 // IBS 12%
		0.03,                 // CBS 3%
	)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	// Verifica split
	if result.Transaction.ValorIBS != 12000 {
		t.Errorf("ValorIBS esperado 12000, got %d", result.Transaction.ValorIBS)
	}
	if result.Transaction.ValorCBS != 3000 {
		t.Errorf("ValorCBS esperado 3000, got %d", result.Transaction.ValorCBS)
	}
	if result.Transaction.Liquido != 85000 {
		t.Errorf("Liquido esperado 85000, got %d", result.Transaction.Liquido)
	}

	// Vendedor não tinha crédito, então creditoUsado = 0
	if result.Transaction.CreditoUsado != 0 {
		t.Errorf("CreditoUsado esperado 0, got %d", result.Transaction.CreditoUsado)
	}

	// Comprador ganhou crédito de 15000 (12000 + 3000)
	if result.CreditoGerado != 15000 {
		t.Errorf("CreditoGerado esperado 15000, got %d", result.CreditoGerado)
	}

	// Verifica saldos
	saldoVendedor, _ := s.GetBalance("11.111.111/0001-11")
	if saldoVendedor != 0 {
		t.Errorf("saldo vendedor esperado 0, got %d", saldoVendedor)
	}

	saldoComprador, _ := s.GetBalance("22.222.222/0001-22")
	if saldoComprador != 15000 {
		t.Errorf("saldo comprador esperado 15000, got %d", saldoComprador)
	}
}

func TestProcessTransaction_WithCreditUsage(t *testing.T) {
	s := setupStore(t)

	// Venda 1: Fazenda → Fábrica, R$ 1.000
	_, err := s.ProcessTransaction("11.111.111/0001-11", "22.222.222/0001-22", 100000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro na venda 1: %v", err)
	}

	// Fábrica agora tem 15000 de crédito
	saldo, _ := s.GetBalance("22.222.222/0001-22")
	if saldo != 15000 {
		t.Fatalf("saldo Fábrica esperado 15000, got %d", saldo)
	}

	// Venda 2: Fábrica → Mercado, R$ 3.000 (IBS 12% = 36000, CBS 3% = 9000, total = 45000)
	result, err := s.ProcessTransaction("22.222.222/0001-22", "33.333.333/0001-33", 300000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro na venda 2: %v", err)
	}

	// Fábrica usou 15000 de crédito (todo o saldo dela)
	if result.Transaction.CreditoUsado != 15000 {
		t.Errorf("CreditoUsado esperado 15000, got %d", result.Transaction.CreditoUsado)
	}

	// Mercado ganhou crédito de 45000 (full IBS+CBS)
	if result.CreditoGerado != 45000 {
		t.Errorf("CreditoGerado esperado 45000, got %d", result.CreditoGerado)
	}

	// Fábrica ficou com saldo 0 (usou tudo)
	saldoFabrica, _ := s.GetBalance("22.222.222/0001-22")
	if saldoFabrica != 0 {
		t.Errorf("saldo Fábrica esperado 0, got %d", saldoFabrica)
	}

	// Mercado ficou com 45000
	saldoMercado, _ := s.GetBalance("33.333.333/0001-33")
	if saldoMercado != 45000 {
		t.Errorf("saldo Mercado esperado 45000, got %d", saldoMercado)
	}
}

func TestProcessTransaction_PartialCreditUsage(t *testing.T) {
	s := setupStore(t)

	// Gera 15000 de crédito para Fábrica
	_, err := s.ProcessTransaction("11.111.111/0001-11", "22.222.222/0001-22", 100000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro na venda 1: %v", err)
	}

	// Venda pequena: Fábrica → Mercado, R$ 500 (IBS 12% = 6000, CBS 3% = 1500, total = 7500)
	result, err := s.ProcessTransaction("22.222.222/0001-22", "33.333.333/0001-33", 50000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro na venda 2: %v", err)
	}

	// CreditoUsado = min(15000, 7500) = 7500
	if result.Transaction.CreditoUsado != 7500 {
		t.Errorf("CreditoUsado esperado 7500, got %d", result.Transaction.CreditoUsado)
	}

	// Fábrica ainda tem 15000 - 7500 = 7500 de crédito
	saldo, _ := s.GetBalance("22.222.222/0001-22")
	if saldo != 7500 {
		t.Errorf("saldo Fábrica esperado 7500, got %d", saldo)
	}
}

func TestProcessTransaction_CompanyNotFound(t *testing.T) {
	s := NewStore()
	_, err := s.ProcessTransaction("00.000.000/0000-00", "11.111.111/0001-11", 100000, 0.12, 0.03)
	if err == nil {
		t.Fatal("esperava erro para vendedor inexistente")
	}
}

func TestListTransactions(t *testing.T) {
	s := setupStore(t)

	_, err := s.ProcessTransaction("11.111.111/0001-11", "22.222.222/0001-22", 100000, 0.12, 0.03)
	if err != nil {
		t.Fatalf("erro na venda 1: %v", err)
	}

	txns := s.ListTransactions()
	if len(txns) != 1 {
		t.Fatalf("esperava 1 transação, got %d", len(txns))
	}

	if txns[0].VendedorCNPJ != "11.111.111/0001-11" {
		t.Errorf("VendedorCNPJ esperado 11.111.111/0001-11, got %s", txns[0].VendedorCNPJ)
	}
	if txns[0].CompradorCNPJ != "22.222.222/0001-22" {
		t.Errorf("CompradorCNPJ esperado 22.222.222/0001-22, got %s", txns[0].CompradorCNPJ)
	}
}
