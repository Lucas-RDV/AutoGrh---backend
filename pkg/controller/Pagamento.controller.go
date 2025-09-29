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

type PagamentoController struct {
	pagamentoService *service.PagamentoService
}

func NewPagamentoController(p *service.PagamentoService) *PagamentoController {
	return &PagamentoController{pagamentoService: p}
}

// GET /pagamentos/{id}
func (c *PagamentoController) GetPagamentoByID(w http.ResponseWriter, r *http.Request) {
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

	pagamento, err := c.pagamentoService.BuscarPagamento(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if pagamento == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Pagamento não encontrado", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, pagamento)
}

// GET /funcionarios/{id}/pagamentos
func (c *PagamentoController) ListarPagamentosFuncionario(w http.ResponseWriter, r *http.Request) {
	funcIDStr := chi.URLParam(r, "id")
	funcID, err := strconv.ParseInt(funcIDStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID de funcionário inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	list, err := c.pagamentoService.ListarPagamentosFuncionario(r.Context(), claims, funcID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}

// PUT /pagamentos/{id}
func (c *PagamentoController) UpdatePagamento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	var input struct {
		Adicional float64 `json:"adicional"`
		INSS      float64 `json:"inss"`
		Familia   float64 `json:"familia"`
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

	if err := c.pagamentoService.AtualizarPagamento(r.Context(), claims, id, input.Adicional, input.INSS, input.Familia); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Pagamento atualizado com sucesso"})
}

// PUT /pagamentos/{id}/pagar
func (c *PagamentoController) MarcarComoPago(w http.ResponseWriter, r *http.Request) {
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

	if err := c.pagamentoService.MarcarPagamentoComoPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Pagamento marcado como pago"})
}
