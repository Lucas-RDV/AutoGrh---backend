package testes

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/service/jwt"
)

/*** ---------------- Mocks (exclusivos p/ Funcionario) ---------------- ***/

type funcionarioFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *funcionarioFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func funcionarioHasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
	for _, e := range entries {
		if e.EventoID == evt && e.UsuarioID != nil && *e.UsuarioID == uid {
			if prefix == "" || strings.HasPrefix(e.Detalhe, prefix) {
				return true
			}
		}
	}
	return false
}

// funcionarioFakeRepo implementa service.FuncionarioRepository
type funcionarioFakeRepo struct {
	seq int64
	m   map[int64]*entity.Funcionario
}

func newFuncionarioFakeRepo() *funcionarioFakeRepo {
	return &funcionarioFakeRepo{m: make(map[int64]*entity.Funcionario)}
}

func (r *funcionarioFakeRepo) Create(ctx context.Context, f *entity.Funcionario) error {
	r.seq++
	cp := *f
	cp.ID = r.seq
	// DB real marca ativo=true na criação; o service não seta. Simulamos aqui.
	cp.Ativo = true
	r.m[cp.ID] = &cp
	f.ID = cp.ID
	return nil
}

func (r *funcionarioFakeRepo) GetByID(ctx context.Context, id int64) (*entity.Funcionario, error) {
	if f, ok := r.m[id]; ok {
		cp := *f
		return &cp, nil
	}
	return nil, nil
}

func (r *funcionarioFakeRepo) Update(ctx context.Context, f *entity.Funcionario) error {
	if f.ID <= 0 {
		return errors.New("id inválido no repo")
	}
	if _, ok := r.m[f.ID]; !ok {
		return errors.New("não encontrado")
	}
	cp := *f
	// preserva Ativo atual se não vier setado
	if old, ok := r.m[f.ID]; ok {
		cp.Ativo = old.Ativo
	}
	r.m[f.ID] = &cp
	return nil
}

func (r *funcionarioFakeRepo) Delete(ctx context.Context, id int64) error {
	if f, ok := r.m[id]; ok {
		cp := *f
		cp.Ativo = false // soft delete
		r.m[id] = &cp
	}
	return nil
}

func (r *funcionarioFakeRepo) ListAtivos(ctx context.Context) ([]*entity.Funcionario, error) {
	var out []*entity.Funcionario
	for _, v := range r.m {
		if v.Ativo {
			cp := *v
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *funcionarioFakeRepo) ListInativos(ctx context.Context) ([]*entity.Funcionario, error) {
	var out []*entity.Funcionario
	for _, v := range r.m {
		if !v.Ativo {
			cp := *v
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *funcionarioFakeRepo) ListTodos(ctx context.Context) ([]*entity.Funcionario, error) {
	var out []*entity.Funcionario
	for _, v := range r.m {
		cp := *v
		out = append(out, &cp)
	}
	return out, nil
}

/*** -------------- SUT helpers -------------- ***/

func newAuthForFuncionarioSuite(lr *funcionarioFakeLogRepo) *service.AuthService {
	cfg := service.AuthConfig{
		Issuer:          "autogrh-test",
		AccessTTL:       15 * time.Minute,
		ClockSkew:       2 * time.Minute,
		LoginSuccessEvt: 1001,
		LoginFailEvt:    1002,
		Timezone:        "America/Campo_Grande",
	}
	return service.NewAuthService(nil, lr, jwtm.NewHS256Manager([]byte("secret")), cfg, nil)
}

func newFuncionarioServiceSUT(repo service.FuncionarioRepository, lr *funcionarioFakeLogRepo) *service.FuncionarioService {
	auth := newAuthForFuncionarioSuite(lr)
	return service.NewFuncionarioService(auth, lr, repo)
}

/*** ---------------- Tests: Create ---------------- ***/

func TestFuncionario_CreateFuncionario_Sucesso(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 99}

	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour).Truncate(time.Second)
	adm := time.Now().Add(-30 * 24 * time.Hour).Truncate(time.Second)

	f := &entity.Funcionario{
		PessoaID:          1,
		PIS:               "123",
		CTPF:              "CT-001",
		Nascimento:        nasc,
		Admissao:          adm,
		Cargo:             "  Analista  ", // valida trim
		SalarioInicial:    3500.0,
		FeriasDisponiveis: 30,
	}
	if err := svc.CreateFuncionario(ctx, claims, f); err != nil {
		t.Fatalf("CreateFuncionario erro: %v", err)
	}
	if f.ID == 0 {
		t.Fatalf("esperava ID preenchido")
	}
	got, _ := repo.GetByID(ctx, f.ID)
	if got.Cargo != "Analista" {
		t.Fatalf("esperava Cargo sem espaços após trim, veio %q", got.Cargo)
	}
	if !got.Ativo {
		t.Fatalf("esperava Ativo=true após criação")
	}
	if !funcionarioHasLogPrefix(lr.entries, 3, 99, "Criou funcionário ID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

func TestFuncionario_CreateFuncionario_Validacoes(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}
	nasc := time.Now().Add(-20 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-10 * 24 * time.Hour)

	cases := []struct {
		name string
		f    *entity.Funcionario
		want string
	}{
		{"PessoaID inválido", &entity.Funcionario{PessoaID: 0, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "PessoaID inválido ou não informado"},
		{"Cargo vazio", &entity.Funcionario{PessoaID: 1, Cargo: "  ", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "cargo não pode ser vazio"},
		{"Salário negativo", &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: -1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "salário inicial não pode ser negativo"},
		{"Férias negativas", &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: -1, Admissao: adm, Nascimento: nasc}, "dias de férias disponíveis não podem ser negativos"},
		{"Admissão zero", &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: time.Time{}, Nascimento: nasc}, "data de admissão inválida"},
		{"Nascimento zero", &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: time.Time{}}, "data de nascimento inválida"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.CreateFuncionario(ctx, claims, tc.f)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("esperava erro %q, veio: %v", tc.want, err)
			}
		})
	}
}

/*** ---------------- Tests: GetByID ---------------- ***/

func TestFuncionario_GetFuncionarioByID(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	if _, err := svc.GetFuncionarioByID(ctx, claims, 0); err == nil {
		t.Fatalf("esperava erro por ID inválido")
	}

	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-30 * 24 * time.Hour)
	f := &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}
	if err := svc.CreateFuncionario(ctx, claims, f); err != nil {
		t.Fatalf("create falhou: %v", err)
	}

	got, err := svc.GetFuncionarioByID(ctx, claims, f.ID)
	if err != nil || got == nil || got.ID != f.ID {
		t.Fatalf("GetByID falhou: got=%+v err=%v", got, err)
	}
}

/*** ---------------- Tests: Update ---------------- ***/

func TestFuncionario_UpdateFuncionario_Sucesso_E_Log(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 7}

	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-30 * 24 * time.Hour)
	f := &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1000, FeriasDisponiveis: 10, Admissao: adm, Nascimento: nasc}
	if err := svc.CreateFuncionario(ctx, claims, f); err != nil {
		t.Fatalf("create falhou: %v", err)
	}

	up := &entity.Funcionario{
		ID:                f.ID,
		PessoaID:          1,
		PIS:               "999",
		CTPF:              "CT-999",
		Nascimento:        nasc,
		Admissao:          adm,
		Cargo:             "  Pleno ",
		SalarioInicial:    4200.00,
		FeriasDisponiveis: 20,
	}
	if err := svc.UpdateFuncionario(ctx, claims, up); err != nil {
		t.Fatalf("update falhou: %v", err)
	}

	got, _ := repo.GetByID(ctx, f.ID)
	if got.Cargo != "Pleno" || got.SalarioInicial != 4200.00 || got.PIS != "999" || got.CTPF != "CT-999" || got.FeriasDisponiveis != 20 {
		t.Fatalf("update não aplicou corretamente: %+v", got)
	}
	if !funcionarioHasLogPrefix(lr.entries, 4, 7, "Atualizou funcionário ID=") {
		t.Errorf("não registrou log de update (EventoID=4)")
	}
}

func TestFuncionario_UpdateFuncionario_Validacoes(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}
	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-30 * 24 * time.Hour)

	cases := []struct {
		name string
		f    *entity.Funcionario
		want string
	}{
		{"ID inválido", &entity.Funcionario{ID: 0, PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "ID do funcionário inválido"},
		{"PessoaID inválido", &entity.Funcionario{ID: 1, PessoaID: 0, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "PessoaID inválido"},
		{"Cargo vazio", &entity.Funcionario{ID: 1, PessoaID: 1, Cargo: "  ", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "cargo não pode ser vazio"},
		{"Salário negativo", &entity.Funcionario{ID: 1, PessoaID: 1, Cargo: "A", SalarioInicial: -1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}, "salário inicial não pode ser negativo"},
		{"Férias negativas", &entity.Funcionario{ID: 1, PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: -1, Admissao: adm, Nascimento: nasc}, "dias de férias disponíveis não podem ser negativos"},
		{"Admissão zero", &entity.Funcionario{ID: 1, PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: time.Time{}, Nascimento: nasc}, "data de admissão inválida"},
		{"Nascimento zero", &entity.Funcionario{ID: 1, PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: time.Time{}}, "data de nascimento inválido"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.UpdateFuncionario(ctx, claims, tc.f)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("esperava erro %q, veio: %v", tc.want, err)
			}
		})
	}
}

/*** ---------------- Tests: Delete ---------------- ***/

func TestFuncionario_DeleteFuncionario_Sucesso_E_Log(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 55}

	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-30 * 24 * time.Hour)
	f := &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}
	if err := svc.CreateFuncionario(ctx, claims, f); err != nil {
		t.Fatalf("create falhou: %v", err)
	}

	if err := svc.DeleteFuncionario(ctx, claims, f.ID); err != nil {
		t.Fatalf("delete falhou: %v", err)
	}
	got, _ := repo.GetByID(ctx, f.ID)
	if got == nil || got.Ativo {
		t.Fatalf("esperava funcionário inativo após delete, veio: %+v", got)
	}
	if !funcionarioHasLogPrefix(lr.entries, 5, 55, "Deletou funcionário ID=") {
		t.Errorf("não registrou log de delete (EventoID=5)")
	}
}

func TestFuncionario_DeleteFuncionario_ValidacaoID(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	if err := svc.DeleteFuncionario(ctx, claims, 0); err == nil {
		t.Fatalf("esperava erro por ID inválido")
	}
}

/*** ---------------- Tests: List ---------------- ***/

func TestFuncionario_ListasAtivosInativosETodos(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &funcionarioFakeLogRepo{}
	repo := newFuncionarioFakeRepo()
	svc := newFuncionarioServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	nasc := time.Now().Add(-25 * 365 * 24 * time.Hour)
	adm := time.Now().Add(-30 * 24 * time.Hour)

	// cria 3 (todos ativos)
	a := &entity.Funcionario{PessoaID: 1, Cargo: "A", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}
	b := &entity.Funcionario{PessoaID: 2, Cargo: "B", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}
	c := &entity.Funcionario{PessoaID: 3, Cargo: "C", SalarioInicial: 1, FeriasDisponiveis: 0, Admissao: adm, Nascimento: nasc}
	_ = svc.CreateFuncionario(ctx, claims, a)
	_ = svc.CreateFuncionario(ctx, claims, b)
	_ = svc.CreateFuncionario(ctx, claims, c)

	// inativar um
	_ = svc.DeleteFuncionario(ctx, claims, b.ID)

	ativos, err := svc.ListFuncionariosAtivos(ctx, claims)
	if err != nil {
		t.Fatalf("ListAtivos erro: %v", err)
	}
	if len(ativos) != 2 {
		t.Fatalf("esperava 2 ativos, veio %d", len(ativos))
	}

	inativos, err := svc.ListFuncionariosInativos(ctx, claims)
	if err != nil {
		t.Fatalf("ListInativos erro: %v", err)
	}
	if len(inativos) != 1 || inativos[0].ID != b.ID {
		t.Fatalf("esperava 1 inativo (B), veio: %+v", inativos)
	}

	todos, err := svc.ListTodosFuncionarios(ctx, claims)
	if err != nil {
		t.Fatalf("ListTodos erro: %v", err)
	}
	if len(todos) != 3 {
		t.Fatalf("esperava 3 no total, veio %d", len(todos))
	}
}
