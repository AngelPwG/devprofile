package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs, err := h.db.GetAuditLogs()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(w).Encode(logs)
}
