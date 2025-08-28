package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

// FeriasRepository define as operações de acesso ao banco
type FeriasRepository interface {
	Create(f *entity.Ferias) error
	GetFeriasByFuncionarioID(funcionarioID int64) ([]*entity.Ferias, error)
	GetByID(id int64) (*entity.Ferias, error)
	Update(f *entity.Ferias) error
	Delete(id int64) error
	List() ([]*entity.Ferias, error)
}

type FeriasService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        FeriasRepository
}

func NewFeriasService(auth *AuthService, logRepo LogRepository, repo FeriasRepository) *FeriasService {
	return &FeriasService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CriarFerias insere um novo direito de férias
func (s *FeriasService) CriarFerias(ctx context.Context, claims Claims, funcionarioID int64, dias int, valor float64, inicio time.Time) (*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:create"); err != nil {
		return nil, err
	}

	f := entity.NewFerias(inicio)
	f.FuncionarioID = funcionarioID
	f.Dias = dias
	f.Valor = valor
	f.Vencido = false

	if err := s.repo.Create(f); err != nil {
		return nil, fmt.Errorf("erro ao criar férias: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias criadas funcionarioID=%d dias=%d inicio=%s", funcionarioID, dias, inicio.Format("2006-01-02")),
	})

	return f, nil
}

// GetFeriasByFuncionarioID retorna as férias de um funcionário
func (s *FeriasService) GetFeriasByFuncionarioID(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:list"); err != nil {
		return nil, err
	}
	return s.repo.GetFeriasByFuncionarioID(funcionarioID)
}

// ListFerias retorna todas as férias cadastradas
func (s *FeriasService) ListFerias(ctx context.Context, claims Claims) ([]*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:list"); err != nil {
		return nil, err
	}
	return s.repo.List()
}

// MarcarComoVencidas atualiza o campo Vencido das férias
func (s *FeriasService) MarcarComoVencidas(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:update"); err != nil {
		return err
	}

	f, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar férias: %w", err)
	}
	if f == nil {
		return fmt.Errorf("férias não encontradas")
	}

	f.Vencido = true
	if err := s.repo.Update(f); err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias marcadas como vencidas id=%d", id),
	})
	return nil
}

// AtualizarFerias altera os dados de férias
func (s *FeriasService) AtualizarFerias(ctx context.Context, claims Claims, f *entity.Ferias) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:update"); err != nil {
		return err
	}
	if err := s.repo.Update(f); err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias atualizadas id=%d", f.ID),
	})
	return nil
}

// DeletarFerias remove um registro de férias
func (s *FeriasService) DeletarFerias(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar férias: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias deletadas id=%d", id),
	})
	return nil
}
