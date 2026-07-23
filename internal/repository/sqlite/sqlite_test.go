package sqlite

import (
	"os"
	"testing"
	"time"

	"github.com/Talen400/sp_b2b/internal/domain"
)

func newTestDB(t *testing.T) *DB {
	t.Helper()
	f, err := os.CreateTemp("", "split-test-*.db")
	if err != nil {
		t.Fatalf("criar temp db: %v", err)
	}
	f.Close()

	db, err := Open(f.Name())
	if err != nil {
		t.Fatalf("abrir db: %v", err)
	}
	t.Cleanup(func() { db.Close(); os.Remove(f.Name()) })

	if err := db.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestCreateCompany(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	err := r.Create("11.111.111/0001-11", "Fazenda Boa Vista")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestCreateCompany_Duplicate(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	r.Create("11.111.111/0001-11", "Fazenda Boa Vista")
	err := r.Create("11.111.111/0001-11", "Fazenda Boa Vista")
	if err == nil {
		t.Fatal("esperava erro para CNPJ duplicado")
	}
}

func TestGetCompany(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	r.Create("11.111.111/0001-11", "Fazenda Boa Vista")

	c, err := r.Get("11.111.111/0001-11")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if c.Nome != "Fazenda Boa Vista" || c.SaldoCredito != 0 {
		t.Errorf("resultado inesperado: %+v", c)
	}
}

func TestGetCompany_NotFound(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	_, err := r.Get("00.000.000/0000-00")
	if err == nil {
		t.Fatal("esperava erro para CNPJ inexistente")
	}
}

func TestListCompanies(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	r.Create("22.222.222/0001-22", "Fábrica de Sucos SA")
	r.Create("11.111.111/0001-11", "Fazenda Boa Vista")

	companies, err := r.List()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(companies) != 2 {
		t.Fatalf("esperava 2 empresas, got %d", len(companies))
	}
	if companies[0].Nome != "Fazenda Boa Vista" {
		t.Errorf("primeira deveria ser Fazenda (ordem alfabética), got %s", companies[0].Nome)
	}
}

func TestUpdateCompanyCredit(t *testing.T) {
	db := newTestDB(t)
	r := NewCompanyRepository(db)
	r.Create("11.111.111/0001-11", "Fazenda Boa Vista")

	err := r.UpdateCredit("11.111.111/0001-11", 15000)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	c, _ := r.Get("11.111.111/0001-11")
	if c.SaldoCredito != 15000 {
		t.Errorf("saldo esperado 15000, got %d", c.SaldoCredito)
	}
}

func TestCreateAndGetTransaction(t *testing.T) {
	db := newTestDB(t)
	cr := NewCompanyRepository(db)
	cr.Create("11.111.111/0001-11", "Fazenda")
	cr.Create("22.222.222/0001-22", "Fábrica")

	tr := NewTransactionRepository(db)
	txn := domain.Transaction{
		ID:            "TXN-00001",
		VendedorCNPJ:  "11.111.111/0001-11",
		CompradorCNPJ: "22.222.222/0001-22",
		ValorBruto:    100000,
		AliquotaIBS:   0.12,
		AliquotaCBS:   0.03,
		ValorLiquido:  85000,
		ValorIBS:      12000,
		ValorCBS:      3000,
		CreditoUsado:  0,
		Timestamp:     time.Now(),
	}

	err := tr.Create(txn)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	got, err := tr.Get("TXN-00001")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got.ValorBruto != 100000 || got.ValorLiquido != 85000 {
		t.Errorf("resultado inesperado: %+v", got)
	}
}

func TestListTransactions_WithFilter(t *testing.T) {
	db := newTestDB(t)
	cr := NewCompanyRepository(db)
	cr.Create("11.111.111/0001-11", "Fazenda")
	cr.Create("22.222.222/0001-22", "Fábrica")
	cr.Create("33.333.333/0001-33", "Mercado")

	tr := NewTransactionRepository(db)
	t1 := domain.Transaction{ID: "TXN-00001", VendedorCNPJ: "11.111.111/0001-11", CompradorCNPJ: "22.222.222/0001-22", ValorBruto: 100000, AliquotaIBS: 0.12, AliquotaCBS: 0.03, ValorLiquido: 85000, ValorIBS: 12000, ValorCBS: 3000, CreditoUsado: 0, Timestamp: time.Now()}
	t2 := domain.Transaction{ID: "TXN-00002", VendedorCNPJ: "22.222.222/0001-22", CompradorCNPJ: "33.333.333/0001-33", ValorBruto: 300000, AliquotaIBS: 0.12, AliquotaCBS: 0.03, ValorLiquido: 255000, ValorIBS: 36000, ValorCBS: 9000, CreditoUsado: 15000, Timestamp: time.Now()}

	tr.Create(t1)
	tr.Create(t2)

	txns, err := tr.List("22.222.222/0001-22")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(txns) != 2 {
		t.Fatalf("esperava 2 transações, got %d", len(txns))
	}
}
