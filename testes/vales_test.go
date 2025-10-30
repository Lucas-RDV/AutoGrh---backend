package testes

import (
	Adapter "AutoGRH/pkg/adapter"
	"context"
	"fmt"
	"testing"
	"time"

	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/service/jwt"
)

/************ Fake Log ************/
type valeFakeLogRepo struct{ entries []service.LogEntry }

func (l *valeFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

/************ Helpers ************/
func seedPessoaFuncionarioVale(t *testing.T) int64 {
	t.Helper()
	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000)
	rg := fmt.Sprintf("%09d", now%1000000000)

	p := &entity.Pessoa{Nome: "Teste Vale", CPF: cpf, RG: rg}
	if err := repository.CreatePessoa(p); err != nil {
		t.Fatalf("seed CreatePessoa erro: %v", err)
	}
	if p.ID == 0 {
		t.Fatalf("seed pessoa sem ID")
	}
	f := &entity.Funcionario{
		PessoaID:          p.ID,
		PIS:               "PIS-VALE",
		CTPF:              "CT-VALE",
		Nascimento:        time.Now().AddDate(-30, 0, 0),
		Admissao:          time.Now().AddDate(-2, 0, 0),
		Cargo:             "Assistente",
		SalarioInicial:    2200,
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

func newAdminAuthVale(lr *valeFakeLogRepo) *service.AuthService {
	cfg := service.AuthConfig{
		Issuer:          "autogrh-test",
		AccessTTL:       15 * time.Minute,
		ClockSkew:       2 * time.Minute,
		LoginSuccessEvt: 1001,
		LoginFailEvt:    1002,
		Timezone:        "America/Campo_Grande",
	}
	perms := service.PermissionMap{
		"admin": {"*": {}}, // libera "vale:update" e "vale:delete"
	}
	return service.NewAuthService(nil, lr, jwtm.NewHS256Manager([]byte("test-secret")), cfg, perms)
}

func newValeServiceWithDB(lr *valeFakeLogRepo) *service.ValeService {
	adp := Adapter.NewValeRepositoryAdapter(
		repository.CreateVale,
		repository.GetValeByID,
		repository.GetValesByFuncionarioID,
		repository.UpdateVale,
		nil, // SoftDelete (inexistente)
		repository.DeleteVale,
		repository.ListValesPendentes,
		repository.ListValesAprovadosNaoPagos,
		repository.ListAllVales,
	)
	return service.NewValeService(adp, newAdminAuthVale(lr), lr)
}

/************ TESTES ************/

func TestVale_Criar_Listar_Aprovar_Pagar_E_Listas(t *testing.T) {
	// DB limpo no início e no fim
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &valeFakeLogRepo{}
	svc := newValeServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 201, Perfil: "admin"}

	funcID := seedPessoaFuncionarioVale(t)

	// Criar vale (pendente)
	data := time.Date(2025, time.January, 10, 0, 0, 0, 0, time.Local)
	v, err := svc.CriarVale(ctx, claims, funcID, 500.00, data)
	if err != nil {
		t.Fatalf("CriarVale erro: %v", err)
	}
	if v == nil || v.ID == 0 || v.Aprovado || v.Pago == true {
		t.Fatalf("vale criado inválido: %+v", v)
	}

	// Listar todos (deve conter 1)
	all, err := svc.ListarVales(ctx, claims)
	if err != nil || len(all) != 1 {
		t.Fatalf("ListarVales esperado 1, got=%d err=%v", len(all), err)
	}

	// Listar pendentes (deve conter 1)
	pend, err := svc.ListarValesPendentes(ctx, claims)
	if err != nil || len(pend) != 1 || pend[0].ID != v.ID {
		t.Fatalf("ListarValesPendentes esperado 1 (id=%d), got=%d", v.ID, len(pend))
	}

	// Listar aprovados NÃO pagos (vazio)
	apNP, err := svc.ListarValesAprovadosNaoPagos(ctx, claims)
	if err != nil || len(apNP) != 0 {
		t.Fatalf("ListarValesAprovadosNaoPagos esperado 0, got=%d err=%v", len(apNP), err)
	}

	// Aprovar (continua não pago)
	if err := svc.AprovarVale(ctx, claims, v.ID); err != nil {
		t.Fatalf("AprovarVale erro: %v", err)
	}

	// Pendentes agora 0
	pend, _ = svc.ListarValesPendentes(ctx, claims)
	if len(pend) != 0 {
		t.Fatalf("pendentes esperado 0 após aprovar, got=%d", len(pend))
	}
	// Aprovados NÃO pagos agora 1
	apNP, _ = svc.ListarValesAprovadosNaoPagos(ctx, claims)
	if len(apNP) != 1 || apNP[0].ID != v.ID || !apNP[0].Aprovado || apNP[0].Pago {
		t.Fatalf("aprovadosNaoPagos esperado 1 (aprovado=true, pago=false): %+v", apNP)
	}

	// Marcar como PAGO
	if err := svc.MarcarValeComoPago(ctx, claims, v.ID); err != nil {
		t.Fatalf("MarcarValeComoPago erro: %v", err)
	}

	// Aprovados NÃO pagos volta a 0
	apNP, _ = svc.ListarValesAprovadosNaoPagos(ctx, claims)
	if len(apNP) != 0 {
		t.Fatalf("aprovadosNaoPagos esperado 0 após pagar, got=%d", len(apNP))
	}

	// Listar todos ainda contém o vale (agora aprovado e pago)
	all, _ = svc.ListarVales(ctx, claims)
	if len(all) != 1 || all[0].ID != v.ID || !all[0].Aprovado || !all[0].Pago {
		t.Fatalf("ListarVales após pagar inválido: %+v", all)
	}

	// Listar por funcionário
	byFunc, err := svc.ListarValesFuncionario(ctx, claims, funcID)
	if err != nil || len(byFunc) == 0 {
		t.Fatalf("ListarValesFuncionario esperado >=1, got=%d err=%v", len(byFunc), err)
	}
}

func TestVale_Update_Get_Delete(t *testing.T) {
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &valeFakeLogRepo{}
	svc := newValeServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 202, Perfil: "admin"}

	funcID := seedPessoaFuncionarioVale(t)

	// cria
	v, err := svc.CriarVale(ctx, claims, funcID, 900.00, time.Date(2025, time.February, 5, 0, 0, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("CriarVale erro: %v", err)
	}

	// update (valor/data)
	v.Valor = 950.00
	v.Data = time.Date(2025, time.March, 6, 0, 0, 0, 0, time.Local)
	if err := svc.AtualizarVale(ctx, claims, v); err != nil {
		t.Fatalf("AtualizarVale erro: %v", err)
	}

	// get
	got, err := svc.GetVale(ctx, claims, v.ID)
	if err != nil || got == nil || got.Valor != 950.00 || got.Data.Year() != 2025 || got.Data.Month() != time.March {
		t.Fatalf("GetVale inválido: %+v err=%v", got, err)
	}

	// delete definitivo
	if err := svc.DeleteVale(ctx, claims, v.ID); err != nil {
		t.Fatalf("DeleteVale erro: %v", err)
	}

	// sumiu das listagens
	all, _ := svc.ListarVales(ctx, claims)
	if len(all) != 0 {
		t.Fatalf("ListarVales esperado 0 após delete, got=%d", len(all))
	}
	byFunc, _ := svc.ListarValesFuncionario(ctx, claims, funcID)
	if len(byFunc) != 0 {
		t.Fatalf("ListarValesFuncionario esperado 0 após delete, got=%d", len(byFunc))
	}
}
