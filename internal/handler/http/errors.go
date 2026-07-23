package httpapi

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse é o formato padronizado de erro da API.
// Todos os endpoints usam este formato — nunca vazar erro interno cru.
// Códigos HTTP: 400 (validação), 404 (não encontrado), 409 (conflito),
// 500 (erro interno — mensagem genérica).
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail carrega o código máquina e a mensagem legível do erro.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError escreve uma resposta de erro padronizada no ResponseWriter.
// O formato segue o contrato definido em AGENT.md regra 5.
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// writeJSON escreve qualquer valor como JSON com o status HTTP informado.
// É o helper usado por todos os handlers para respostas de sucesso.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
