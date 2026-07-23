package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Talen400/sp_b2b/internal/domain"
)

type createCompanyRequest struct {
	CNPJ string `json:"cnpj"`
	Nome string `json:"nome"`
}

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
