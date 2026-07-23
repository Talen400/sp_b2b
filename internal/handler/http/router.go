package httpapi

import (
	"context"
	"net/http"

	"github.com/Talen400/sp_b2b/internal/domain"
	"github.com/Talen400/sp_b2b/internal/repository"
)

// PPNotification is the result of notifying the Plataforma Pública mock.
type PPNotification struct {
	Status          string `json:"status"`                      // "skipped", "sent", "failed"
	Arrangement     string `json:"arrangement,omitempty"`       // e.g. "boleto"
	PPMessageID     string `json:"pp_message_id,omitempty"`     // Message-Id sent
	PPCorrelationID string `json:"pp_correlation_id,omitempty"` // Correlation-Id sent
	Error           string `json:"error,omitempty"`             // error detail if failed
}

// PPNotifyFunc is called after a transaction is created to optionally notify the PP mock.
// The function should not block the response — implementations should use a short timeout.
type PPNotifyFunc func(ctx context.Context, txn domain.Transaction) *PPNotification

type Handler struct {
	companyRepo     repository.CompanyRepo
	transactionRepo repository.TransactionRepo
	ppNotify        PPNotifyFunc
}

func NewHandler(cr repository.CompanyRepo, tr repository.TransactionRepo) *Handler {
	return &Handler{
		companyRepo:     cr,
		transactionRepo: tr,
	}
}

func NewHandlerWithPP(cr repository.CompanyRepo, tr repository.TransactionRepo, ppn PPNotifyFunc) *Handler {
	return &Handler{
		companyRepo:     cr,
		transactionRepo: tr,
		ppNotify:        ppn,
	}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/healthz", h.Healthz)
	mux.HandleFunc("POST /api/v1/companies", h.CreateCompany)
	mux.HandleFunc("GET /api/v1/companies/{cnpj}", h.GetCompany)
	mux.HandleFunc("GET /api/v1/companies", h.ListCompanies)
	mux.HandleFunc("POST /api/v1/transactions", h.CreateTransaction)
	mux.HandleFunc("GET /api/v1/transactions/{id}", h.GetTransaction)
	mux.HandleFunc("GET /api/v1/transactions", h.ListTransactions)

	return mux
}
