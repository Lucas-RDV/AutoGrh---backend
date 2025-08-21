package middleware

import (
	"context"
	"net/http"
	"strings"

	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/service"
)

type ctxKey int

const claimsKey ctxKey = 1

func WithClaims(ctx context.Context, c service.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

func GetClaims(ctx context.Context) (service.Claims, bool) {
	c, ok := ctx.Value(claimsKey).(service.Claims)
	return c, ok
}

// Agora compatível com chi.Use()
func RequireAuth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearer(r.Header.Get("Authorization"))
			if token == "" {
				httpjson.Unauthorized(w, "TOKEN_MISSING", "token ausente")
				return
			}
			claims, err := auth.ValidateToken(r.Context(), token)
			if err != nil {
				httpjson.Unauthorized(w, "TOKEN_INVALID", "token inválido")
				return
			}
			ctxWithClaims := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctxWithClaims))
		})
	}
}

// Agora compatível com chi.Use()
func RequirePerm(auth *service.AuthService, perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequireAuth(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				httpjson.Unauthorized(w, "NO_CLAIMS", "sem claims")
				return
			}
			if err := auth.Authorize(r.Context(), claims, perm); err != nil {
				httpjson.Forbidden(w, "não autorizado")
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

func extractBearer(h string) string {
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
