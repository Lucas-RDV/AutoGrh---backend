package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FeriasController struct {
	feriasService *service.FeriasService
}

func NewFeriasController(s *service.FeriasService) *FeriasController {
	return &FeriasController{feriasService: s}
}

// Listar todas as férias
func (c *FeriasController) ListFerias(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	ferias, err := c.feriasService.ListFerias(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, ferias)
}

// Buscar férias por ID
func (c *FeriasController) GetFeriasByID(w http.ResponseWriter, r *http.Request) {
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

	f, err := c.feriasService.GetFeriasByID(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if f == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Férias não encontradas", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, f)
}

// Buscar férias de um funcionário
func (c *FeriasController) GetFeriasByFuncionarioID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	funcionarioIDStr := chi.URLParam(r, "funcionarioID")
	funcionarioID, err := strconv.ParseInt(funcionarioIDStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	lista, err := c.feriasService.GetFeriasByFuncionarioID(r.Context(), claims, funcionarioID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, lista)
}

// Marcar férias como vencidas
func (c *FeriasController) MarcarComoVencidas(w http.ResponseWriter, r *http.Request) {
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

	if err := c.feriasService.MarcarComoVencidas(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Marcar terço como pago
func (c *FeriasController) MarcarTercoComoPago(w http.ResponseWriter, r *http.Request) {
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

	if err := c.feriasService.MarcarTercoComoPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSaldoFerias retorna os dias e valores restantes das férias
func (c *FeriasController) GetSaldoFerias(w http.ResponseWriter, r *http.Request) {
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

	// Buscar férias pelo ID
	ferias, err := c.feriasService.GetFeriasByID(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if ferias == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Férias não encontradas", Code: "NOT_FOUND"})
		return
	}

	salario := ferias.Valor

	// Calcular saldo
	saldo, err := c.feriasService.CalcularSaldo(ferias, salario)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, saldo)
}
