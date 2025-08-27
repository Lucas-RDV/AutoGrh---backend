package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// DocumentoRepository define as operações necessárias para documentos
type DocumentoRepository interface {
	Create(ctx context.Context, d *entity.Documento) error
	GetByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error)
	GetByID(ctx context.Context, id int64) (*entity.Documento, error) // novo
	List(ctx context.Context) ([]*entity.Documento, error)
	Delete(ctx context.Context, id int64) error
}

type DocumentoService struct {
	authService *AuthService
	logRepo     LogRepository
	docRepo     DocumentoRepository
}

func NewDocumentoService(auth *AuthService, logRepo LogRepository, docRepo DocumentoRepository) *DocumentoService {
	return &DocumentoService{
		authService: auth,
		logRepo:     logRepo,
		docRepo:     docRepo,
	}
}

func getBaseDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// fallback para diretório atual se não conseguir obter home
		home = "."
	}
	return filepath.Join(home, "Documents", "AutoGRH")
}

// SalvarDocumento grava o arquivo no disco e persiste o caminho relativo no banco
func (s *DocumentoService) SalvarDocumento(
	ctx context.Context,
	claims Claims,
	funcionarioID int64,
	file multipart.File,
	originalName string,
) (*entity.Documento, error) {

	if err := s.authService.Authorize(ctx, claims, "documento:create"); err != nil {
		return nil, err
	}

	// Base padrão: ~/Documents/AutoGRH
	baseDir := getBaseDir()

	// Diretório relativo e absoluto do funcionário
	relDir := filepath.Join("documentos", fmt.Sprintf("%d", funcionarioID))
	dir := filepath.Join(baseDir, relDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de documentos: %w", err)
	}

	// Nome final do arquivo: timestamp_nomeOriginal
	nomeArquivo := fmt.Sprintf("%d_%s", time.Now().Unix(), originalName)

	// Caminhos relativo (banco) e absoluto (disco)
	relPath := filepath.Join(relDir, nomeArquivo)
	fullPath := filepath.Join(baseDir, relPath)

	// Salvar arquivo em disco
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo em disco: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, fmt.Errorf("erro ao gravar arquivo em disco: %w", err)
	}

	// Criar entidade e salvar no banco (somente caminho relativo)
	d := &entity.Documento{
		FuncionarioID: funcionarioID,
		Caminho:       relPath,
	}
	if err := s.docRepo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("erro ao salvar documento no banco: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Documento criado funcionarioID=%d caminho=%s", funcionarioID, relPath),
	})

	return d, nil
}

// GetDocumentosByFuncionarioID retorna documentos de um funcionário
func (s *DocumentoService) GetDocumentosByFuncionarioID(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Documento, error) {
	if err := s.authService.Authorize(ctx, claims, "documento:list"); err != nil {
		return nil, err
	}
	return s.docRepo.GetByFuncionarioID(ctx, funcionarioID)
}

// ListDocumentos retorna todos os documentos
func (s *DocumentoService) ListDocumentos(ctx context.Context, claims Claims) ([]*entity.Documento, error) {
	if err := s.authService.Authorize(ctx, claims, "documento:list"); err != nil {
		return nil, err
	}
	return s.docRepo.List(ctx)
}

// GetDocumentoPath retorna o caminho absoluto de um documento para download
func (s *DocumentoService) GetDocumentoPath(ctx context.Context, claims Claims, id int64) (string, error) {
	if err := s.authService.Authorize(ctx, claims, "documento:list"); err != nil {
		return "", err
	}

	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar documento: %w", err)
	}
	if doc == nil {
		return "", fmt.Errorf("documento não encontrado")
	}

	baseDir := getBaseDir()
	fullPath := filepath.Join(baseDir, doc.Caminho)
	return fullPath, nil
}

// DeleteDocumento remove documento do banco e do disco
func (s *DocumentoService) DeleteDocumento(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "documento:delete"); err != nil {
		return err
	}

	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar documento: %w", err)
	}
	if doc == nil {
		return fmt.Errorf("documento não encontrado")
	}

	baseDir := getBaseDir()
	fullPath := filepath.Join(baseDir, doc.Caminho)

	// Apagar do disco
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao remover arquivo do disco: %w", err)
	}

	// Apagar do banco
	if err := s.docRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao remover documento do banco: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Documento deletado id=%d caminho=%s", id, doc.Caminho),
	})

	return nil
}
