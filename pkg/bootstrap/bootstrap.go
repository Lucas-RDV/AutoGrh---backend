package Bootstrap

import (
	"context"
	"os"
	"strconv"
	"time"

	"AutoGRH/pkg/Adapter"
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	jwtm "AutoGRH/pkg/service/jwt"
)

type AppConfig struct {
	JWTSecret string
	Auth      service.AuthConfig
	Perms     service.PermissionMap
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvIntDefault(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvInt64Default(k string, def int64) int64 {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

func Load() AppConfig {
	issuer := getenvDefault("JWT_ISSUER", "AutoGRH")
	ttlH := getenvIntDefault("JWT_EXPIRES_IN_HOURS", 6)
	skewS := getenvIntDefault("JWT_CLOCK_SKEW_SECONDS", 60)
	successID := getenvInt64Default("EVENT_LOGIN_SUCCESS_ID", 1)
	failID := getenvInt64Default("EVENT_LOGIN_FAIL_ID", 2)

	cfg := service.AuthConfig{
		Issuer:          issuer,
		AccessTTL:       time.Duration(ttlH) * time.Hour,
		ClockSkew:       time.Duration(skewS) * time.Second,
		LoginSuccessEvt: successID,
		LoginFailEvt:    failID,
		Timezone:        "America/Campo_Grande",
		CutoffHours:     []int{12, 19},
	}

	perms := service.PermissionMap{
		"admin":   {"*": {}},
		"usuario": {"self:read": {}, "self:update": {}, "ferias:request": {}},
	}

	return AppConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
		Auth:      cfg,
		Perms:     perms,
	}
}

func ConnectDB() error {
	// adapte para a sua função concreta, se necessário
	if repository.DB == nil {
		repository.ConnectDB()
	}
	return repository.DB.PingContext(context.Background())
}

func BuildAuth(app AppConfig) *service.AuthService {
	tm := jwtm.NewHS256Manager([]byte(app.JWTSecret))
	userRepo := Adapter.NewUserRepositoryAdapter(repository.GetUsuarioByUsername, nil)
	createLog := func(ctx context.Context, l *Entity.Log) (int64, error) { return 0, repository.CreateLog(l) }
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)
	return service.NewAuthService(userRepo, logRepo, tm, app.Auth, app.Perms)
}
