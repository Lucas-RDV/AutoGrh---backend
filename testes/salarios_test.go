package testes

import (
	Adapter "AutoGRH/pkg/adapter"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/service/jwt"
)

/*** ----------------- Log fake compartilhado ----------------- ***/

type salarioFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *salarioFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func salarioHasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
	for _, e := range entries {
		if e.EventoID == evt && e.UsuarioID != nil && *e.UsuarioID == uid {
			if prefix == "" || strings.HasPrefix(e.Detalhe, prefix) {
				return true
			}
		}
	}
	return false
}

/*** ----------------- Helpers de seed ----------------- ***/

// Cria uma pessoa e um funcionário reais no DB de teste e retorna o funcionarioID.
// Isso satisfaz as FKs de salario/salario_real.
func seedPessoaFuncionarioForSalary(t *testing.T) int64 {
	t.Helper()

	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000) // 11 dígitos
	rg := fmt.Sprintf("%09d", now%1000000000)     // 9 dígitos (seguro p/ coluna comum)

	p := &entity.Pessoa{
		Nome: "Teste Salarios",
		CPF:  cpf,
		RG:   rg,
	}
	if err := repository.CreatePessoa(p); err != nil {
		t.Fatalf("seed CreatePessoa erro: %v", err)
	}
	if p.ID == 0 {
		t.Fatalf("seed pessoa sem ID")
	}

	f := &entity.Funcionario{
		PessoaID:          p.ID,
		PIS:               "PIS-TESTE",
		CTPF:              "CT-TESTE",
		Nascimento:        time.Now().AddDate(-25, 0, 0),
		Admissao:          time.Now().AddDate(-1, 0, 0),
		Cargo:             "Analista",
		SalarioInicial:    2500.00,
		FeriasDisponiveis: 0,
	}
	if err := repository.CreateFuncionario(f); err != nil {
		t.Fatalf("seed CreateFuncionario erro: %v", err)
	}
	if f.ID == 0 {
		t.Fatalf("seed funcionario sem ID")
	}
	return f.ID
}

/*** ----------------- Helpers de Auth/Services ----------------- ***/

func newAdminAuthForSalary(lr *salarioFakeLogRepo) *service.AuthService {
	cfg := service.AuthConfig{
		Issuer:          "autogrh-test",
		AccessTTL:       10 * time.Minute,
		ClockSkew:       2 * time.Minute,
		LoginSuccessEvt: 1001,
		LoginFailEvt:    1002,
		Timezone:        "America/Campo_Grande",
	}
	// perfil admin com "*"
	perms := service.PermissionMap{
		"admin": {"*": {}},
	}
	return service.NewAuthService(nil, lr, jwtm.NewHS256Manager([]byte("secret")), cfg, perms)
}

func newSalarioServiceWithDB(lr *salarioFakeLogRepo) *service.SalarioService {
	auth := newAdminAuthForSalary(lr)
	// usa Adapter apontando para o repositório real (DB de teste)
	repoAdapter := Adapter.NewSalarioRepositoryAdapter(
		repository.CreateSalario,
		repository.GetSalariosByFuncionarioID,
		repository.UpdateSalario,
		repository.DeleteSalario,
	)
	return service.NewSalarioService(auth, lr, repoAdapter)
}

func newSalarioRealServiceWithDB(lr *salarioFakeLogRepo) *service.SalarioRealService {
	auth := newAdminAuthForSalary(lr)
	repoAdapter := Adapter.NewSalarioRealRepositoryAdapter(
		repository.CreateSalarioReal,
		repository.GetSalariosReaisByFuncionarioID,
		repository.GetSalarioRealAtual,
		repository.UpdateSalarioReal,
		repository.DeleteSalarioReal,
	)
	return service.NewSalarioRealService(auth, lr, repoAdapter)
}

/*** ===================== TESTES: SALÁRIO ===================== ***/

// cria salário sem atual → novo ativo + log EventoID=3
func TestSalario_CriarSemAtual_CriaAtivoELoga(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 10, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	s, err := svc.CriarSalario(ctx, claims, funcID, 3000.00)
	if err != nil {
		t.Fatalf("CriarSalario erro: %v", err)
	}
	if s == nil || s.ID == 0 {
		t.Fatalf("esperava salário criado com ID")
	}

	// atual deve ser o recém-criado
	atual, err := repository.GetSalarioAtual(funcID)
	if err != nil || atual == nil || atual.ID != s.ID {
		t.Fatalf("GetSalarioAtual incorreto: got=%+v err=%v", atual, err)
	}

	if !salarioHasLogPrefix(lr.entries, 3, 10, "Salário criado funcionarioID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

// cria salário quando já existe um atual → fecha o anterior e cria novo
func TestSalario_CriarComAtual_FechaAnteriorECriaNovo(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 11, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	// cria um salário "atual" diretamente no repo (fim=nil)
	prev := &entity.Salario{
		FuncionarioID: funcID,
		Inicio:        time.Now().Add(-30 * 24 * time.Hour),
		Valor:         2500.00,
	}
	if err := repository.CreateSalario(prev); err != nil {
		t.Fatalf("seed CreateSalario erro: %v", err)
	}

	// sanity: deve estar como atual
	a1, err := repository.GetSalarioAtual(funcID)
	if err != nil || a1 == nil || a1.ID != prev.ID || a1.Fim != nil {
		t.Fatalf("pré-condição de atual falhou: got=%+v err=%v", a1, err)
	}

	// cria novo via service (deve encerrar o anterior e criar um novo atual)
	newS, err := svc.CriarSalario(ctx, claims, funcID, 2800.00)
	if err != nil {
		t.Fatalf("CriarSalario erro: %v", err)
	}

	// anterior deve estar com FIM preenchido (no histórico)
	list, err := repository.GetSalariosByFuncionarioID(funcID)
	if err != nil {
		t.Fatalf("GetSalariosByFuncionarioID erro: %v", err)
	}
	var prevClosed bool
	for _, s := range list {
		if s.ID == prev.ID && s.Fim != nil {
			prevClosed = true
		}
	}
	if !prevClosed {
		t.Fatalf("esperava salário anterior encerrado (Fim != nil)")
	}

	// atual agora deve ser o novo
	a2, err := repository.GetSalarioAtual(funcID)
	if err != nil || a2 == nil || a2.ID != newS.ID {
		t.Fatalf("atual após criação não é o novo: got=%+v err=%v", a2, err)
	}

	if !salarioHasLogPrefix(lr.entries, 3, 11, "Salário criado funcionarioID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

func TestSalario_List_Atualizar_Excluir_ComLogs(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 12, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	// cria dois salários sequenciais (service já encerra o anterior)
	s1, _ := svc.CriarSalario(ctx, claims, funcID, 2000.00)
	time.Sleep(10 * time.Millisecond) // separa instantes
	s2, _ := svc.CriarSalario(ctx, claims, funcID, 2200.00)

	// list
	list, err := svc.ListSalarios(ctx, claims, funcID)
	if err != nil || len(list) < 2 {
		t.Fatalf("ListSalarios erro ou insuficiente: len=%d err=%v", len(list), err)
	}

	// atualizar s2
	s2.Valor = 2300.00
	if err := svc.AtualizarSalario(ctx, claims, s2); err != nil {
		t.Fatalf("AtualizarSalario erro: %v", err)
	}
	list2, _ := repository.GetSalariosByFuncionarioID(funcID)
	var foundUpdated bool
	for _, s := range list2 {
		if s.ID == s2.ID && s.Valor == 2300.00 {
			foundUpdated = true
		}
	}
	if !foundUpdated {
		t.Fatalf("esperava salário atualizado com valor=2300.00")
	}
	if !salarioHasLogPrefix(lr.entries, 4, 12, "Salário registrado atualizado id=") {
		t.Errorf("não registrou log de update (EventoID=4)")
	}

	// deletar s1 (precisa permissão; nosso admin tem "*")
	if err := svc.DeletarSalario(ctx, claims, s1.ID); err != nil {
		t.Fatalf("DeletarSalario erro: %v", err)
	}
	list3, _ := repository.GetSalariosByFuncionarioID(funcID)
	for _, s := range list3 {
		if s.ID == s1.ID {
			t.Fatalf("salário deletado ainda presente no histórico")
		}
	}
	if !salarioHasLogPrefix(lr.entries, 5, 12, "Salário registrado deletado id=") {
		t.Errorf("não registrou log de delete (EventoID=5)")
	}
}

/*** ===================== TESTES: SALÁRIO REAL ===================== ***/

// cria salário real sem atual → novo ativo + log
func TestSalarioReal_CriarSemAtual_CriaAtivoELoga(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioRealServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 20, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	r, err := svc.CriarSalarioReal(ctx, claims, funcID, 3100.00)
	if err != nil {
		t.Fatalf("CriarSalarioReal erro: %v", err)
	}
	if r == nil || r.ID == 0 {
		t.Fatalf("esperava salário real criado com ID")
	}

	atual, err := repository.GetSalarioRealAtual(funcID)
	if err != nil || atual == nil || atual.ID != r.ID {
		t.Fatalf("GetSalarioRealAtual incorreto: got=%+v err=%v", atual, err)
	}

	if !salarioHasLogPrefix(lr.entries, 3, 20, "Salário real criado funcionarioID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

// cria salário real quando já existe atual → fecha o anterior e cria novo
func TestSalarioReal_CriarComAtual_FechaAnteriorECriaNovo(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioRealServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 21, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	// seed: salário real atual direto no repo
	prev := &entity.SalarioReal{
		FuncionarioID: funcID,
		Inicio:        time.Now().Add(-15 * 24 * time.Hour),
		Valor:         2700.00,
	}
	if err := repository.CreateSalarioReal(prev); err != nil {
		t.Fatalf("seed CreateSalarioReal erro: %v", err)
	}

	a1, err := repository.GetSalarioRealAtual(funcID)
	if err != nil || a1 == nil || a1.ID != prev.ID || a1.Fim != nil {
		t.Fatalf("pré-condição de atual falhou: got=%+v err=%v", a1, err)
	}

	newR, err := svc.CriarSalarioReal(ctx, claims, funcID, 2900.00)
	if err != nil {
		t.Fatalf("CriarSalarioReal erro: %v", err)
	}

	// anterior encerrado?
	hist, err := repository.GetSalariosReaisByFuncionarioID(funcID)
	if err != nil {
		t.Fatalf("GetSalariosReaisByFuncionarioID erro: %v", err)
	}
	var prevClosed bool
	for _, s := range hist {
		if s.ID == prev.ID && s.Fim != nil {
			prevClosed = true
		}
	}
	if !prevClosed {
		t.Fatalf("esperava salário real anterior encerrado (Fim != nil)")
	}

	// novo é o atual
	a2, err := repository.GetSalarioRealAtual(funcID)
	if err != nil || a2 == nil || a2.ID != newR.ID {
		t.Fatalf("atual após criação não é o novo: got=%+v err=%v", a2, err)
	}

	if !salarioHasLogPrefix(lr.entries, 3, 21, "Salário real criado funcionarioID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

func TestSalarioReal_List_GetAtual_Delete_ComLogs(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &salarioFakeLogRepo{}
	svc := newSalarioRealServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 22, Perfil: "admin"}

	funcID := seedPessoaFuncionarioForSalary(t)

	r1, _ := svc.CriarSalarioReal(ctx, claims, funcID, 1800.00)
	time.Sleep(10 * time.Millisecond)
	r2, _ := svc.CriarSalarioReal(ctx, claims, funcID, 2000.00)

	// list
	list, err := svc.ListSalariosReais(ctx, claims, funcID)
	if err != nil || len(list) < 2 {
		t.Fatalf("ListSalariosReais erro ou insuficiente: len=%d err=%v", len(list), err)
	}

	// get atual
	atual, err := svc.GetSalarioRealAtual(ctx, claims, funcID)
	if err != nil || atual == nil || atual.ID != r2.ID {
		t.Fatalf("GetSalarioRealAtual incorreto: got=%+v err=%v", atual, err)
	}

	// deletar r1
	if err := svc.DeleteSalarioReal(ctx, claims, r1.ID); err != nil {
		t.Fatalf("DeleteSalarioReal erro: %v", err)
	}
	list2, _ := repository.GetSalariosReaisByFuncionarioID(funcID)
	for _, s := range list2 {
		if s.ID == r1.ID {
			t.Fatalf("salário real deletado ainda presente no histórico")
		}
	}
	if !salarioHasLogPrefix(lr.entries, 5, 22, "Salário real deletado id=") {
		t.Errorf("não registrou log de delete (EventoID=5)")
	}
}
