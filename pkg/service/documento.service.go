package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
)

// DocumentoRepository interface usada pelo service
type DocumentoRepository interface {
	Create(ctx context.Context, d *entity.Documento) error
	GetByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error)
	List(ctx context.Context) ([]*entity.Documento, error)
	Delete(ctx context.Context, id int64) error
}

type DocumentoService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        DocumentoRepository
}

func NewDocumentoService(auth *AuthService, logRepo LogRepository, repo DocumentoRepository) *DocumentoService {
	return &DocumentoService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CreateDocumento insere um novo documento para um funcionário
func (s *DocumentoService) CreateDocumento(ctx context.Context, claims Claims, d *entity.Documento) error {
	if d.FuncionarioID <= 0 {
		return fmt.Errorf("funcionarioID inválido")
	}
	if len(d.Doc) == 0 {
		return fmt.Errorf("conteúdo do documento não pode ser vazio")
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Upload documento ID=%d para funcionario ID=%d", d.ID, d.FuncionarioID),
	})

	return nil
}

// GetDocumentosByFuncionarioID lista documentos de um funcionário
func (s *DocumentoService) GetDocumentosByFuncionarioID(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Documento, error) {
	if funcionarioID <= 0 {
		return nil, fmt.Errorf("funcionarioID inválido")
	}
	return s.repo.GetByFuncionarioID(ctx, funcionarioID)
}

// ListDocumentos retorna todos os documentos cadastrados
func (s *DocumentoService) ListDocumentos(ctx context.Context, claims Claims) ([]*entity.Documento, error) {
	return s.repo.List(ctx)
}

// DeleteDocumento remove um documento (somente admin)
func (s *DocumentoService) DeleteDocumento(ctx context.Context, claims Claims, id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID inválido")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Deletou documento ID=%d", id),
	})

	return nil
}
