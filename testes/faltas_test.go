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

/************ Log fake ************/

type faltaFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *faltaFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func faltaHasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
	for _, e := range entries {
		if e.EventoID == evt && e.UsuarioID != nil && *e.UsuarioID == uid {
			if prefix == "" || strings.HasPrefix(e.Detalhe, prefix) {
				return true
			}
		}
	}
	return false
}

/************ Helpers ************/

// Pessoa (CPF 11) + Funcionário válidos p/ FK
func seedPessoaFuncionarioFalta(t *testing.T) int64 {
	t.Helper()
	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000)
	rg := fmt.Sprintf("%09d", now%1000000000)

	p := &entity.Pessoa{Nome: "Teste Falta", CPF: cpf, RG: rg}
	if err := repository.CreatePessoa(p); err != nil {
		t.Fatalf("seed CreatePessoa erro: %v", err)
	}
	if p.ID == 0 {
		t.Fatalf("seed pessoa sem ID")
	}

	f := &entity.Funcionario{
		PessoaID:          p.ID,
		PIS:               "PIS-FALTA",
		CTPF:              "CT-FALTA",
		Nascimento:        time.Now().AddDate(-25, 0, 0),
		Admissao:          time.Now().AddDate(-2, 0, 0),
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

func newAdminAuthFalta(lr *faltaFakeLogRepo) *service.AuthService {
	cfg := service.AuthConfig{
		Issuer:          "autogrh-test",
		AccessTTL:       10 * time.Minute,
		ClockSkew:       2 * time.Minute,
		LoginSuccessEvt: 1001,
		LoginFailEvt:    1002,
		Timezone:        "America/Campo_Grande",
	}
	perms := service.PermissionMap{
		"admin": {"*": {}},
	}
	return service.NewAuthService(nil, lr, jwtm.NewHS256Manager([]byte("secret")), cfg, perms)
}

func newFaltaServiceWithDB(lr *faltaFakeLogRepo) *service.FaltaService {
	auth := newAdminAuthFalta(lr)
	adp := Adapter.NewFaltaRepositoryAdapter(
		repository.CreateFalta,
		repository.UpdateFalta,
		repository.DeleteFalta,
		repository.GetFaltaByID,
		repository.GetFaltasByFuncionarioID,
		repository.ListAllFaltas,
	)
	return service.NewFaltaService(auth, lr, adp)
}

/************ TESTES ************/

func TestFalta_Create_Get_Update_List_Delete_WithLogs(t *testing.T) {
	// DB limpo no início e no fim do teste
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &faltaFakeLogRepo{}
	svc := newFaltaServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 101, Perfil: "admin"}

	funcID := seedPessoaFuncionarioFalta(t)

	// Create
	mesRef := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)
	falta := entity.NewFalta(2, mesRef, funcID)
	if err := svc.CreateFalta(ctx, claims, falta); err != nil {
		t.Fatalf("CreateFalta erro: %v", err)
	}
	if falta.ID == 0 {
		t.Fatalf("esperava falta com ID")
	}
	if !faltaHasLogPrefix(lr.entries, 3, 101, "Falta criada ID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}

	// GetByID
	got, err := svc.GetFaltaByID(ctx, claims, falta.ID)
	if err != nil {
		t.Fatalf("GetFaltaByID erro: %v", err)
	}
	if got == nil || got.ID != falta.ID || got.FuncionarioID != funcID {
		t.Fatalf("GetFaltaByID inválido: %+v", got)
	}

	// GetByFuncionarioID
	listFunc, err := svc.GetFaltasByFuncionarioID(ctx, claims, funcID)
	if err != nil {
		t.Fatalf("GetFaltasByFuncionarioID erro: %v", err)
	}
	if len(listFunc) == 0 {
		t.Fatalf("esperava faltas do funcionário")
	}

	// Update (qtd e mês)
	got.Quantidade = 3
	got.Mes = time.Date(2025, time.February, 1, 0, 0, 0, 0, time.Local)
	if err := svc.UpdateFalta(ctx, claims, got); err != nil {
		t.Fatalf("UpdateFalta erro: %v", err)
	}
	if !faltaHasLogPrefix(lr.entries, 4, 101, "Falta atualizada ID=") {
		t.Errorf("não registrou log de update (EventoID=4)")
	}

	// List all
	all, err := svc.ListAllFaltas(ctx, claims)
	if err != nil {
		t.Fatalf("ListAllFaltas erro: %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("esperava ao menos 1 registro de falta")
	}

	// Delete
	if err := svc.DeleteFalta(ctx, claims, falta.ID); err != nil {
		t.Fatalf("DeleteFalta erro: %v", err)
	}
	if !faltaHasLogPrefix(lr.entries, 5, 101, "Falta deletada ID=") {
		t.Errorf("não registrou log de deleção (EventoID=5)")
	}
	// Confirma remoção
	after, err := repository.GetFaltaByID(falta.ID)
	if err != nil {
		t.Fatalf("GetFaltaByID pós-delete erro: %v", err)
	}
	if after != nil {
		t.Fatalf("falta ainda existe após delete: %+v", after)
	}
}

func TestFalta_UpsertMensal_Cria_Atualiza_Zera(t *testing.T) {
	// DB limpo no início e no fim do teste
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &faltaFakeLogRepo{}
	svc := newFaltaServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 102, Perfil: "admin"}

	funcID := seedPessoaFuncionarioFalta(t)
	mes, ano := 1, 2025 // janeiro/2025

	// cria (insere) 5 faltas
	if err := svc.UpsertMensal(ctx, claims, funcID, mes, ano, 5); err != nil {
		t.Fatalf("UpsertMensal (insert) erro: %v", err)
	}
	total, err := repository.GetTotalFaltasByFuncionarioMesAno(funcID, mes, ano)
	if err != nil || total != 5 {
		t.Fatalf("total esperado 5, veio %d (err=%v)", total, err)
	}

	// atualiza para 8
	if err := svc.UpsertMensal(ctx, claims, funcID, mes, ano, 8); err != nil {
		t.Fatalf("UpsertMensal (update) erro: %v", err)
	}
	total, _ = repository.GetTotalFaltasByFuncionarioMesAno(funcID, mes, ano)
	if total != 8 {
		t.Fatalf("total esperado 8, veio %d", total)
	}

	// zera (mantém linha com 0 via update; não cria se não houver)
	if err := svc.UpsertMensal(ctx, claims, funcID, mes, ano, 0); err != nil {
		t.Fatalf("UpsertMensal (zero) erro: %v", err)
	}
	total, _ = repository.GetTotalFaltasByFuncionarioMesAno(funcID, mes, ano)
	if total != 0 {
		t.Fatalf("total esperado 0 após zerar, veio %d", total)
	}
}

func TestFalta_Create_ComQuantidadeZero_DeveErro(t *testing.T) {
	// DB limpo no início e no fim do teste
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &faltaFakeLogRepo{}
	svc := newFaltaServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 103, Perfil: "admin"}

	funcID := seedPessoaFuncionarioFalta(t)
	f0 := entity.NewFalta(0, time.Date(2025, time.March, 1, 0, 0, 0, 0, time.Local), funcID)
	if err := svc.CreateFalta(ctx, claims, f0); err == nil {
		t.Fatalf("esperava erro ao criar falta com quantidade 0")
	}
}
