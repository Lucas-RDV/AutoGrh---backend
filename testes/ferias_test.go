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

/************ Log fake compartilhado ************/

type fdFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *fdFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func hasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
	for _, e := range entries {
		if e.EventoID == evt && e.UsuarioID != nil && *e.UsuarioID == uid {
			if prefix == "" || strings.HasPrefix(e.Detalhe, prefix) {
				return true
			}
		}
	}
	return false
}

/************ Helpers de seed ************/

// Pessoa (CPF 11 dígitos) + Funcionário válidos para satisfazer FKs
func seedPessoaFuncionarioFD(t *testing.T) int64 {
	t.Helper()
	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000) // 11 dígitos
	rg := fmt.Sprintf("%09d", now%1000000000)     // 9 dígitos

	p := &entity.Pessoa{
		Nome: "Teste Ferias/Descanso",
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
		PIS:               "PIS-FD",
		CTPF:              "CT-FD",
		Nascimento:        time.Now().AddDate(-25, 0, 0),
		Admissao:          time.Now().AddDate(-2, 0, 0), // 2 anos atrás para viabilizar concessões
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

/************ Helpers de Auth/Services ************/

func newAdminAuthFD(lr *fdFakeLogRepo) *service.AuthService {
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

func newFeriasServiceWithDB(lr *fdFakeLogRepo) *service.FeriasService {
	auth := newAdminAuthFD(lr)
	fRepo := Adapter.NewFeriasRepositoryAdapter(
		repository.CreateFerias,
		repository.GetFeriasByFuncionarioID,
		repository.GetFeriasByID,
		repository.UpdateFerias,
		repository.DeleteFerias,
		repository.ListFerias,
	)
	return service.NewFeriasService(auth, lr, fRepo)
}

func newDescansoServiceWithDB(lr *fdFakeLogRepo) *service.DescansoService {
	auth := newAdminAuthFD(lr)
	dRepo := Adapter.NewDescansoRepositoryAdapter(
		repository.CreateDescanso,
		repository.GetDescansoByID,
		repository.UpdateDescanso,
		repository.DeleteDescanso,
		repository.GetDescansosByFeriasID,
		repository.GetDescansosByFuncionarioID,
		repository.GetDescansosAprovados,
		repository.GetDescansosPendentes,
	)
	return service.NewDescansoService(auth, lr, dRepo)
}

/************ TESTES: FÉRIAS ************/

func TestFerias_Criar_Listar_Atualizar_Excluir(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &fdFakeLogRepo{}
	fsvc := newFeriasServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 70, Perfil: "admin"}

	funcID := seedPessoaFuncionarioFD(t)

	// Criar
	inicio := time.Now().AddDate(0, 0, -5) // data válida
	f, err := fsvc.CriarFerias(ctx, claims, funcID, 30, 3000.00, inicio)
	if err != nil {
		t.Fatalf("CriarFerias erro: %v", err)
	}
	if f == nil || f.ID == 0 {
		t.Fatalf("esperava ferias criadas com ID")
	}
	if f.Terco <= 0 {
		t.Fatalf("terço não atribuído")
	}
	if !hasLogPrefix(lr.entries, 3, 70, "Férias criadas") {
		t.Errorf("não registrou log de criação de férias")
	}

	// GetByID
	got, err := fsvc.GetFeriasByID(ctx, claims, f.ID)
	if err != nil || got == nil || got.ID != f.ID {
		t.Fatalf("GetFeriasByID incorreto: got=%+v err=%v", got, err)
	}

	// List por funcionário
	list, err := fsvc.GetFeriasByFuncionarioID(ctx, claims, funcID)
	if err != nil || len(list) == 0 {
		t.Fatalf("GetFeriasByFuncionarioID falhou: len=%d err=%v", len(list), err)
	}

	// Atualizar (ex.: marcar vencido = true)
	f.Vencido = true
	if err := fsvc.AtualizarFerias(ctx, claims, f); err != nil {
		t.Fatalf("AtualizarFerias erro: %v", err)
	}
	if !hasLogPrefix(lr.entries, 4, 70, "Férias atualizadas id=") {
		t.Errorf("não registrou log de update de férias")
	}

	// Deletar
	if err := fsvc.DeletarFerias(ctx, claims, f.ID); err != nil {
		t.Fatalf("DeletarFerias erro: %v", err)
	}
	// confirmar remoção
	check, err := repository.GetFeriasByID(f.ID)
	if err != nil {
		t.Fatalf("GetFeriasByID pós-delete erro: %v", err)
	}
	if check != nil {
		t.Fatalf("férias ainda existem após delete")
	}
	if !hasLogPrefix(lr.entries, 5, 70, "Férias deletadas id=") {
		t.Errorf("não registrou log de delete de férias")
	}
}

/************ TESTES: DESCANSO + integração com FÉRIAS ************/

func TestDescanso_FluxoCompleto_Create_Aprovar_Pagar_ComFechamentoDeFerias(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &fdFakeLogRepo{}
	fsvc := newFeriasServiceWithDB(lr)
	dsvc := newDescansoServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 80, Perfil: "admin"}

	funcID := seedPessoaFuncionarioFD(t)

	// TZ fixa e datas estáveis
	loc, _ := time.LoadLocation("America/Campo_Grande")
	day := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	}
	const Y = 2025
	const M = time.January

	// 1) Cria FÉRIAS com 30 dias a partir de 01/01/2025
	inicioFerias := day(Y, M, 1)
	f, err := fsvc.CriarFerias(ctx, claims, funcID, 30, 3000.00, inicioFerias)
	if err != nil {
		t.Fatalf("CriarFerias erro: %v", err)
	}

	// 2) Cria **30 descansos pendentes** de 1 dia (não aprova ainda)
	var descansoIDs []int64
	for i := 0; i < 30; i++ {
		start := inicioFerias.AddDate(0, 0, i) // 01→30/jan
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
		end := start

		// Checa saldo "de criação": f.DiasRestantes considera SOMA dos já criados (pendentes), mas f.Dias ainda é 30
		fAtual, err := repository.GetFeriasByID(f.ID)
		if err != nil || fAtual == nil {
			t.Fatalf("GetFeriasByID loop (criação) erro: %v", err)
		}

		d := entity.NewDescanso(start, end, f.ID)
		if err := dsvc.CreateDescanso(ctx, claims, d); err != nil {
			t.Fatalf("CreateDescanso (1 dia) erro no dia %02d/%02d/%d: %v", start.Day(), start.Month(), start.Year(), err)
		}
		// NÃO aprovar aqui — evita “dupla subtração” na próxima criação
		restante := fAtual.DiasRestantes() - d.DuracaoEmDias()
		if restante < 0 {
			t.Fatalf("Saldo ficou negativo após criação (inesperado): restante=%d", restante)
		}
		descansoIDs = append(descansoIDs, d.ID)
	}

	// 3) Aprova todos os descansos criados (agora sim consome no banco)
	for _, id := range descansoIDs {
		if err := dsvc.AprovarDescanso(ctx, claims, id); err != nil {
			t.Fatalf("AprovarDescanso erro para ID=%d: %v", id, err)
		}
	}

	// Confere que zerou os dias
	fZero, err := repository.GetFeriasByID(f.ID)
	if err != nil || fZero == nil {
		t.Fatalf("GetFeriasByID pós-aprovação erro: %v", err)
	}
	if fZero.Dias != 0 {
		t.Fatalf("esperava dias=0 após aprovar 30 descansos, veio %v", fZero.Dias)
	}

	// 4) Marca TERÇO como pago
	if err := fsvc.MarcarTercoComoPago(ctx, claims, f.ID); err != nil {
		t.Fatalf("MarcarTercoComoPago erro: %v", err)
	}

	// 5) Paga o **último** descanso → deve fechar as férias como pagas
	lastDescansoID := descansoIDs[len(descansoIDs)-1]
	if err := dsvc.MarcarComoPago(ctx, claims, lastDescansoID); err != nil {
		t.Fatalf("MarcarComoPago(último) erro: %v", err)
	}
	fFinal, _ := repository.GetFeriasByID(f.ID)
	if fFinal == nil || !fFinal.Pago {
		t.Fatalf("Férias deveriam estar PAGAS após pagar último descanso com terçoPago=true e dias=0")
	}

	// Sanidade das listagens
	porFerias, err := dsvc.ListarPorFerias(ctx, claims, f.ID)
	if err != nil || len(porFerias) != 30 {
		t.Fatalf("ListarPorFerias esperado 30 descansos, got=%d err=%v", len(porFerias), err)
	}
	aprovados, err := dsvc.ListarAprovados(ctx, claims)
	if err != nil || len(aprovados) != 30 {
		t.Fatalf("ListarAprovados esperado 30 aprovados, got=%d err=%v", len(aprovados), err)
	}
}
