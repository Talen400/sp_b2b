package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Talen400/sp_b2b/internal/domain"
)

// createCompanyRequest é o payload esperado pelo POST /api/v1/companies.
type createCompanyRequest struct {
	CNPJ string `json:"cnpj"` // CNPJ da empresa (string livre)
	Nome string `json:"nome"` // Nome fantasia
}

// CreateCompany lida com POST /api/v1/companies.
// Cria uma nova empresa com saldo de crédito inicial zero.
// Responde 201 com a empresa criada, ou 400/409 em caso de erro.
func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req createCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "corpo da requisição inválido")
		return
	}
	if req.CNPJ == "" || req.Nome == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "cnpj e nome são obrigatórios")
		return
	}

	if err := h.companyRepo.Create(req.CNPJ, req.Nome); err != nil {
		code := "INTERNAL_ERROR"
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "já existe") {
			code = "CONFLICT"
			status = http.StatusConflict
		}
		writeError(w, status, code, err.Error())
		return
	}

	company, _ := h.companyRepo.Get(req.CNPJ)
	writeJSON(w, http.StatusCreated, company)
}

// GetCompany lida com GET /api/v1/companies/{cnpj}.
// Retorna os dados da empresa + saldo de crédito atual.
// Responde 200 ou 404.
func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	cnpj := r.PathValue("cnpj")
	if cnpj == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "cnpj é obrigatório")
		return
	}

	company, err := h.companyRepo.Get(cnpj)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao consultar empresa")
		return
	}
	writeJSON(w, http.StatusOK, company)
}

// ListCompanies lida com GET /api/v1/companies.
// Retorna a lista de todas as empresas cadastradas, ordenadas por nome.
// Sempre retorna um array (pode ser vazio).
func (h *Handler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companyRepo.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "erro ao listar empresas")
		return
	}
	if companies == nil {
		companies = []domain.Company{}
	}
	writeJSON(w, http.StatusOK, companies)
}
