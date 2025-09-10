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

type SalarioController struct {
	salarioService *service.SalarioService
}

func NewSalarioController(salarioService *service.SalarioService) *SalarioController {
	return &SalarioController{salarioService: salarioService}
}

// POST /funcionarios/{id}/salarios
func (c *SalarioController) Create(w http.ResponseWriter, r *http.Request) {
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

	salario, err := c.salarioService.CriarSalario(r.Context(), claims, funcID, req.Valor)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, salario)
}

// GET /funcionarios/{id}/salarios
func (c *SalarioController) ListByFuncionario(w http.ResponseWriter, r *http.Request) {
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

	list, err := c.salarioService.ListSalarios(r.Context(), claims, funcID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}

// PUT /salarios/{id}
func (c *SalarioController) Update(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Valor float64 `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	salario := &entity.Salario{
		ID:    id,
		Valor: req.Valor,
	}

	if err := c.salarioService.AtualizarSalario(r.Context(), claims, salario); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "salário atualizado"})
}

// DELETE /salarios/{id}
func (c *SalarioController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.salarioService.DeletarSalario(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "salário deletado"})
}
