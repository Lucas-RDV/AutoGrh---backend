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
	"time"

	"github.com/go-chi/chi/v5"
)

type DescansoController struct {
	descansoService *service.DescansoService
}

func NewDescansoController(descansoService *service.DescansoService) *DescansoController {
	return &DescansoController{descansoService: descansoService}
}

// POST /descansos
func (c *DescansoController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	var req struct {
		FeriasID int64     `json:"ferias_id"`
		Inicio   time.Time `json:"inicio"`
		Fim      time.Time `json:"fim"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	d := &entity.Descanso{
		FeriasID: req.FeriasID,
		Inicio:   req.Inicio,
		Fim:      req.Fim,
	}
	if err := c.descansoService.CreateDescanso(r.Context(), claims, d); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusCreated, d)
}

// PUT /descansos/{id}/aprovar  (admin)
func (c *DescansoController) Aprovar(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	if err := c.descansoService.AprovarDescanso(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso aprovado"})
}

// PUT /descansos/{id}/pagar
func (c *DescansoController) Pagar(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	if err := c.descansoService.MarcarComoPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso pago"})
}

// PUT /descansos/{id}/desmarcar-pago  (admin)
func (c *DescansoController) DesmarcarPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	if err := c.descansoService.DesmarcarPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso desmarcado como pago"})
}

// GET /funcionarios/{id}/descansos
func (c *DescansoController) ListByFuncionario(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	list, err := c.descansoService.ListarPorFuncionario(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// GET /ferias/{id}/descansos
func (c *DescansoController) ListByFerias(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	list, err := c.descansoService.ListarPorFerias(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// GET /descansos/aprovados
func (c *DescansoController) ListAprovados(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	list, err := c.descansoService.ListarAprovados(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// GET /descansos/pendentes
func (c *DescansoController) ListPendentes(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	list, err := c.descansoService.ListarPendentes(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, list)
}

// DELETE /descansos/{id}  (reprovar)
func (c *DescansoController) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpjson.BadRequest(w, "id inválido")
		return
	}
	if err := c.descansoService.DeleteDescanso(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso deletado"})
}

func (c *DescansoController) CreateAuto(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	funcID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || funcID <= 0 {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}
	var in struct {
		Inicio string `json:"inicio"` // "YYYY-MM-DD"
		Fim    string `json:"fim"`    // "YYYY-MM-DD"
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}
	if in.Inicio == "" || in.Fim == "" {
		httpjson.BadRequest(w, "campos 'inicio' e 'fim' são obrigatórios")
		return
	}
	ini, err := dateStringToTime.DateStringToTime(in.Inicio)
	if err != nil {
		httpjson.Internal(w, "data 'inicio' inválida: "+err.Error())
		return
	}
	fim, err := dateStringToTime.DateStringToTime(in.Fim)
	if err != nil {
		httpjson.Internal(w, "data 'fim' inválida: "+err.Error())
		return
	}
	if err := c.descansoService.CreateDescansoAuto(r.Context(), claims, funcID, ini, fim); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, "ok")
}
