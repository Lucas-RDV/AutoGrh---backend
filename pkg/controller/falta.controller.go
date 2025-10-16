package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FaltaController struct {
	faltaService *service.FaltaService
}

func NewFaltaController(s *service.FaltaService) *FaltaController {
	return &FaltaController{faltaService: s}
}

// Criar nova falta
func (c *FaltaController) CreateFalta(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	var f entity.Falta
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	if err := c.faltaService.CreateFalta(r.Context(), claims, &f); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, f)
}

// Atualizar falta
func (c *FaltaController) UpdateFalta(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	var f entity.Falta
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}
	f.ID = id

	if err := c.faltaService.UpdateFalta(r.Context(), claims, &f); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, f)
}

// Deletar falta
func (c *FaltaController) DeleteFalta(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	if err := c.faltaService.DeleteFalta(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "falta deletada"})
}

// Buscar falta por ID
func (c *FaltaController) GetFaltaByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	f, err := c.faltaService.GetFaltaByID(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, f)
}

// Listar todas as faltas de um funcionário
func (c *FaltaController) GetFaltasByFuncionarioID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	funcionarioIDStr := chi.URLParam(r, "id")
	funcionarioID, err := strconv.ParseInt(funcionarioIDStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	lista, err := c.faltaService.GetFaltasByFuncionarioID(r.Context(), claims, funcionarioID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, lista)
}

// Listar todas as faltas
func (c *FaltaController) ListAllFaltas(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	lista, err := c.faltaService.ListAllFaltas(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, lista)
}

func (c *FaltaController) UpsertMensal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	funcionarioID, _ := strconv.ParseInt(idStr, 10, 64)

	var q struct {
		Mes        int `json:"mes"`
		Ano        int `json:"ano"`
		Quantidade int `json:"quantidade"`
	}
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.faltaService.UpsertMensal(r.Context(), claims, funcionarioID, q.Mes, q.Ano, q.Quantidade); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent) // sem corpo
}
