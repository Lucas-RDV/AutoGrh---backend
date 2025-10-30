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

/*** ---------------- Mocks (nomes exclusivos p/ Pessoa) ---------------- ***/

type pessoaFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *pessoaFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func pessoaHasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
	for _, e := range entries {
		if e.EventoID == evt && e.UsuarioID != nil && *e.UsuarioID == uid {
			if prefix == "" || strings.HasPrefix(e.Detalhe, prefix) {
				return true
			}
		}
	}
	return false
}

// pessoaFakePessoaRepo implementa service.PessoaRepository
type pessoaFakePessoaRepo struct {
	seq int64

	byID  map[int64]*entity.Pessoa
	byCPF map[string]*entity.Pessoa
	byRG  map[string]*entity.Pessoa
	list  []*entity.Pessoa

	forceExistsCPF bool
	forceExistsRG  bool

	errCreate error
	errGet    error
	errUpdate error
	errDelete error
}

func newPessoaFakeRepo() *pessoaFakePessoaRepo {
	return &pessoaFakePessoaRepo{
		byID:  make(map[int64]*entity.Pessoa),
		byCPF: make(map[string]*entity.Pessoa),
		byRG:  make(map[string]*entity.Pessoa),
	}
}

func (r *pessoaFakePessoaRepo) Create(ctx context.Context, p *entity.Pessoa) error {
	if r.errCreate != nil {
		return r.errCreate
	}
	r.seq++
	cp := *p
	cp.ID = r.seq
	r.byID[cp.ID] = &cp
	if cp.CPF != "" {
		r.byCPF[cp.CPF] = &cp
	}
	if cp.RG != "" {
		r.byRG[cp.RG] = &cp
	}
	r.list = append(r.list, &cp)
	p.ID = cp.ID // service espera o ID preenchido
	return nil
}

func (r *pessoaFakePessoaRepo) GetByID(ctx context.Context, id int64) (*entity.Pessoa, error) {
	if r.errGet != nil {
		return nil, r.errGet
	}
	p, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

func (r *pessoaFakePessoaRepo) GetByCPF(ctx context.Context, cpf string) (*entity.Pessoa, error) {
	if r.errGet != nil {
		return nil, r.errGet
	}
	if p, ok := r.byCPF[cpf]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, nil
}

func (r *pessoaFakePessoaRepo) Update(ctx context.Context, p *entity.Pessoa) error {
	if r.errUpdate != nil {
		return r.errUpdate
	}
	if p.ID <= 0 {
		return errors.New("id inválido no repo")
	}
	if _, ok := r.byID[p.ID]; !ok {
		return errors.New("não encontrado")
	}
	cp := *p
	r.byID[p.ID] = &cp
	if cp.CPF != "" {
		r.byCPF[cp.CPF] = &cp
	}
	if cp.RG != "" {
		r.byRG[cp.RG] = &cp
	}
	// atualiza em list (mantendo ordem)
	for i, x := range r.list {
		if x.ID == cp.ID {
			r.list[i] = &cp
			break
		}
	}
	return nil
}

func (r *pessoaFakePessoaRepo) Delete(ctx context.Context, id int64) error {
	if r.errDelete != nil {
		return r.errDelete
	}
	if _, ok := r.byID[id]; !ok {
		return nil
	}
	cp := r.byID[id]
	delete(r.byID, id)
	if cp.CPF != "" {
		delete(r.byCPF, cp.CPF)
	}
	if cp.RG != "" {
		delete(r.byRG, cp.RG)
	}
	// remove de list
	newList := make([]*entity.Pessoa, 0, len(r.list))
	for _, x := range r.list {
		if x.ID != id {
			newList = append(newList, x)
		}
	}
	r.list = newList
	return nil
}

func (r *pessoaFakePessoaRepo) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	if r.forceExistsCPF {
		return true, nil
	}
	_, ok := r.byCPF[cpf]
	return ok, nil
}

func (r *pessoaFakePessoaRepo) ExistsByRG(ctx context.Context, rg string) (bool, error) {
	if r.forceExistsRG {
		return true, nil
	}
	_, ok := r.byRG[rg]
	return ok, nil
}

func (r *pessoaFakePessoaRepo) SearchByNome(ctx context.Context, nome string) ([]*entity.Pessoa, error) {
	var out []*entity.Pessoa
	for _, p := range r.list {
		if strings.Contains(strings.ToLower(p.Nome), strings.ToLower(nome)) {
			cp := *p
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *pessoaFakePessoaRepo) ListAll(ctx context.Context) ([]*entity.Pessoa, error) {
	out := make([]*entity.Pessoa, 0, len(r.list))
	for _, p := range r.list {
		cp := *p
		out = append(out, &cp)
	}
	return out, nil
}

/*** -------------- SUT helpers (nomes exclusivos) -------------- ***/

func newAuthForPessoaSuite(lr *pessoaFakeLogRepo) *service.AuthService {
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

func newPessoaServiceSUT(repo service.PessoaRepository, lr *pessoaFakeLogRepo) *service.PessoaService {
	auth := newAuthForPessoaSuite(lr)
	return service.NewPessoaService(auth, lr, repo)
}

/*** ---------------- Tests ---------------- ***/

func TestPessoa_CreatePessoa_Sucesso(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 99}

	// dados com espaçamentos para validar Trim
	p := &entity.Pessoa{
		Nome:              "  Ana Maria  ",
		CPF:               "  12345678901 ",
		RG:                "  998877 ",
		Endereco:          "Rua A, 100",
		Contato:           "99999-0000",
		ContatoEmergencia: "88888-0000",
	}
	if err := svc.CreatePessoa(ctx, claims, p); err != nil {
		t.Fatalf("CreatePessoa erro: %v", err)
	}
	if p.ID == 0 {
		t.Fatalf("esperava ID após create")
	}

	// verificação: persistiu TRIM
	got, _ := repo.GetByID(ctx, p.ID)
	if got.Nome != "Ana Maria" || got.CPF != "12345678901" || got.RG != "998877" {
		t.Errorf("esperava campos TRIM: %+v", got)
	}

	// log de criação (EventoID=3)
	if !pessoaHasLogPrefix(lr.entries, 3, 99, "Criou pessoa ID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}
}

func TestPessoa_CreatePessoa_DuplicadoCPF_RG(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	// primeiro insert
	if err := svc.CreatePessoa(ctx, claims, &entity.Pessoa{
		Nome: "X", CPF: "111", RG: "222",
	}); err != nil {
		t.Fatalf("primeiro create falhou: %v", err)
	}

	// duplicado por CPF
	err := svc.CreatePessoa(ctx, claims, &entity.Pessoa{
		Nome: "Y", CPF: "111", RG: "333",
	})
	if err == nil || !strings.Contains(err.Error(), "já existe uma pessoa com este CPF") {
		t.Fatalf("esperava erro de CPF duplicado, veio: %v", err)
	}

	// duplicado por RG
	err = svc.CreatePessoa(ctx, claims, &entity.Pessoa{
		Nome: "Z", CPF: "444", RG: "222",
	})
	if err == nil || !strings.Contains(err.Error(), "já existe uma pessoa com este RG") {
		t.Fatalf("esperava erro de RG duplicado, veio: %v", err)
	}
}

func TestPessoa_CreatePessoa_CamposObrigatorios(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	cases := []struct {
		name string
		p    *entity.Pessoa
		want string
	}{
		{"Nome vazio", &entity.Pessoa{Nome: " ", CPF: "1", RG: "2"}, "nome não pode ser vazio"},
		{"CPF vazio", &entity.Pessoa{Nome: "A", CPF: " ", RG: "2"}, "CPF não pode ser vazio"},
		{"RG vazio", &entity.Pessoa{Nome: "A", CPF: "1", RG: " "}, "RG não pode ser vazio"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.CreatePessoa(ctx, claims, tc.p)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("esperava erro %q, veio: %v", tc.want, err)
			}
		})
	}
}

func TestPessoa_GetPessoaByID(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	// inválido
	if _, err := svc.GetPessoaByID(ctx, claims, 0); err == nil {
		t.Fatalf("esperava erro por ID inválido")
	}

	// cria e busca
	p := &entity.Pessoa{Nome: "Ana", CPF: "1", RG: "2"}
	if err := svc.CreatePessoa(ctx, claims, p); err != nil {
		t.Fatalf("create falhou: %v", err)
	}
	got, err := svc.GetPessoaByID(ctx, claims, p.ID)
	if err != nil || got == nil || got.ID != p.ID {
		t.Fatalf("GetByID falhou: got=%+v err=%v", got, err)
	}
}

func TestPessoa_UpdatePessoa_Sucesso_E_Log(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 7}

	// cria base
	p := &entity.Pessoa{Nome: "Ana", CPF: "1", RG: "2", Endereco: "R1", Contato: "C1", ContatoEmergencia: "E1"}
	if err := svc.CreatePessoa(ctx, claims, p); err != nil {
		t.Fatalf("create falhou: %v", err)
	}

	// atualiza com TRIM e novos valores
	up := &entity.Pessoa{
		ID:                p.ID,
		Nome:              "  Ana Maria ",
		CPF:               "  111 ",
		RG:                "  222 ",
		Endereco:          "R2",
		Contato:           "C2",
		ContatoEmergencia: "E2",
	}
	if err := svc.UpdatePessoa(ctx, claims, up); err != nil {
		t.Fatalf("update falhou: %v", err)
	}

	// verificação persistida (TRIM aplicado)
	got, _ := repo.GetByID(ctx, p.ID)
	if got.Nome != "Ana Maria" || got.CPF != "111" || got.RG != "222" || got.Endereco != "R2" || got.Contato != "C2" || got.ContatoEmergencia != "E2" {
		t.Fatalf("update não aplicou corretamente: %+v", got)
	}

	// log de update (EventoID=4)
	if !pessoaHasLogPrefix(lr.entries, 4, 7, "Atualizou pessoa ID=") {
		t.Errorf("não registrou log de update (EventoID=4)")
	}
}

func TestPessoa_UpdatePessoa_Validacoes(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	// ID inválido
	err := svc.UpdatePessoa(ctx, claims, &entity.Pessoa{ID: 0, Nome: "A", CPF: "1", RG: "2"})
	if err == nil || !strings.Contains(err.Error(), "ID inválido") {
		t.Fatalf("esperava erro por ID inválido, veio: %v", err)
	}
	// Vazios
	err = svc.UpdatePessoa(ctx, claims, &entity.Pessoa{ID: 1, Nome: " ", CPF: "1", RG: "2"})
	if err == nil || !strings.Contains(err.Error(), "não podem ser vazios") {
		t.Fatalf("esperava erro por nome vazio, veio: %v", err)
	}
	err = svc.UpdatePessoa(ctx, claims, &entity.Pessoa{ID: 1, Nome: "A", CPF: " ", RG: "2"})
	if err == nil || !strings.Contains(err.Error(), "não podem ser vazios") {
		t.Fatalf("esperava erro por CPF vazio, veio: %v", err)
	}
	err = svc.UpdatePessoa(ctx, claims, &entity.Pessoa{ID: 1, Nome: "A", CPF: "1", RG: " "})
	if err == nil || !strings.Contains(err.Error(), "não podem ser vazios") {
		t.Fatalf("esperava erro por RG vazio, veio: %v", err)
	}
}

func TestPessoa_DeletePessoa_Sucesso_E_Log(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 55}

	// cria
	p := &entity.Pessoa{Nome: "Ana", CPF: "1", RG: "2"}
	if err := svc.CreatePessoa(ctx, claims, p); err != nil {
		t.Fatalf("create falhou: %v", err)
	}

	// deleta
	if err := svc.DeletePessoa(ctx, claims, p.ID); err != nil {
		t.Fatalf("delete falhou: %v", err)
	}
	if got, _ := repo.GetByID(ctx, p.ID); got != nil {
		t.Fatalf("esperava pessoa removida, ainda existe: %+v", got)
	}
	// log de delete (EventoID=5)
	if !pessoaHasLogPrefix(lr.entries, 5, 55, "Deletou pessoa ID=") {
		t.Errorf("não registrou log de delete (EventoID=5)")
	}
}

func TestPessoa_DeletePessoa_ValidacaoID(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	if err := svc.DeletePessoa(ctx, claims, 0); err == nil {
		t.Fatalf("esperava erro por ID inválido")
	}
}

func TestPessoa_ListPessoas_RetornaTudo(t *testing.T) {
	defer func() { _ = truncateAll() }()

	lr := &pessoaFakeLogRepo{}
	repo := newPessoaFakeRepo()
	svc := newPessoaServiceSUT(repo, lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 1}

	_ = svc.CreatePessoa(ctx, claims, &entity.Pessoa{Nome: "A", CPF: "1", RG: "2"})
	_ = svc.CreatePessoa(ctx, claims, &entity.Pessoa{Nome: "B", CPF: "3", RG: "4"})
	_ = svc.CreatePessoa(ctx, claims, &entity.Pessoa{Nome: "C", CPF: "5", RG: "6"})

	list, err := svc.ListPessoas(ctx, claims)
	if err != nil {
		t.Fatalf("ListPessoas erro: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("esperava 3 pessoas, veio %d", len(list))
	}
}
