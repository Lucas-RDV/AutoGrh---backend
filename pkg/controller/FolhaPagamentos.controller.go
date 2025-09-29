package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	mw "AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FolhaPagamentoController struct {
	service *service.FolhaPagamentoService
}

func NewFolhaPagamentoController(s *service.FolhaPagamentoService) *FolhaPagamentoController {
	return &FolhaPagamentoController{service: s}
}

// CriarFolhaVale cria manualmente uma folha de VALE
func (c *FolhaPagamentoController) CriarFolhaVale(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Mes int `json:"mes"`
		Ano int `json:"ano"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	folha, err := c.service.CriarFolhaVale(r.Context(), claims, input.Mes, input.Ano)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, folha)
}

// ListarFolhas retorna todas as folhas
func (c *FolhaPagamentoController) ListarFolhas(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	folhas, err := c.service.ListarFolhas(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, folhas)
}

// BuscarFolha retorna uma folha pelo ID
func (c *FolhaPagamentoController) BuscarFolha(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	folha, err := c.service.BuscarFolha(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if folha == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Folha não encontrada", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, folha)
}

// BuscarFolhaPorMesAnoTipo retorna uma folha pelo mês, ano e tipo
func (c *FolhaPagamentoController) BuscarFolhaPorMesAnoTipo(w http.ResponseWriter, r *http.Request) {
	mesStr := chi.URLParam(r, "mes")
	anoStr := chi.URLParam(r, "ano")
	tipo := chi.URLParam(r, "tipo")

	mes, err := strconv.Atoi(mesStr)
	if err != nil {
		httpjson.BadRequest(w, "Mês inválido")
		return
	}
	ano, err := strconv.Atoi(anoStr)
	if err != nil {
		httpjson.BadRequest(w, "Ano inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	folha, err := c.service.BuscarFolhaPorMesAnoTipo(r.Context(), claims, mes, ano, tipo)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if folha == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Folha não encontrada", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, folha)
}

// RecalcularFolha limpa os pagamentos e recalcula
func (c *FolhaPagamentoController) RecalcularFolha(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.service.RecalcularFolha(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FecharFolha marca a folha como paga
func (c *FolhaPagamentoController) FecharFolha(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.service.FecharFolha(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExcluirFolha remove uma folha permanentemente
func (c *FolhaPagamentoController) ExcluirFolha(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.service.ExcluirFolha(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RecalcularFolhaVale limpa os pagamentos e recalcula a folha de VALE
func (c *FolhaPagamentoController) RecalcularFolhaVale(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.service.RecalcularFolhaVale(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
