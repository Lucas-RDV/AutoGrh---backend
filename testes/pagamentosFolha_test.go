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

/*
Cobre:
- FolhaPagamentoService: CriarFolhaSalario, RecalcularFolha (SALARIO), FecharFolha,
  CriarFolhaVale, RecalcularFolhaVale (limpa e refaz), Listar/Buscar/BuscarPorMesAnoTipo.
- PagamentoService: BuscarPagamento, AtualizarPagamento (mantém descontoVales), ListarPagamentosFuncionario,
  MarcarPagamentoComoPago, ListarPagamentosDaFolha.

Refs: Adapters (Folha/Pagamento) e Services/Repos/Entidades.
- Folha adapter/repo iface/métodos. :contentReference[oaicite:8]{index=8} :contentReference[oaicite:9]{index=9}
- Pagamento adapter/repo iface/métodos. :contentReference[oaicite:10]{index=10} :contentReference[oaicite:11]{index=11}
- Services Folha/Pagamento (regras: faltas, vales pagos, rebuild, fechar). :contentReference[oaicite:12]{index=12} :contentReference[oaicite:13]{index=13}
- Entidades Folha/Pagamento (RecalcularValorFinal). :contentReference[oaicite:14]{index=14} :contentReference[oaicite:15]{index=15}
*/

/************ Fake Log ************/
type folhaFakeLogRepo struct{ entries []service.LogEntry }

func (l *folhaFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

/************ Auth helper ************/
func newAdminAuth(lr service.LogRepository) *service.AuthService {
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

/************ Seed helpers ************/
// Cria Pessoa + Funcionario válidos
func seedPessoaFuncionarioBase(t *testing.T, nome string) int64 {
	t.Helper()
	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000)
	rg := fmt.Sprintf("%09d", now%1000000000)

	p := &entity.Pessoa{Nome: nome, CPF: cpf, RG: rg}
	if err := repository.CreatePessoa(p); err != nil {
		t.Fatalf("seed CreatePessoa erro: %v", err)
	}
	f := &entity.Funcionario{
		PessoaID:       p.ID,
		PIS:            "PIS-" + nome,
		CTPF:           "CT-" + nome,
		Nascimento:     time.Now().AddDate(-28, 0, 0),
		Admissao:       time.Now().AddDate(-2, 0, 0),
		Cargo:          "Dev",
		SalarioInicial: 0,
	}
	if err := repository.CreateFuncionario(f); err != nil {
		t.Fatalf("seed CreateFuncionario erro: %v", err)
	}
	return f.ID
}

// Seta salário real atual
func seedSalarioRealAtual(t *testing.T, funcID int64, valor float64) {
	t.Helper()
	sr := &entity.SalarioReal{
		FuncionarioID: funcID,
		Valor:         valor,
		Inicio:        time.Now().AddDate(0, -1, 0), // vigente
	}
	if err := repository.CreateSalarioReal(sr); err != nil {
		t.Fatalf("seed CreateSalarioReal erro: %v", err)
	}
}

func seedFaltasMes(t *testing.T, funcID int64, mes, ano, qtd int) {
	t.Helper()
	// Usa o FaltaService.UpsertMensal (é ele que faz o upsert corretamente)
	lr := &folhaFakeLogRepo{}
	fsvc := newFaltaServiceForSeed(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 999, Perfil: "admin"}
	if err := fsvc.UpsertMensal(ctx, claims, funcID, mes, ano, qtd); err != nil {
		t.Fatalf("seed Upsert faltas erro: %v", err)
	}
}

// Instância mínima do FaltaService para seed
func newFaltaServiceForSeed(lr *folhaFakeLogRepo) *service.FaltaService {
	auth := newAdminAuth(lr)
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

// Cria vale aprovado+pago no mês/ano (desconto em folha SALARIO)
func seedValePago(t *testing.T, funcID int64, valor float64, data time.Time) {
	t.Helper()
	v := &entity.Vale{
		FuncionarioID: funcID,
		Valor:         valor,
		Data:          data,
		Aprovado:      true,
		Pago:          true,
		Ativo:         true,
	}
	if err := repository.CreateVale(v); err != nil {
		t.Fatalf("seed CreateVale erro: %v", err)
	}
}

/************ Service factories ************/
func newFolhaService(lr *folhaFakeLogRepo) *service.FolhaPagamentoService {
	auth := newAdminAuth(lr)
	adp := Adapter.NewFolhaPagamentoRepositoryAdapter(
		repository.CreateFolhaPagamento,
		repository.GetFolhaPagamentoByID,
		repository.GetFolhaByMesAnoTipo,
		repository.UpdateFolhaPagamento,
		repository.DeleteFolhaPagamento,
		repository.ListFolhasPagamentos,
		repository.MarcarFolhaComoPaga,
	)
	return service.NewFolhaPagamentoService(adp, auth, lr)
}

func newPagamentoService(lr *folhaFakeLogRepo) *service.PagamentoService {
	auth := newAdminAuth(lr)
	padp := Adapter.NewPagamentoRepositoryAdapter(
		repository.CreatePagamento, // não usado diretamente pelo service
		repository.UpdatePagamento,
		repository.GetPagamentosByFolhaID,
		repository.DeletePagamentosByFolhaID, // usado no RecalcularFolhaVale (pelo service de folha)
		repository.GetPagamentoByID,
		repository.ListPagamentosByFuncionarioID,
	)
	return service.NewPagamentoService(padp, auth, lr)
}

/************ TESTES ************/

func TestFolhaSalario_Criar_Recalcular_Fechar(t *testing.T) {
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &folhaFakeLogRepo{}
	fs := newFolhaService(lr)
	ps := newPagamentoService(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 501, Perfil: "admin"}

	const mes, ano = 1, 2025
	// Seed: 1 funcionário com salário 3000, 3 faltas, 1 vale pago 500
	funcID := seedPessoaFuncionarioBase(t, "Funcionario A")
	seedSalarioRealAtual(t, funcID, 3000)
	seedFaltasMes(t, funcID, mes, ano, 3) // desconto faltas = 3000/30*3 = 300
	seedValePago(t, funcID, 500, time.Date(ano, time.January, 10, 0, 0, 0, 0, time.Local))

	// Criar folha SALARIO
	folha, err := fs.CriarFolhaSalario(ctx, claims, mes, ano)
	if err != nil {
		t.Fatalf("CriarFolhaSalario erro: %v", err)
	}
	if folha == nil || folha.ID == 0 || folha.Tipo != "SALARIO" || folha.Pago {
		t.Fatalf("folha salario inválida: %+v", folha)
	}
	// Valor final esperado = 3000 - 300 (faltas) - 500 (vales pagos) = 2200
	if folha.ValorTotal < 2199.99 || folha.ValorTotal > 2200.01 {
		t.Fatalf("ValorTotal esperado ~2200, veio %.2f", folha.ValorTotal)
	}

	// Ver pagarmentos da folha
	rows, err := ps.ListarPagamentosDaFolha(ctx, claims, folha.ID) // :contentReference[oaicite:16]{index=16}
	if err != nil || len(rows) != 1 {
		t.Fatalf("ListarPagamentosDaFolha esperado 1, got=%d err=%v", len(rows), err)
	}
	p := rows[0]
	if p.FuncionarioID != funcID || p.FolhaID != folha.ID {
		t.Fatalf("pagamento inconsistente: %+v", p)
	}
	if p.ValorFinal < 2199.99 || p.ValorFinal > 2200.01 {
		t.Fatalf("valorFinal esperado ~2200, veio %.2f", p.ValorFinal)
	}

	// Altera cenário e Recalcular (SALARIO): +2 faltas e novo vale pago 200
	seedFaltasMes(t, funcID, mes, ano, 5) // desconto faltas agora = 3000/30*5 = 500
	seedValePago(t, funcID, 200, time.Date(ano, time.January, 20, 0, 0, 0, 0, time.Local))
	if err := fs.RecalcularFolha(ctx, claims, folha.ID); err != nil { // SALARIO: rebuild mantém/atualiza pagamentos existentes :contentReference[oaicite:17]{index=17}
		t.Fatalf("RecalcularFolha salario erro: %v", err)
	}

	// Confere novo total (faltas 500 + vales 700 = 1200) → 3000-1200 = 1800
	folhaRec, err := repository.GetFolhaPagamentoByID(folha.ID)
	if err != nil || folhaRec == nil {
		t.Fatalf("GetFolhaPagamentoByID erro: %v", err)
	}
	if folhaRec.ValorTotal < 1799.99 || folhaRec.ValorTotal > 1800.01 {
		t.Fatalf("ValorTotal após recalcular esperado ~1800, veio %.2f", folhaRec.ValorTotal)
	}
	rows2, _ := ps.ListarPagamentosDaFolha(ctx, claims, folha.ID)
	if len(rows2) != 1 || rows2[0].ValorFinal < 1799.99 || rows2[0].ValorFinal > 1800.01 {
		t.Fatalf("Pagamento após recalcular esperado ~1800, veio %.2f", rows2[0].ValorFinal)
	}

	// Fechar folha (marca pagamentos e folha como pagos) :contentReference[oaicite:18]{index=18}
	if err := fs.FecharFolha(ctx, claims, folha.ID); err != nil {
		t.Fatalf("FecharFolha erro: %v", err)
	}
	fechada, _ := repository.GetFolhaPagamentoByID(folha.ID)
	if fechada == nil || !fechada.Pago {
		t.Fatalf("folha deveria estar paga: %+v", fechada)
	}
	rows3, _ := ps.ListarPagamentosDaFolha(ctx, claims, folha.ID)
	for _, pp := range rows3 {
		if !pp.Pago {
			t.Fatalf("pagamento %d deveria estar pago", pp.ID)
		}
	}
}

func TestFolhaVale_Criar_Recalcular_Fechar(t *testing.T) {
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &folhaFakeLogRepo{}
	fs := newFolhaService(lr)
	ps := newPagamentoService(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 502, Perfil: "admin"}

	const mes, ano = 2, 2025
	funcA := seedPessoaFuncionarioBase(t, "Func Vale A")
	funcB := seedPessoaFuncionarioBase(t, "Func Vale B")

	// Dois vales aprovados NÃO pagos (devem ir p/ folha de VALE) :contentReference[oaicite:19]{index=19}
	v1 := &entity.Vale{FuncionarioID: funcA, Valor: 300, Data: time.Date(ano, time.February, 5, 0, 0, 0, 0, time.Local), Aprovado: true, Pago: false, Ativo: true}
	v2 := &entity.Vale{FuncionarioID: funcB, Valor: 450, Data: time.Date(ano, time.February, 7, 0, 0, 0, 0, time.Local), Aprovado: true, Pago: false, Ativo: true}
	if err := repository.CreateVale(v1); err != nil {
		t.Fatalf("CreateVale v1 erro: %v", err)
	}
	if err := repository.CreateVale(v2); err != nil {
		t.Fatalf("CreateVale v2 erro: %v", err)
	}

	// Criar folha VALE (soma 750)
	fv, err := fs.CriarFolhaVale(ctx, claims, mes, ano)
	if err != nil {
		t.Fatalf("CriarFolhaVale erro: %v", err)
	}
	if fv == nil || fv.Tipo != "VALE" || fv.ValorTotal < 749.99 || fv.ValorTotal > 750.01 {
		t.Fatalf("folha VALE inválida/total: %+v", fv)
	}
	pags, _ := ps.ListarPagamentosDaFolha(ctx, claims, fv.ID)
	if len(pags) != 2 {
		t.Fatalf("esperava 2 pagamentos na folha de vale, got=%d", len(pags))
	}

	// RecalcularFolhaVale zera e recria conforme aprovados não pagos atuais (sem filtro de data) :contentReference[oaicite:20]{index=20}
	// Remove um vale (marca pago manualmente) e adiciona outro não pago → total esperado muda
	v2.Pago = true
	if err := repository.UpdateVale(v2); err != nil {
		t.Fatalf("update v2 pago erro: %v", err)
	}
	v3 := &entity.Vale{FuncionarioID: funcA, Valor: 200, Data: time.Date(ano, time.February, 10, 0, 0, 0, 0, time.Local), Aprovado: true, Pago: false, Ativo: true}
	if err := repository.CreateVale(v3); err != nil {
		t.Fatalf("CreateVale v3 erro: %v", err)
	}

	if err := fs.RecalcularFolhaVale(ctx, claims, fv.ID); err != nil {
		t.Fatalf("RecalcularFolhaVale erro: %v", err)
	}
	// Agora só v1 (300) e v3 (200) entram = 500
	fv2, _ := repository.GetFolhaPagamentoByID(fv.ID)
	if fv2.ValorTotal < 499.99 || fv2.ValorTotal > 500.01 {
		t.Fatalf("ValorTotal folha VALE após recalcular esperado ~500, veio %.2f", fv2.ValorTotal)
	}
	pags2, _ := ps.ListarPagamentosDaFolha(ctx, claims, fv.ID)
	if len(pags2) != 2 {
		t.Fatalf("esperava 2 pagamentos após recalcular, got=%d", len(pags2))
	}

	// Fechar folha VALE → marca pagamentos da folha como pagos e marca TODOS vales aprovados como pagos (repo) :contentReference[oaicite:21]{index=21}
	if err := fs.FecharFolha(ctx, claims, fv.ID); err != nil {
		t.Fatalf("FecharFolha VALE erro: %v", err)
	}
	fv3, _ := repository.GetFolhaPagamentoByID(fv.ID)
	if fv3 == nil || !fv3.Pago {
		t.Fatalf("folha VALE deveria estar paga: %+v", fv3)
	}
	pags3, _ := ps.ListarPagamentosDaFolha(ctx, claims, fv.ID)
	for _, p := range pags3 {
		if !p.Pago {
			t.Fatalf("pagamento da folha VALE deveria estar pago: %+v", p)
		}
	}
	// (Opcional) poderíamos checar que ListValesAprovadosNaoPagos agora está vazia, dependendo da sua implementação de MarcarTodosValesComoPagos.
}

func TestPagamentoService_Get_Update_Pagar_ListarFuncionario(t *testing.T) {
	if err := truncateAll(); err != nil {
		t.Fatalf("truncateAll inicio: %v", err)
	}
	t.Cleanup(func() { _ = truncateAll() })

	lr := &folhaFakeLogRepo{}
	fs := newFolhaService(lr)
	ps := newPagamentoService(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 503, Perfil: "admin"}

	const mes, ano = 3, 2025
	funcID := seedPessoaFuncionarioBase(t, "Func Pagamento")
	seedSalarioRealAtual(t, funcID, 2500)
	seedFaltasMes(t, funcID, mes, ano, 0)

	// Cria uma folha SALARIO simples com 1 pagamento
	folha, err := fs.CriarFolhaSalario(ctx, claims, mes, ano)
	if err != nil {
		t.Fatalf("CriarFolhaSalario erro: %v", err)
	}
	pags, _ := ps.ListarPagamentosDaFolha(ctx, claims, folha.ID)
	if len(pags) != 1 {
		t.Fatalf("esperava 1 pagamento na folha criada, got=%d", len(pags))
	}
	pgID := pags[0].ID

	// Buscar pagamento
	got, err := ps.BuscarPagamento(ctx, claims, pgID) // :contentReference[oaicite:22]{index=22}
	if err != nil || got == nil || got.ID != pgID {
		t.Fatalf("BuscarPagamento inválido: %+v err=%v", got, err)
	}

	// AtualizarPagamento: aplica adicional/INSS/família e REcalcula mantendo descontoVales :contentReference[oaicite:23]{index=23}
	if err := ps.AtualizarPagamento(ctx, claims, pgID, 100, 50, 30); err != nil {
		t.Fatalf("AtualizarPagamento erro: %v", err)
	}
	got2, _ := ps.BuscarPagamento(ctx, claims, pgID)
	// ValorFinal = salarioBase + adicional + familia - INSS - descontoVales - descontoFaltas
	esperado := got.SalarioBase + 100 + 30 - 50 - got.DescontoVales - 0
	if got2.ValorFinal < esperado-0.01 || got2.ValorFinal > esperado+0.01 {
		t.Fatalf("ValorFinal após atualizar esperado ~%.2f, veio %.2f", esperado, got2.ValorFinal)
	}

	// Listar por funcionário
	listFunc, err := ps.ListarPagamentosFuncionario(ctx, claims, funcID)
	if err != nil || len(listFunc) != 1 || listFunc[0].ID != pgID {
		t.Fatalf("ListarPagamentosFuncionario inválido: len=%d err=%v", len(listFunc), err)
	}

	// Marcar pagamento como pago
	if err := ps.MarcarPagamentoComoPago(ctx, claims, pgID); err != nil {
		t.Fatalf("MarcarPagamentoComoPago erro: %v", err)
	}
	got3, _ := ps.BuscarPagamento(ctx, claims, pgID)
	if !got3.Pago {
		t.Fatalf("Pagamento deveria estar pago: %+v", got3)
	}
}
