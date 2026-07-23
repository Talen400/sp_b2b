// Pacote httpapi implementa a camada HTTP da API REST.
// Cada handler é um método de Handler, que recebe as interfaces de
// repositório e um callback opcional para notificar a Plataforma Pública.
//
// Nenhuma regra de negócio vaza para dentro destes handlers — eles
// decodificam o request, chamam domínio/repositório, e codificam a resposta.
package httpapi

import (
	"context"
	"net/http"

	"github.com/Talen400/sp_b2b/internal/domain"
	"github.com/Talen400/sp_b2b/internal/repository"
)

// PPNotification é o resultado da notificação opcional à Plataforma Pública.
// Aparece no response do POST /transactions quando -pp-url está configurado.
type PPNotification struct {
	Status          string `json:"status"`                      // "skipped", "sent" ou "failed"
	Arrangement     string `json:"arrangement,omitempty"`       // arranjo notificado (ex: "boleto")
	PPMessageID     string `json:"pp_message_id,omitempty"`     // Message-Id enviado à PP
	PPCorrelationID string `json:"pp_correlation_id,omitempty"` // Correlation-Id enviado à PP
	Error           string `json:"error,omitempty"`             // detalhe do erro se status = "failed"
}

// PPNotifyFunc é o callback disparado após criar uma transação com sucesso.
// A implementação é injetada pelo main.go a partir da flag -pp-url.
// Deve respeitar o contexto (timeout de 5s) para não travar a resposta.
type PPNotifyFunc func(ctx context.Context, txn domain.Transaction) *PPNotification

// Handler agrupa os dependências compartilhadas por todos os handlers HTTP.
type Handler struct {
	companyRepo     repository.CompanyRepo
	transactionRepo repository.TransactionRepo
	ppNotify        PPNotifyFunc
}

// NewHandler cria um Handler sem notificação PP.
func NewHandler(cr repository.CompanyRepo, tr repository.TransactionRepo) *Handler {
	return &Handler{
		companyRepo:     cr,
		transactionRepo: tr,
	}
}

// NewHandlerWithPP cria um Handler com callback de notificação PP.
// O callback é disparado após cada POST /transactions bem-sucedido.
func NewHandlerWithPP(cr repository.CompanyRepo, tr repository.TransactionRepo, ppn PPNotifyFunc) *Handler {
	return &Handler{
		companyRepo:     cr,
		transactionRepo: tr,
		ppNotify:        ppn,
	}
}

// Router monta o http.ServeMux com todas as rotas da API.
// Usa o roteamento por método+path do Go 1.22+ (ServeMux nativo).
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
