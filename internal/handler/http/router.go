package httpapi

import (
	"net/http"

	"github.com/Talen400/sp_b2b/internal/repository"
)

type Handler struct {
	companyRepo     repository.CompanyRepo
	transactionRepo repository.TransactionRepo
}

func NewHandler(cr repository.CompanyRepo, tr repository.TransactionRepo) *Handler {
	return &Handler{
		companyRepo:     cr,
		transactionRepo: tr,
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
