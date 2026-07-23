package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Talen400/sp_b2b/internal/domain"
)

// createTransactionRequest é o payload esperado pelo POST /api/v1/transactions.
// vendedor_cnpj e comprador_cnpj devem ser CNPJs de empresas já cadastradas.
// valor_bruto é em centavos (int64). aliquotas em decimal (ex: 0.12 = 12%).
type createTransactionRequest struct {
	VendedorCNPJ  string  `json:"vendedor_cnpj"`
	CompradorCNPJ string  `json:"comprador_cnpj"`
	ValorBruto    int64   `json:"valor_bruto"`
	AliquotaIBS   float64 `json:"aliquota_ibs"`
	AliquotaCBS   float64 `json:"aliquota_cbs"`
}

// createTransactionResponse é o retorno do POST /api/v1/transactions.
// Inclui a transação persistida, os valores de crédito usado/gerado, e
// opcionalmente o resultado da notificação à Plataforma Pública.
type createTransactionResponse struct {
	Transaction    domain.Transaction `json:"transaction"`
	CreditoUsado   int64              `json:"credito_usado"`
	CreditoGerado  int64              `json:"credito_gerado"`
	PPNotification *PPNotification    `json:"pp_notification,omitempty"`
}

// CreateTransaction lida com POST /api/v1/transactions.
//
// Fluxo:
//  1. Valida o payload (campos obrigatórios, valor bruto positivo).
//  2. Calcula o split (CalculateSplit).
//  3. Consulta vendedor e comprador no repositório.
//  4. Calcula crédito usado (UseCredit) e gerado (ApplyCredit).
//  5. Persiste a transação e atualiza os saldos de crédito.
//  6. Se houver PP notifier configurado, dispara notificação ao mock.
//
// Responde 201 com a transação criada, ou 400/404/500 em caso de erro.
// O campo pp_notification aparece apenas quando -pp-url está configurado
// (status "sent") ou quando não está (status "skipped").
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "corpo da requisição inválido")
		return
	}

	if req.VendedorCNPJ == "" || req.CompradorCNPJ == "" || req.ValorBruto <= 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR",
			"vendedor_cnpj, comprador_cnpj e valor_bruto (positivo) são obrigatórios")
		return
	}

	split, err := domain.CalculateSplit(req.ValorBruto, req.AliquotaIBS, req.AliquotaCBS)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	vendedor, err := h.companyRepo.Get(req.VendedorCNPJ)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao consultar vendedor")
		return
	}

	comprador, err := h.companyRepo.Get(req.CompradorCNPJ)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao consultar comprador")
		return
	}

	creditoUsado := domain.UseCredit(&vendedor, split)
	creditoGerado := domain.ApplyCredit(&comprador, split)

	id := fmt.Sprintf("TXN-%d", time.Now().UnixNano())

	txn := domain.Transaction{
		ID:            id,
		VendedorCNPJ:  req.VendedorCNPJ,
		CompradorCNPJ: req.CompradorCNPJ,
		ValorBruto:    req.ValorBruto,
		AliquotaIBS:   req.AliquotaIBS,
		AliquotaCBS:   req.AliquotaCBS,
		ValorLiquido:  split.Liquido,
		ValorIBS:      split.ValorIBS,
		ValorCBS:      split.ValorCBS,
		CreditoUsado:  creditoUsado,
		Timestamp:     time.Now(),
	}

	if err := h.transactionRepo.Create(txn); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao salvar transação")
		return
	}
	if err := h.companyRepo.UpdateCredit(req.VendedorCNPJ, vendedor.SaldoCredito); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao atualizar crédito do vendedor")
		return
	}
	if err := h.companyRepo.UpdateCredit(req.CompradorCNPJ, comprador.SaldoCredito); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao atualizar crédito do comprador")
		return
	}

	resp := createTransactionResponse{
		Transaction:   txn,
		CreditoUsado:  creditoUsado,
		CreditoGerado: creditoGerado,
	}

	if h.ppNotify != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		resp.PPNotification = h.ppNotify(ctx, txn)
	} else {
		resp.PPNotification = &PPNotification{Status: "skipped"}
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetTransaction lida com GET /api/v1/transactions/{id}.
// Retorna os detalhes de uma transação pelo ID.
// Responde 200 ou 404.
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id é obrigatório")
		return
	}

	txn, err := h.transactionRepo.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao consultar transação")
		return
	}
	writeJSON(w, http.StatusOK, txn)
}

// ListTransactions lida com GET /api/v1/transactions.
// Aceita filtro opcional por CNPJ (?cnpj=...). Se fornecido, retorna
// transações onde o CNPJ é vendedor ou comprador.
// Sempre retorna um array (pode ser vazio).
func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	txns, err := h.transactionRepo.List(cnpj)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao listar transações")
		return
	}
	if txns == nil {
		txns = []domain.Transaction{}
	}
	writeJSON(w, http.StatusOK, txns)
}
