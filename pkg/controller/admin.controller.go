package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AdminController struct {
	users *service.UsuarioService
}

func NewAdminController(users *service.UsuarioService) *AdminController {
	return &AdminController{users: users}
}

// GET /admin/usuarios
func (c *AdminController) UsuariosList(w http.ResponseWriter, r *http.Request) {
	list, err := c.users.List(r.Context())
	if err != nil {
		httpjson.Internal(w, "erro ao listar usuários")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// POST /admin/usuarios
func (c *AdminController) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUsuarioInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	reativado, err := c.users.Create(r.Context(), req)
	if err != nil {
		httpjson.BadRequest(w, err.Error())
		return
	}

	mensagem := "usuário criado com sucesso"
	status := http.StatusCreated

	if reativado {
		mensagem = "usuário reativado com sucesso"
		status = http.StatusOK
	}

	httpjson.WriteJSON(w, status, map[string]string{
		"mensagem": mensagem,
	})
}

// PUT /admin/usuarios/{id}
func (c *AdminController) UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}
	var req service.UpdateUsuarioInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}
	if err := c.users.Update(r.Context(), id, req); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"mensagem": "usuário atualizado com sucesso"}`))
}

// DELETE /admin/usuarios/{id}
func (c *AdminController) DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}
	if err := c.users.Delete(r.Context(), id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"mensagem": "usuário desativado com sucesso"})
}
