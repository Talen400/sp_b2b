package httpapi

import "net/http"

// Healthz lida com GET /api/v1/healthz.
// Responde 200 com {"status":"ok"} — usado para liveness check.
// Não depende de banco ou qualquer recurso externo.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
