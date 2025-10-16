package controller

import (
	"AutoGRH/pkg/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type LogController struct{}

func NewLogController() *LogController { return &LogController{} }

// GET /admin/logs?limit=200&usuarioId=123
func (c *LogController) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit := 200
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	if uidStr := q.Get("usuarioId"); uidStr != "" {
		if uid, err := strconv.ParseInt(uidStr, 10, 64); err == nil && uid > 0 {
			logs, err := repository.GetLogsByUsuarioID(uid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(logs)
			return
		}
	}

	logs, err := repository.ListAllLogsView(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(logs)
}
