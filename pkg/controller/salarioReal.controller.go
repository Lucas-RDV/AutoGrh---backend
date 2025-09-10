package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SalarioRealController struct {
	salarioRealService *service.SalarioRealService
}

func NewSalarioRealController(salarioRealService *service.SalarioRealService) *SalarioRealController {
	return &SalarioRealController{salarioRealService: salarioRealService}
}

// POST /funcionarios/{id}/salarios-reais
func (c *SalarioRealController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	funcID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	var req struct {
		Valor float64 `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	// ⬇️ agora capturamos os DOIS retornos: criado, err
	criado, err := c.salarioRealService.CriarSalarioReal(r.Context(), claims, funcID, req.Valor)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, criado)
}

// GET /funcionarios/{id}/salarios-reais
func (c *SalarioRealController) ListByFuncionario(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	funcID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	list, err := c.salarioRealService.ListSalariosReais(r.Context(), claims, funcID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}

// GET /funcionarios/{id}/salario-real-atual
func (c *SalarioRealController) GetAtual(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	funcID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	atual, err := c.salarioRealService.GetSalarioRealAtual(r.Context(), claims, funcID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if atual == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Nenhum salário real vigente encontrado", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, atual)
}

// DELETE /salarios-reais/{id}
func (c *SalarioRealController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.salarioRealService.DeleteSalarioReal(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "salário real deletado"})
}
