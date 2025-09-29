package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/utils/dateStringToTime"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ValeController struct {
	valeService *service.ValeService
}

func NewValeController(valeService *service.ValeService) *ValeController {
	return &ValeController{valeService: valeService}
}

// CriarVale (RH)
func (c *ValeController) CriarVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	var input struct {
		FuncionarioID int64   `json:"funcionarioID"`
		Valor         float64 `json:"valor"`
		Data          string  `json:"data"` // formato YYYY-MM-DD
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	data, err := dateStringToTime.DateStringToTime(input.Data)
	if err != nil {
		httpjson.BadRequest(w, "Data inválida")
		return
	}

	v, err := c.valeService.CriarVale(r.Context(), claims, input.FuncionarioID, input.Valor, data)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, v)
}

// GetVale por ID
func (c *ValeController) GetVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	v, err := c.valeService.GetVale(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if v == nil {
		httpjson.BadRequest(w, "Vale não encontrado")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, v)
}

// ListarVales (todos)
func (c *ValeController) ListarVales(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	vales, err := c.valeService.ListarVales(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, vales)
}

// ListarValesFuncionario
func (c *ValeController) ListarValesFuncionario(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	funcionarioID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	vales, err := c.valeService.ListarValesFuncionario(r.Context(), claims, funcionarioID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, vales)
}

// ListarValesPendentes
func (c *ValeController) ListarValesPendentes(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	vales, err := c.valeService.ListarValesPendentes(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, vales)
}

// ListarValesAprovadosNaoPagos
func (c *ValeController) ListarValesAprovadosNaoPagos(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	vales, err := c.valeService.ListarValesAprovadosNaoPagos(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, vales)
}

// AtualizarVale
func (c *ValeController) AtualizarVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var v entity.Vale
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}
	v.ID = id

	if err := c.valeService.AtualizarVale(r.Context(), claims, &v); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Vale atualizado com sucesso"})
}

// SoftDeleteVale
func (c *ValeController) SoftDeleteVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := c.valeService.SoftDeleteVale(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Vale inativado"})
}

// AprovarVale
func (c *ValeController) AprovarVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := c.valeService.AprovarVale(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Vale aprovado"})
}

// MarcarValeComoPago
func (c *ValeController) MarcarValeComoPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := c.valeService.MarcarValeComoPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Vale marcado como pago"})
}

// DeleteVale (permanente)
func (c *ValeController) DeleteVale(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := c.valeService.DeleteVale(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "Vale excluído permanentemente"})
}
