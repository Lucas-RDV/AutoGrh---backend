package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Usar "*" para permitir qualquer origem
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func defaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposeHeaders:    nil,
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// NewCORS retorna um middleware que aplica CORS.
func NewCORS(cfg CORSConfig) func(http.Handler) http.Handler {
	// completa defaults
	if len(cfg.AllowedOrigins) == 0 && len(cfg.AllowedMethods) == 0 && len(cfg.AllowedHeaders) == 0 {
		cfg = defaultCORSConfig()
	}
	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = defaultCORSConfig().AllowedMethods
	}
	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = defaultCORSConfig().AllowedHeaders
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if allow := matchOrigin(origin, cfg.AllowedOrigins); allow != "" {
					w.Header().Set("Access-Control-Allow-Origin", allow)
					if cfg.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}
					if len(cfg.ExposeHeaders) > 0 {
						w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
					}
					w.Header().Add("Vary", "Origin")
				}
			}

			// Preflight
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				if cfg.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(cfg.MaxAge/time.Second)))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func matchOrigin(origin string, allowed []string) string {
	if origin == "" {
		return ""
	}
	for _, a := range allowed {
		if a == "*" || a == origin {
			if a == "*" {
				return "*"
			}
			return origin
		}
	}
	return ""
}
