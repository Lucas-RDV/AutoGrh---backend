package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
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

	// DTO só com os campos que o cliente deve enviar
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
		// Valor será calculado no service
	}

	if err := c.descansoService.CreateDescanso(r.Context(), claims, d); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, d)
}

// PUT /descansos/{id}/aprovar
func (c *DescansoController) Aprovar(w http.ResponseWriter, r *http.Request) {
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

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	if err := c.descansoService.MarcarComoPago(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso pago"})
}

// GET /funcionarios/{id}/descansos
func (c *DescansoController) ListByFuncionario(w http.ResponseWriter, r *http.Request) {
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

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
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

// DELETE /descansos/{id}
func (c *DescansoController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.descansoService.DeleteDescanso(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"message": "descanso deletado"})
}
