package testes

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"AutoGRH/pkg/service"
	"AutoGRH/pkg/service/jwt"
)

// ---------------------- Mocks/Helpers ----------------------

type fakeUserRepo struct {
	users           map[string]*service.UserRecord // login -> record
	lastLoginCalled bool
	lastLoginID     int64
	lastLoginWhen   time.Time
}

func (r *fakeUserRepo) GetByLogin(ctx context.Context, login string) (*service.UserRecord, error) {
	if u, ok := r.users[login]; ok {
		return u, nil
	}
	return nil, nil
}

func (r *fakeUserRepo) UpdateLastLogin(ctx context.Context, userID int64, when time.Time) error {
	r.lastLoginCalled = true
	r.lastLoginID = userID
	r.lastLoginWhen = when
	return nil
}

type fakeLogRepo struct {
	entries []service.LogEntry
}

func (l *fakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func hasLog(entries []service.LogEntry, evt int64, uid *int64, contains string) bool {
	for _, e := range entries {
		if e.EventoID == evt {
			if (uid == nil && e.UsuarioID == nil) || (uid != nil && e.UsuarioID != nil && *uid == *e.UsuarioID) {
				if contains == "" || (contains != "" && (e.Detalhe == contains || (len(e.Detalhe) >= len(contains) && e.Detalhe[:len(contains)] == contains))) {
					return true
				}
			}
		}
	}
	return false
}

func perms(adminAll bool) service.PermissionMap {
	pm := service.PermissionMap{
		"usuario": {"coisa:list": {}},
	}
	if adminAll {
		pm["admin"] = map[string]struct{}{"*": {}}
	} else {
		pm["admin"] = map[string]struct{}{"coisa:list": {}, "coisa:create": {}}
	}
	return pm
}

// ---------------------- Construtor do SUT ----------------------

func newAuthServiceForTest(userRepo service.UserRepository, logRepo service.LogRepository) *service.AuthService {
	cfg := service.AuthConfig{
		Issuer:          "autogrh-test",
		AccessTTL:       30 * time.Minute,
		ClockSkew:       2 * time.Minute,
		LoginSuccessEvt: 1001,
		LoginFailEvt:    1002,
		Timezone:        "America/Campo_Grande",
		CutoffHours:     []int{12, 19},
	}
	tokens := jwtm.NewHS256Manager([]byte("secret-test"))
	return service.NewAuthService(userRepo, logRepo, tokens, cfg, perms(true))
}

// ---------------------- Testes: Login ----------------------

func TestAuth_Login_Sucesso_GeraToken_AtualizaLastLogin_Audita(t *testing.T) {
	defer func() { _ = truncateAll() }()

	// usuário ativo com senha correta (senha hash bcrypt)
	hash, err := service.NewAuthService(nil, nil, nil, service.AuthConfig{}, nil).HashPassword("senha123")
	if err != nil {
		t.Fatalf("erro ao hashear senha de teste: %v", err)
	}
	u := &service.UserRecord{
		ID:        42,
		Nome:      "Alice",
		Login:     "alice",
		Perfil:    "admin",
		Ativo:     true,
		SenhaHash: hash,
	}
	ur := &fakeUserRepo{users: map[string]*service.UserRecord{"alice": u}}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	now := time.Now()
	tok, exp, usr, err := auth.Login(context.Background(), "alice", "senha123")
	if err != nil {
		t.Fatalf("Login erro: %v", err)
	}
	if tok == "" {
		t.Fatalf("esperava token não vazio")
	}
	if usr.ID != 42 || usr.Login != "alice" || usr.Perfil != "admin" {
		t.Errorf("UserMinimal inválido: %+v", usr)
	}
	// exp deve estar entre now e now+TTL (minTime garante <= TTL)
	if !exp.After(now) || exp.After(now.Add(30*time.Minute+5*time.Second)) {
		t.Errorf("ExpiresAt fora do intervalo esperado: %v", exp)
	}
	// UpdateLastLogin chamado
	if !ur.lastLoginCalled || ur.lastLoginID != 42 {
		t.Errorf("UpdateLastLogin não foi chamado corretamente: called=%v id=%d", ur.lastLoginCalled, ur.lastLoginID)
	}
	// Audit de sucesso
	if !hasLog(lr.entries, 1001, &u.ID, "login ok") {
		t.Errorf("não registrou log de sucesso (evento=1001)")
	}
}

func TestAuth_Login_UsuarioInativo(t *testing.T) {
	defer func() { _ = truncateAll() }()

	hash, _ := service.NewAuthService(nil, nil, nil, service.AuthConfig{}, nil).HashPassword("x")
	u := &service.UserRecord{ID: 7, Nome: "Bob", Login: "bob", Perfil: "usuario", Ativo: false, SenhaHash: hash}
	ur := &fakeUserRepo{users: map[string]*service.UserRecord{"bob": u}}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	_, _, _, err := auth.Login(context.Background(), "bob", "x")
	if !errors.Is(err, service.ErrInactiveUser) {
		t.Fatalf("esperava ErrInactiveUser, veio: %v", err)
	}
	// Audit de falha (inativo)
	if !hasLog(lr.entries, 1002, &u.ID, "usuário inativo") {
		t.Errorf("não registrou log de falha por usuário inativo (evento=1002)")
	}
}

func TestAuth_Login_Inexistente_ou_SenhaIncorreta(t *testing.T) {
	defer func() { _ = truncateAll() }()

	// Cenário 1: usuário inexistente
	ur := &fakeUserRepo{users: map[string]*service.UserRecord{}}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	_, _, _, err := auth.Login(context.Background(), "ghost", "qualquer")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("inexistente: esperava ErrInvalidCredentials, veio: %v", err)
	}
	// Audit de falha sem userID
	if !hasLog(lr.entries, 1002, nil, "login=ghost") {
		t.Errorf("não registrou log de falha de login inexistente (evento=1002)")
	}

	// Cenário 2: senha incorreta
	hash, _ := service.NewAuthService(nil, nil, nil, service.AuthConfig{}, nil).HashPassword("correta")
	u := &service.UserRecord{ID: 9, Nome: "Cara", Login: "cara", Perfil: "usuario", Ativo: true, SenhaHash: hash}
	ur.users["cara"] = u

	_, _, _, err = auth.Login(context.Background(), "cara", "errada")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("senha incorreta: esperava ErrInvalidCredentials, veio: %v", err)
	}
	if !hasLog(lr.entries, 1002, &u.ID, "senha incorreta") {
		t.Errorf("não registrou log de senha incorreta (evento=1002)")
	}
}

// ---------------------- Testes: ValidateToken ----------------------

func TestAuth_ValidateToken_Sucesso_Issuer_e_ExpOk(t *testing.T) {
	defer func() { _ = truncateAll() }()

	ur := &fakeUserRepo{}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	// gerar token real via SignAccess
	now := time.Now()
	claims := service.Claims{
		UserID:    1,
		Nome:      "Teste",
		Perfil:    "admin",
		Issuer:    "autogrh-test",
		IssuedAt:  now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
	tok, err := authTokens(auth).SignAccess(claims)
	if err != nil {
		t.Fatalf("erro ao assinar token de teste: %v", err)
	}

	out, err := auth.ValidateToken(context.Background(), tok)
	if err != nil {
		t.Fatalf("ValidateToken erro: %v", err)
	}
	if out.UserID != 1 || out.Issuer != "autogrh-test" {
		t.Errorf("claims inválidas: %+v", out)
	}
}

func TestAuth_ValidateToken_IssuerInvalido(t *testing.T) {
	defer func() { _ = truncateAll() }()

	ur := &fakeUserRepo{}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	now := time.Now()
	claims := service.Claims{
		UserID:    1,
		Nome:      "Teste",
		Perfil:    "admin",
		Issuer:    "issuer-errado", // não bate com cfg.Issuer
		IssuedAt:  now,
		ExpiresAt: now.Add(10 * time.Minute),
	}
	tok, _ := authTokens(auth).SignAccess(claims)

	_, err := auth.ValidateToken(context.Background(), tok)
	if err == nil || !strings.HasPrefix(err.Error(), "issuer inválido") {
		t.Fatalf("esperava erro de issuer inválido, veio: %v", err)
	}
}

func TestAuth_ValidateToken_Expirado_ComClockSkew(t *testing.T) {
	defer func() { _ = truncateAll() }()

	// Usamos o auth padrão (ClockSkew=2m), mas expiramos há mais de 2m para garantir falha
	ur := &fakeUserRepo{}
	lr := &fakeLogRepo{}
	auth := newAuthServiceForTest(ur, lr)

	now := time.Now()
	claims := service.Claims{
		UserID:    1,
		Nome:      "Teste",
		Perfil:    "admin",
		Issuer:    "autogrh-test",
		IssuedAt:  now.Add(-15 * time.Minute),
		ExpiresAt: now.Add(-5 * time.Minute),
	}
	tok, _ := authTokens(auth).SignAccess(claims)

	_, err := auth.ValidateToken(context.Background(), tok)
	if err == nil || err.Error() != "token expirado" {
		t.Fatalf("esperava 'token expirado', veio: %v", err)
	}
}

// ---------------------- Testes: Authorize ----------------------

func TestAuth_Authorize(t *testing.T) {
	defer func() { _ = truncateAll() }()

	ur := &fakeUserRepo{}
	lr := &fakeLogRepo{}
	// admin com "*" (libera tudo)
	authAdminAll := service.NewAuthService(ur, lr, jwtm.NewHS256Manager([]byte("s")), service.AuthConfig{Issuer: "autogrh-test"}, perms(true))
	// admin sem "*" (somente algumas)
	authAdminSome := service.NewAuthService(ur, lr, jwtm.NewHS256Manager([]byte("s")), service.AuthConfig{Issuer: "autogrh-test"}, perms(false))

	admin := service.Claims{Perfil: "admin"}
	user := service.Claims{Perfil: "usuario"}

	// adminAll permite tudo
	if err := authAdminAll.Authorize(context.Background(), admin, "qualquer:coisa"); err != nil {
		t.Fatalf("adminAll deveria permitir tudo, erro: %v", err)
	}
	// adminSome permite apenas list/create
	if err := authAdminSome.Authorize(context.Background(), admin, "coisa:create"); err != nil {
		t.Fatalf("adminSome deveria permitir coisa:create, erro: %v", err)
	}
	if err := authAdminSome.Authorize(context.Background(), admin, "coisa:delete"); err == nil {
		t.Fatalf("adminSome NÃO deveria permitir coisa:delete")
	}
	// usuario só tem coisa:list
	if err := authAdminSome.Authorize(context.Background(), user, "coisa:list"); err != nil {
		t.Fatalf("usuario deveria poder coisa:list, erro: %v", err)
	}
	if err := authAdminSome.Authorize(context.Background(), user, "coisa:create"); err == nil {
		t.Fatalf("usuario NÃO deveria poder coisa:create")
	}
}

// ---------------------- Testes: Hash/Verify ----------------------

func TestAuth_HashPassword_e_VerifyPassword(t *testing.T) {
	defer func() { _ = truncateAll() }()

	auth := newAuthServiceForTest(&fakeUserRepo{}, &fakeLogRepo{})

	hash, err := auth.HashPassword("segredo")
	if err != nil || hash == "" {
		t.Fatalf("HashPassword falhou: %v (%s)", err, hash)
	}
	if !auth.VerifyPassword("segredo", hash) {
		t.Fatalf("VerifyPassword deveria aceitar a senha correta")
	}
	if auth.VerifyPassword("errada", hash) {
		t.Fatalf("VerifyPassword não deveria aceitar senha incorreta")
	}
}

// ---------------------- util: acesso ao TokenManager do auth ----------------------

// como AuthService não expõe o TokenManager, criamos esta ajudante
type tokenMgrAccessor interface {
	SignAccess(c service.Claims) (string, error)
	ParseAccess(token string) (service.Claims, error)
}

func authTokens(a *service.AuthService) tokenMgrAccessor {
	// neste teste, sempre criamos com jwtm.HS256Manager; portanto podemos só recriar igual
	// para assinar tokens com mesmo segredo da suíte "newAuthServiceForTest"
	// (para evitar refletir em campo privado)
	return jwtm.NewHS256Manager([]byte("secret-test"))
}
