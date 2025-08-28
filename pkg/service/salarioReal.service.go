package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

// SalarioRealRepository define as operações necessárias para salários reais
type SalarioRealRepository interface {
	Create(s *entity.SalarioReal) error
	GetByFuncionarioID(funcionarioID int64) ([]*entity.SalarioReal, error)
	GetAtual(funcionarioID int64) (*entity.SalarioReal, error)
	Update(s *entity.SalarioReal) error
	Delete(id int64) error
}

// SalarioRealService encapsula a lógica de negócio para salários reais
type SalarioRealService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        SalarioRealRepository
}

func NewSalarioRealService(auth *AuthService, logRepo LogRepository, repo SalarioRealRepository) *SalarioRealService {
	return &SalarioRealService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CriarSalarioReal encerra o salário atual (se houver) e insere um novo
func (s *SalarioRealService) CriarSalarioReal(ctx context.Context, claims Claims, funcionarioID int64, valor float64) (*entity.SalarioReal, error) {
	if err := s.authService.Authorize(ctx, claims, "salario_real:create"); err != nil {
		return nil, err
	}

	// Encerrar salário atual, se existir
	atual, err := s.repo.GetAtual(funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar salário atual: %w", err)
	}
	if atual != nil {
		now := time.Now()
		atual.Fim = &now
		if err := s.repo.Update(atual); err != nil {
			return nil, fmt.Errorf("erro ao encerrar salário atual: %w", err)
		}
	}

	// Criar novo salário
	novo := &entity.SalarioReal{
		FuncionarioID: funcionarioID,
		Inicio:        time.Now(),
		Valor:         valor,
	}
	if err := s.repo.Create(novo); err != nil {
		return nil, fmt.Errorf("erro ao criar salário real: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Salário real criado funcionarioID=%d valor=%.2f", funcionarioID, valor),
	})

	return novo, nil
}

// GetSalarioRealAtual retorna o salário real atual de um funcionário
func (s *SalarioRealService) GetSalarioRealAtual(ctx context.Context, claims Claims, funcionarioID int64) (*entity.SalarioReal, error) {
	if err := s.authService.Authorize(ctx, claims, "salario_real:list"); err != nil {
		return nil, err
	}
	return s.repo.GetAtual(funcionarioID)
}

// ListSalariosReais retorna o histórico de salários reais de um funcionário
func (s *SalarioRealService) ListSalariosReais(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.SalarioReal, error) {
	if err := s.authService.Authorize(ctx, claims, "salario_real:list"); err != nil {
		return nil, err
	}
	return s.repo.GetByFuncionarioID(funcionarioID)
}

// DeleteSalarioReal exclui um registro (se permitido pela regra do cliente)
func (s *SalarioRealService) DeleteSalarioReal(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "salario_real:delete"); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar salário real: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Salário real deletado id=%d", id),
	})

	return nil
}
