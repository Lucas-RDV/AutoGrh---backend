package controller

import (
	"AutoGRH/pkg/controller/httpjson"
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
		httpjson.BadRequest(w, "Request inválido")
		return
	}
	if req.Login == "" || req.Senha == "" {
		httpjson.BadRequest(w, "login e senha são obrigatórios")
		return
	}
	tok, exp, user, err := c.auth.Login(r.Context(), req.Login, req.Senha)
	if err != nil {
		httpjson.Unauthorized(w, "INVALID_CREDENTIALS", "credenciais inválidas")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Em produção (HTTPS) use Secure:true
		Secure:  false,
		Expires: exp,
	})
	// corpo sem expor o token (front não precisa mais dele)
	resp := loginResponse{ExpiresAt: exp, Usuario: user}
	httpjson.WriteJSON(w, http.StatusOK, resp)
}

// POST /auth/logout — apaga o cookie httpOnly
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	httpjson.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
