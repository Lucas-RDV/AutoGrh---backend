package testes

import (
	Adapter "AutoGRH/pkg/adapter"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/service/jwt"
)

/************ Log fake ************/

type docFakeLogRepo struct {
	entries []service.LogEntry
}

func (l *docFakeLogRepo) Create(ctx context.Context, e service.LogEntry) (int64, error) {
	l.entries = append(l.entries, e)
	return int64(len(l.entries)), nil
}

func docHasLogPrefix(entries []service.LogEntry, evt int64, uid int64, prefix string) bool {
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

// Cria Pessoa (CPF 11 dígitos) + Funcionário para satisfazer FK
func seedPessoaFuncionarioDoc(t *testing.T) int64 {
	t.Helper()
	now := time.Now().UnixNano()
	cpf := fmt.Sprintf("%011d", now%100000000000)
	rg := fmt.Sprintf("%09d", now%1000000000)

	p := &entity.Pessoa{Nome: "Teste Documento", CPF: cpf, RG: rg}
	if err := repository.CreatePessoa(p); err != nil {
		t.Fatalf("seed CreatePessoa erro: %v", err)
	}
	if p.ID == 0 {
		t.Fatalf("seed pessoa sem ID")
	}

	f := &entity.Funcionario{
		PessoaID:          p.ID,
		PIS:               "PIS-DOC",
		CTPF:              "CT-DOC",
		Nascimento:        time.Now().AddDate(-25, 0, 0),
		Admissao:          time.Now().AddDate(-1, 0, 0),
		Cargo:             "Analista",
		SalarioInicial:    2500,
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

func newAdminAuthDoc(lr *docFakeLogRepo) *service.AuthService {
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

func newDocumentoServiceWithDB(lr *docFakeLogRepo) *service.DocumentoService {
	auth := newAdminAuthDoc(lr)
	adapter := Adapter.NewDocumentoRepositoryAdapter(
		repository.CreateDocumento,
		repository.GetDocumentosByFuncionarioID,
		repository.GetByID,
		repository.ListDocumentos,
		repository.DeleteDocumento,
	)
	return service.NewDocumentoService(auth, lr, adapter)
}

// cria um arquivo temporário e reabre para leitura (multipart.File requer *os.File)
func createTempUploadFile(t *testing.T, name, content string) *os.File {
	t.Helper()
	tmpDir := t.TempDir()
	full := filepath.Join(tmpDir, name)
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file erro: %v", err)
	}
	f, err := os.Open(full)
	if err != nil {
		t.Fatalf("open temp file erro: %v", err)
	}
	return f
}

/************ TESTES ************/

func TestDocumento_Salvar_Listar_DownloadPath_Delete(t *testing.T) {
	defer func() { _ = truncateAll() }()

	// Isola o "home" do Windows para este teste
	home := t.TempDir()
	t.Setenv("USERPROFILE", home) // Windows prioriza USERPROFILE
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	lr := &docFakeLogRepo{}
	svc := newDocumentoServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 90, Perfil: "admin"}

	funcID := seedPessoaFuncionarioDoc(t)

	// Prepara diretórios que o service espera: ~/Documents/AutoGRH/documentos/{funcID}
	baseDir := filepath.Join(home, "Documents", "AutoGRH")
	destDir := filepath.Join(baseDir, "documentos", strconv.FormatInt(funcID, 10))
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatalf("prep dirs erro: %v", err)
	}

	// --- SalvarDocumento ---
	up := createTempUploadFile(t, "comprovante.pdf", "PDF FAKE CONTENT")
	defer up.Close()

	doc, err := svc.SalvarDocumento(ctx, claims, funcID, up, "comprovante.pdf")
	if err != nil {
		t.Fatalf("SalvarDocumento erro: %v", err)
	}
	if doc == nil || doc.ID == 0 {
		t.Fatalf("esperava documento criado com ID")
	}
	if doc.FuncionarioID != funcID || doc.Caminho == "" {
		t.Fatalf("documento fora do esperado: %+v", doc)
	}
	if !docHasLogPrefix(lr.entries, 3, 90, "Documento criado funcionarioID=") {
		t.Errorf("não registrou log de criação (EventoID=3)")
	}

	// Caminho físico deve existir em ~/Documents/AutoGRH/<doc.Caminho>
	fullPath := filepath.Join(baseDir, doc.Caminho)
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("arquivo não foi gravado no disco em: %s (err=%v)", fullPath, err)
	}

	// --- GetDocumentosByFuncionarioID ---
	listByFunc, err := svc.GetDocumentosByFuncionarioID(ctx, claims, funcID)
	if err != nil {
		t.Fatalf("GetDocumentosByFuncionarioID erro: %v", err)
	}
	if len(listByFunc) == 0 {
		t.Fatalf("esperava ao menos 1 documento do funcionário")
	}

	// --- ListDocumentos (geral) ---
	listAll, err := svc.ListDocumentos(ctx, claims)
	if err != nil {
		t.Fatalf("ListDocumentos erro: %v", err)
	}
	if len(listAll) == 0 {
		t.Fatalf("esperava ao menos 1 documento no sistema")
	}

	// --- GetDocumentoPath ---
	gotPath, err := svc.GetDocumentoPath(ctx, claims, doc.ID)
	if err != nil {
		t.Fatalf("GetDocumentoPath erro: %v", err)
	}
	if gotPath != fullPath {
		t.Fatalf("GetDocumentoPath retornou caminho diferente:\n got: %s\nwant: %s", gotPath, fullPath)
	}
	if data, err := os.ReadFile(gotPath); err != nil || len(data) == 0 {
		t.Fatalf("arquivo de download não legível: err=%v len=%d", err, len(data))
	}

	// --- DeleteDocumento ---
	if err := svc.DeleteDocumento(ctx, claims, doc.ID); err != nil {
		t.Fatalf("DeleteDocumento erro: %v", err)
	}
	// registro some do banco
	dcheck, err := repository.GetByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("GetByID pós-delete erro: %v", err)
	}
	if dcheck != nil {
		t.Fatalf("documento ainda existe no banco após delete: %+v", dcheck)
	}
	// arquivo some do disco (ou foi ignorado se já inexistente)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		if err == nil {
			t.Fatalf("arquivo ainda existe no disco após delete: %s", fullPath)
		}
	}

	if !docHasLogPrefix(lr.entries, 5, 90, "Documento deletado id=") {
		t.Errorf("não registrou log de deleção (EventoID=5)")
	}
}

// Smoke: upload de arquivo "grande" em memória (apenas para exercitar io.Copy com buffer maior)
func TestDocumento_Upload_BufferMaior(t *testing.T) {
	defer func() { _ = truncateAll() }()

	// Isola o "home" do Windows para este teste
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	lr := &docFakeLogRepo{}
	svc := newDocumentoServiceWithDB(lr)
	ctx := context.Background()
	claims := service.Claims{UserID: 91, Perfil: "admin"}

	funcID := seedPessoaFuncionarioDoc(t)

	// Prepara diretórios esperados pelo service
	baseDir := filepath.Join(home, "Documents", "AutoGRH")
	destDir := filepath.Join(baseDir, "documentos", strconv.FormatInt(funcID, 10))
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatalf("prep dirs erro: %v", err)
	}

	// cria arquivo temporário com ~1MB
	var buf bytes.Buffer
	chunk := bytes.Repeat([]byte("X"), 1024) // 1KB
	for i := 0; i < 1024; i++ {              // ~1MB
		if _, err := buf.Write(chunk); err != nil {
			t.Fatalf("prep buffer erro: %v", err)
		}
	}
	tmp := t.TempDir()
	path := filepath.Join(tmp, "grande.bin")
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		t.Fatalf("write grande erro: %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open grande erro: %v", err)
	}
	defer f.Close()

	doc, err := svc.SalvarDocumento(ctx, claims, funcID, f, "grande.bin")
	if err != nil {
		t.Fatalf("SalvarDocumento(grande) erro: %v", err)
	}
	if doc == nil || doc.ID == 0 {
		t.Fatalf("esperava doc criado com ID")
	}

	// valida leitura
	full := filepath.Join(baseDir, doc.Caminho)
	gf, err := os.Open(full)
	if err != nil {
		t.Fatalf("open saved erro: %v", err)
	}
	defer gf.Close()

	n, err := io.Copy(io.Discard, gf)
	if err != nil || n == 0 {
		t.Fatalf("erro lendo arquivo salvo: n=%d err=%v", n, err)
	}
}
