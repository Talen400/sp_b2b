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

type createTransactionRequest struct {
	VendedorCNPJ  string  `json:"vendedor_cnpj"`
	CompradorCNPJ string  `json:"comprador_cnpj"`
	ValorBruto    int64   `json:"valor_bruto"`
	AliquotaIBS   float64 `json:"aliquota_ibs"`
	AliquotaCBS   float64 `json:"aliquota_cbs"`
}

type createTransactionResponse struct {
	Transaction    domain.Transaction `json:"transaction"`
	CreditoUsado   int64              `json:"credito_usado"`
	CreditoGerado  int64              `json:"credito_gerado"`
	PPNotification *PPNotification    `json:"pp_notification,omitempty"`
}

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
