package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

// DescansoRepository define as operações necessárias
type DescansoRepository interface {
	Create(d *entity.Descanso) error
	GetDescansoByID(id int64) (*entity.Descanso, error)
	Update(d *entity.Descanso) error
	Delete(id int64) error
	GetDescansosByFeriasID(feriasID int64) ([]*entity.Descanso, error)
	GetDescansosByFuncionarioID(funcionarioID int64) ([]*entity.Descanso, error)
	GetDescansosAprovados() ([]*entity.Descanso, error)
	GetDescansosPendentes() ([]*entity.Descanso, error)
}

type DescansoService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        DescansoRepository
}

// Construtor
func NewDescansoService(auth *AuthService, logRepo LogRepository, repo DescansoRepository) *DescansoService {
	return &DescansoService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// Criar um novo descanso
func (s *DescansoService) CreateDescanso(ctx context.Context, claims Claims, d *entity.Descanso) error {
	if err := s.authService.Authorize(ctx, claims, "descanso:create"); err != nil {
		return err
	}

	if d.Inicio.After(d.Fim) {
		return fmt.Errorf("data inicial não pode ser depois da final")
	}

	if err := s.repo.Create(d); err != nil {
		return fmt.Errorf("erro ao criar descanso: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso criado ID=%d FeriasID=%d", d.ID, d.FeriasID),
	})

	return nil
}

// Aprovar descanso (admin)
func (s *DescansoService) AprovarDescanso(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "descanso:update"); err != nil {
		return err
	}

	descanso, err := s.repo.GetDescansoByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar descanso: %w", err)
	}
	if descanso == nil {
		return fmt.Errorf("descanso não encontrado")
	}

	descanso.Aprovado = true
	if err := s.repo.Update(descanso); err != nil {
		return fmt.Errorf("erro ao aprovar descanso: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso aprovado ID=%d", id),
	})

	return nil
}

// Marcar descanso como pago
func (s *DescansoService) MarcarComoPago(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "descanso:update"); err != nil {
		return err
	}

	descanso, err := s.repo.GetDescansoByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar descanso: %w", err)
	}
	if descanso == nil {
		return fmt.Errorf("descanso não encontrado")
	}

	descanso.Pago = true
	if err := s.repo.Update(descanso); err != nil {
		return fmt.Errorf("erro ao marcar descanso como pago: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso pago ID=%d", id),
	})

	return nil
}

// Listar descansos por férias
func (s *DescansoService) ListarPorFerias(ctx context.Context, claims Claims, feriasID int64) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, "descanso:list"); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosByFeriasID(feriasID)
}

// Listar descansos por funcionário
func (s *DescansoService) ListarPorFuncionario(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, "descanso:list"); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosByFuncionarioID(funcionarioID)
}

// Listar descansos aprovados
func (s *DescansoService) ListarAprovados(ctx context.Context, claims Claims) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, "descanso:list"); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosAprovados()
}

// Listar descansos pendentes
func (s *DescansoService) ListarPendentes(ctx context.Context, claims Claims) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, "descanso:list"); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosPendentes()
}

// Excluir descanso
func (s *DescansoService) DeleteDescanso(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "descanso:delete"); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar descanso: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso deletado ID=%d", id),
	})

	return nil
}
