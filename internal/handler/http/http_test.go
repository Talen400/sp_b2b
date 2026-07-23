package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Talen400/sp_b2b/internal/repository/sqlite"
)

func setupTestHandler(t *testing.T) *Handler {
	t.Helper()
	f, err := os.CreateTemp("", "split-http-test-*.db")
	if err != nil {
		t.Fatalf("criar temp db: %v", err)
	}
	f.Close()

	db, err := sqlite.Open(f.Name())
	if err != nil {
		t.Fatalf("abrir db: %v", err)
	}
	t.Cleanup(func() { db.Close(); os.Remove(f.Name()) })

	if err := db.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return NewHandler(
		sqlite.NewCompanyRepository(db),
		sqlite.NewTransactionRepository(db),
	)
}

func TestHealthz(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperava 200, got %d", resp.StatusCode)
	}
}

func TestCreateCompany_HappyPath(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	body := `{"cnpj":"11.111.111/0001-11","nome":"Fazenda Boa Vista"}`
	resp, err := http.Post(srv.URL+"/api/v1/companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperava 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["cnpj"] != "11.111.111/0001-11" {
		t.Errorf("CNPJ inesperado: %v", result["cnpj"])
	}
}

func TestCreateCompany_Duplicate(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	body := `{"cnpj":"11.111.111/0001-11","nome":"Fazenda"}`
	http.Post(srv.URL+"/api/v1/companies", "application/json", strings.NewReader(body))

	resp, err := http.Post(srv.URL+"/api/v1/companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperava 409, got %d", resp.StatusCode)
	}
}

func TestGetCompany_NotFound(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/companies/00.000.000/0000-00")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("esperava 404, got %d", resp.StatusCode)
	}
}

func TestCreateTransaction_HappyPath(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"11.111.111/0001-11","nome":"Fazenda"}`))
	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"22.222.222/0001-22","nome":"Fábrica"}`))

	body := `{"vendedor_cnpj":"11.111.111/0001-11","comprador_cnpj":"22.222.222/0001-22","valor_bruto":100000,"aliquota_ibs":0.12,"aliquota_cbs":0.03}`
	resp, err := http.Post(srv.URL+"/api/v1/transactions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperava 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	txn, ok := result["transaction"].(map[string]interface{})
	if !ok {
		t.Fatal("resposta sem campo transaction")
	}
	if txn["valor_bruto"].(float64) != 100000 {
		t.Errorf("valor_bruto esperado 100000, got %v", txn["valor_bruto"])
	}
	if txn["credito_usado"].(float64) != 0 {
		t.Errorf("credito_usado esperado 0, got %v", txn["credito_usado"])
	}
}

func TestCreateTransaction_SellerNotFound(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"22.222.222/0001-22","nome":"Fábrica"}`))

	body := `{"vendedor_cnpj":"11.111.111/0001-11","comprador_cnpj":"22.222.222/0001-22","valor_bruto":100000}`
	resp, err := http.Post(srv.URL+"/api/v1/transactions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("esperava 404, got %d", resp.StatusCode)
	}
}

func TestListTransactions(t *testing.T) {
	h := setupTestHandler(t)
	srv := httptest.NewServer(h.Router())
	defer srv.Close()

	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"11.111.111/0001-11","nome":"Fazenda"}`))
	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"22.222.222/0001-22","nome":"Fábrica"}`))
	http.Post(srv.URL+"/api/v1/companies", "application/json",
		strings.NewReader(`{"cnpj":"33.333.333/0001-33","nome":"Mercado"}`))

	http.Post(srv.URL+"/api/v1/transactions", "application/json",
		strings.NewReader(`{"vendedor_cnpj":"11.111.111/0001-11","comprador_cnpj":"22.222.222/0001-22","valor_bruto":100000,"aliquota_ibs":0.12,"aliquota_cbs":0.03}`))
	http.Post(srv.URL+"/api/v1/transactions", "application/json",
		strings.NewReader(`{"vendedor_cnpj":"22.222.222/0001-22","comprador_cnpj":"33.333.333/0001-33","valor_bruto":300000,"aliquota_ibs":0.12,"aliquota_cbs":0.03}`))

	resp, err := http.Get(srv.URL + "/api/v1/transactions")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperava 200, got %d", resp.StatusCode)
	}

	var txns []interface{}
	json.NewDecoder(resp.Body).Decode(&txns)
	if len(txns) != 2 {
		t.Errorf("esperava 2 transações, got %d", len(txns))
	}
}
