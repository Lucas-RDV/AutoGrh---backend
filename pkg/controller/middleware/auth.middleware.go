package middleware

import (
	"context"
	"net/http"
	"strings"

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

func RequireAuth(auth *service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearer(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "token ausente", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
	})
}

func RequirePerm(auth *service.AuthService, perm string, next http.Handler) http.Handler {
	return RequireAuth(auth, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(w, "sem claims", http.StatusUnauthorized)
			return
		}
		if err := auth.Authorize(r.Context(), claims, perm); err != nil {
			http.Error(w, "não autorizado", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

func extractBearer(h string) string {
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
