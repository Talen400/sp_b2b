package pp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecimal18_2_Marshal(t *testing.T) {
	tests := []struct {
		input    Decimal18_2
		expected string
	}{
		{Decimal18_2(0), `"0.00"`},
		{Decimal18_2(1000), `"10.00"`},
		{Decimal18_2(100050), `"1000.50"`},
		{Decimal18_2(1), `"0.01"`},
		{Decimal18_2(123456789), `"1234567.89"`},
		{Decimal18_2(-500), `"-5.00"`},
	}

	for _, tt := range tests {
		b, err := tt.input.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON(%d) error: %v", tt.input, err)
		}
		if string(b) != tt.expected {
			t.Errorf("MarshalJSON(%d) = %s, want %s", tt.input, string(b), tt.expected)
		}
	}
}

func TestDecimal18_2_Unmarshal(t *testing.T) {
	tests := []struct {
		input    string
		expected Decimal18_2
	}{
		{`"0.00"`, Decimal18_2(0)},
		{`"10.00"`, Decimal18_2(1000)},
		{`"1000.50"`, Decimal18_2(100050)},
		{`"0.01"`, Decimal18_2(1)},
		{`"1234567.89"`, Decimal18_2(123456789)},
		{`"-5.00"`, Decimal18_2(-500)},
	}

	for _, tt := range tests {
		var d Decimal18_2
		if err := d.UnmarshalJSON([]byte(tt.input)); err != nil {
			t.Errorf("UnmarshalJSON(%s) error: %v", tt.input, err)
		}
		if d != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %d, want %d", tt.input, d, tt.expected)
		}
	}
}

func TestInformeTransacaoIniciada_Serialization(t *testing.T) {
	req := InformeTransacaoIniciadaRequest{
		NSUId:       "123456",
		CNPJRec:     "11.111.111/0001-11",
		VLInf:       Decimal18_2(100000),
		VLIbs:       Decimal18_2(12000),
		VLCbs:       Decimal18_2(3000),
		DtHrCriacao: "2026-07-23T10:00:00Z",
	}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if !strings.Contains(string(b), `"1000.00"`) {
		t.Errorf("VLInf should be 1000.00, got: %s", string(b))
	}
}

func TestSendInformeTransacaoIniciada_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Message-Id") == "" {
			t.Error("missing Message-Id header")
		}
		if r.Header.Get("Correlation-Id") == "" {
			t.Error("missing Correlation-Id header")
		}
		if r.Header.Get("Tenant-Id") != "psp-teste" {
			t.Errorf("Tenant-Id = %s, want psp-teste", r.Header.Get("Tenant-Id"))
		}
		if r.Header.Get("Timestamp") == "" {
			t.Error("missing Timestamp header")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	client := New(srv.URL, "psp-teste")
	req := &InformeTransacaoIniciadaRequest{
		NSUId: "1", CNPJRec: "11.111.111/0001-11",
		VLInf: 1000, VLIbs: 120, VLCbs: 30,
	}

	status, ppErr, err := client.SendInformeTransacaoIniciada(context.Background(), "boleto", req)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("status = %d, want 201", status)
	}
	if ppErr != nil {
		t.Errorf("unexpected PP error: %v", ppErr)
	}
}

func TestSendInformeTransacaoIniciada_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"type":"https://split.rfb.gov.br/erros/campo-invalido","title":"Campo inválido","status":422,"detail":"vlInf inválido","instance":"/api/v1/boleto"}`))
	}))
	defer srv.Close()

	client := New(srv.URL, "psp-teste")
	req := &InformeTransacaoIniciadaRequest{NSUId: "1", CNPJRec: "11.111.111/0001-11"}
	status, ppErr, err := client.SendInformeTransacaoIniciada(context.Background(), "boleto", req)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}
	if status != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", status)
	}
	if ppErr == nil || ppErr.Title != "Campo inválido" {
		t.Errorf("unexpected PP error: %v", ppErr)
	}
}

func TestLongPolling_204(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("proximoToken", "/api/v1/boleto/psp-1/tributos/stream/token=R1")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := New(srv.URL, "psp-1")
	lpr, ppErr, err := client.StartLongPolling(context.Background(), "boleto", "psp-1")
	if err != nil {
		t.Fatalf("LongPolling error: %v", err)
	}
	if ppErr != nil {
		t.Errorf("unexpected PP error: %v", ppErr)
	}
	if !lpr.NoContent {
		t.Error("expected NoContent=true")
	}
	if lpr.ProximoToken != "/api/v1/boleto/psp-1/tributos/stream/token=R1" {
		t.Errorf("token = %s, want /api/v1/.../token=R1", lpr.ProximoToken)
	}
}
