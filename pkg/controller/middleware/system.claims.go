package middleware

import (
	"AutoGRH/pkg/service"
	"time"
)

// SystemClaims retorna um conjunto de claims de "usuário técnico"
// usado por workers ou processos internos do sistema.
func SystemClaims() service.Claims {
	return service.Claims{
		UserID:    0,
		Nome:      "system",
		Perfil:    "admin",
		Issuer:    "autogrh-worker",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}
