package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"net/http"
)

type AvisoController struct {
	svc *service.AvisoService
}

func NewAvisoController(s *service.AvisoService) *AvisoController {
	return &AvisoController{svc: s}
}

func (c *AvisoController) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}
	avisos, err := c.svc.List(r.Context(), claims)
	if err != nil {
		httpjson.Forbidden(w, "não autorizado")
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, avisos)
}
