package Controller

import (
	"encoding/json"
	"net/http"
	"time"

	"AutoGRH/pkg/service"
)

type AuthController struct {
	auth *service.AuthService
}

func NewAuthController(auth *service.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

type loginRequest struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

type loginResponse struct {
	Token     string              `json:"token"`
	ExpiresAt time.Time           `json:"expiresAt"`
	Usuario   service.UserMinimal `json:"usuario"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Senha == "" {
		http.Error(w, "login e senha são obrigatórios", http.StatusBadRequest)
		return
	}
	tok, exp, user, err := c.auth.Login(r.Context(), req.Login, req.Senha)
	if err != nil {
		http.Error(w, "credenciais inválidas", http.StatusUnauthorized)
		return
	}
	resp := loginResponse{Token: tok, ExpiresAt: exp, Usuario: user}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
