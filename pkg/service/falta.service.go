package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

// FaltaRepository define as operações de persistência necessárias
type FaltaRepository interface {
	Create(f *entity.Falta) error
	GetFaltaByID(id int64) (*entity.Falta, error)
	Update(f *entity.Falta) error
	Delete(id int64) error
	GetFaltasByFuncionarioID(funcionarioID int64) ([]*entity.Falta, error)
	ListAll() ([]*entity.Falta, error)
}

type FaltaService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        FaltaRepository
}

// Construtor
func NewFaltaService(auth *AuthService, logRepo LogRepository, repo FaltaRepository) *FaltaService {
	return &FaltaService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// Criar nova falta
func (s *FaltaService) CreateFalta(ctx context.Context, claims Claims, f *entity.Falta) error {
	if err := s.authService.Authorize(ctx, claims, "falta:create"); err != nil {
		return err
	}

	if f.Quantidade <= 0 {
		return fmt.Errorf("quantidade de faltas deve ser maior que zero")
	}

	if err := s.repo.Create(f); err != nil {
		return fmt.Errorf("erro ao registrar falta: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Falta criada ID=%d FuncionarioID=%d Qtd=%d", f.ID, f.FuncionarioID, f.Quantidade),
	})

	return nil
}

// Atualizar falta
func (s *FaltaService) UpdateFalta(ctx context.Context, claims Claims, f *entity.Falta) error {
	if err := s.authService.Authorize(ctx, claims, "falta:update"); err != nil {
		return err
	}

	if f.Quantidade <= 0 {
		return fmt.Errorf("quantidade de faltas deve ser maior que zero")
	}

	if err := s.repo.Update(f); err != nil {
		return fmt.Errorf("erro ao atualizar falta: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Falta atualizada ID=%d FuncionarioID=%d Qtd=%d", f.ID, f.FuncionarioID, f.Quantidade),
	})

	return nil
}

// Deletar falta
func (s *FaltaService) DeleteFalta(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "falta:delete"); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar falta: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Falta deletada ID=%d", id),
	})

	return nil
}

// Buscar falta por ID
func (s *FaltaService) GetFaltaByID(ctx context.Context, claims Claims, id int64) (*entity.Falta, error) {
	if err := s.authService.Authorize(ctx, claims, "falta:read"); err != nil {
		return nil, err
	}

	falta, err := s.repo.GetFaltaByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar falta: %w", err)
	}
	if falta == nil {
		return nil, fmt.Errorf("falta não encontrada")
	}
	return falta, nil
}

// Listar todas as faltas de um funcionário
func (s *FaltaService) GetFaltasByFuncionarioID(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Falta, error) {
	if err := s.authService.Authorize(ctx, claims, "falta:list"); err != nil {
		return nil, err
	}
	return s.repo.GetFaltasByFuncionarioID(funcionarioID)
}

// Listar todas as faltas
func (s *FaltaService) ListAllFaltas(ctx context.Context, claims Claims) ([]*entity.Falta, error) {
	if err := s.authService.Authorize(ctx, claims, "falta:list"); err != nil {
		return nil, err
	}
	return s.repo.ListAll()
}
