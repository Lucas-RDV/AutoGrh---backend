package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthConfig struct {
	Issuer          string
	AccessTTL       time.Duration
	ClockSkew       time.Duration
	LoginSuccessEvt int64
	LoginFailEvt    int64

	// Novos (opcionais): se vazios, helpers usam defaults
	Timezone    string // ex.: "America/Campo_Grande"
	CutoffHours []int  // ex.: []int{12, 19}
}

type Claims struct {
	UserID    int64
	Nome      string
	Perfil    string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type UserMinimal struct {
	ID     int64  `json:"id"`
	Nome   string `json:"nome"`
	Login  string `json:"login"`
	Perfil string `json:"perfil"`
}

type PermissionMap map[string]map[string]struct{}

func set(perms ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		m[p] = struct{}{}
	}
	return m
}

type TokenManager interface {
	SignAccess(c Claims) (string, error)
	ParseAccess(token string) (Claims, error)
}

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*UserRecord, error)
	UpdateLastLogin(ctx context.Context, userID int64, when time.Time) error // se não existir, implemente no-op
}

type UserRecord struct {
	ID        int64
	Nome      string
	Login     string
	Perfil    string
	Ativo     bool
	SenhaHash string
}

type LogRepository interface {
	Create(ctx context.Context, entry LogEntry) (int64, error)
}

type LogEntry struct {
	EventoID  int64
	UsuarioID *int64
	Quando    time.Time
	Detalhe   string
}

type AuthService struct {
	userRepo UserRepository
	logRepo  LogRepository
	cfg      AuthConfig
	perms    PermissionMap
	tokens   TokenManager
	clock    func() time.Time
}

func NewAuthService(userRepo UserRepository, logRepo LogRepository, tokens TokenManager, cfg AuthConfig, perms PermissionMap) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logRepo:  logRepo,
		cfg:      cfg,
		perms:    perms,
		tokens:   tokens,
		clock:    time.Now,
	}
}

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrInactiveUser       = errors.New("usuário inativo")
	ErrUnauthorized       = errors.New("não autorizado")
)

// Bcrypt hash dummy pré-computado para timing-safe em usuário inexistente.
// Gerado em tempo de inicialização a partir do texto "invalid".
var dummyBcryptHash = func() []byte {
	h, _ := bcrypt.GenerateFromPassword([]byte("invalid"), bcrypt.DefaultCost)
	return h
}()

func (s *AuthService) Login(ctx context.Context, login, senha string) (accessToken string, expiresAt time.Time, usuario UserMinimal, err error) {
	now := s.clock()
	l := strings.TrimSpace(login)

	// Buscar usuário
	rec, uErr := s.userRepo.GetByLogin(ctx, l)
	if uErr != nil || rec == nil {
		// verificação fake para igualar tempo
		_ = bcrypt.CompareHashAndPassword(dummyBcryptHash, []byte(senha))
		// logar falha (sem userID)
		_ = s.safeAudit(ctx, s.cfg.LoginFailEvt, nil, now, fmt.Sprintf("login=%s", l))
		return "", time.Time{}, UserMinimal{}, ErrInvalidCredentials
	}

	if !rec.Ativo {
		_ = s.safeAudit(ctx, s.cfg.LoginFailEvt, &rec.ID, now, "usuário inativo")
		return "", time.Time{}, UserMinimal{}, ErrInactiveUser
	}

	//Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(rec.SenhaHash), []byte(senha)); err != nil {
		_ = s.safeAudit(ctx, s.cfg.LoginFailEvt, &rec.ID, now, "senha incorreta")
		return "", time.Time{}, UserMinimal{}, ErrInvalidCredentials
	}

	expires := minTime(
		now.Add(s.cfg.AccessTTL),
		nextCutoffLocal(now, s.cfg.Timezone, s.cfg.CutoffHours),
	)

	claims := Claims{
		UserID:    rec.ID,
		Nome:      rec.Nome,
		Perfil:    rec.Perfil,
		Issuer:    s.cfg.Issuer,
		IssuedAt:  now,
		ExpiresAt: expires,
	}
	tok, err := s.tokens.SignAccess(claims)
	if err != nil {
		return "", time.Time{}, UserMinimal{}, fmt.Errorf("emitir token: %w", err)
	}

	_ = s.userRepo.UpdateLastLogin(ctx, rec.ID, now)

	_ = s.safeAudit(ctx, s.cfg.LoginSuccessEvt, &rec.ID, now, "login ok")

	return tok, claims.ExpiresAt, UserMinimal{
		ID:     rec.ID,
		Nome:   rec.Nome,
		Login:  rec.Login,
		Perfil: rec.Perfil,
	}, nil
}

func (s *AuthService) safeAudit(ctx context.Context, eventoID int64, uid *int64, when time.Time, detalhe string) error {
	if s.logRepo == nil || eventoID == 0 {
		return nil
	}
	_, err := s.logRepo.Create(ctx, LogEntry{
		EventoID:  eventoID,
		UsuarioID: uid,
		Quando:    when,
		Detalhe:   detalhe,
	})
	return err
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (Claims, error) {
	claims, err := s.tokens.ParseAccess(token)
	if err != nil {
		return Claims{}, err
	}
	// issuer
	if claims.Issuer != s.cfg.Issuer {
		return Claims{}, fmt.Errorf("issuer inválido: %s", claims.Issuer)
	}
	// expiração
	now := s.clock()
	if now.After(claims.ExpiresAt.Add(s.cfg.ClockSkew)) {
		return Claims{}, errors.New("token expirado")
	}
	return claims, nil
}

func (s *AuthService) Authorize(_ context.Context, c Claims, requiredPerm string) error {
	role := strings.TrimSpace(strings.ToLower(c.Perfil))
	pm, ok := s.perms[role]
	if !ok {
		return ErrUnauthorized
	}

	if _, all := pm["*"]; all {
		return nil
	}
	if requiredPerm == "" {
		return nil
	}
	if _, ok := pm[requiredPerm]; !ok {
		return ErrUnauthorized
	}
	return nil
}

func (s *AuthService) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("senha vazia")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash senha: %w", err)
	}
	return string(hash), nil
}

func (s *AuthService) VerifyPassword(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func parseSubjectToID(sub string) (int64, error) {
	return strconv.ParseInt(sub, 10, 64)
}

// minTime retorna o menor dos dois instantes.
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// nextCutoffLocal calcula o próximo "corte" (12:00 ou 19:00) no fuso definido.
// Se tz/hours vierem vazios, usa defaults: tz="America/Campo_Grande", horas={12,19}.
func nextCutoffLocal(now time.Time, tz string, hours []int) time.Time {
	if tz == "" {
		tz = "America/Campo_Grande"
	}
	if len(hours) == 0 {
		hours = []int{12, 19}
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = now.Location()
	}
	nl := now.In(loc)
	y, m, d := nl.Date()

	// cortes candidatos hoje
	for _, h := range hours {
		c := time.Date(y, m, d, h, 0, 0, 0, loc)
		if nl.Before(c) {
			return c
		}
	}
	// se todos passaram, primeiro corte do dia seguinte
	first := hours[0]
	return time.Date(y, m, d+1, first, 0, 0, 0, loc)
}
