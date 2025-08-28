package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

// SalarioRepository define as operações de acesso a dados para salários registrados
type SalarioRepository interface {
	Create(s *entity.Salario) error
	GetSalariosByFuncionarioID(funcionarioID int64) ([]*entity.Salario, error)
	Update(s *entity.Salario) error
	Delete(id int64) error
}

// SalarioService encapsula a lógica de negócio para salários registrados
type SalarioService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        SalarioRepository
}

func NewSalarioService(auth *AuthService, logRepo LogRepository, repo SalarioRepository) *SalarioService {
	return &SalarioService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CriarSalario insere um novo salário registrado (não encerra automaticamente o anterior)
func (s *SalarioService) CriarSalario(ctx context.Context, claims Claims, funcionarioID int64, valor float64) (*entity.Salario, error) {
	if err := s.authService.Authorize(ctx, claims, "salario:create"); err != nil {
		return nil, err
	}

	novo := &entity.Salario{
		FuncionarioID: funcionarioID,
		Inicio:        time.Now(),
		Valor:         valor,
		Fim:           nil,
	}
	if err := s.repo.Create(novo); err != nil {
		return nil, fmt.Errorf("erro ao criar salário: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Salário registrado criado funcionarioID=%d valor=%.2f", funcionarioID, valor),
	})

	return novo, nil
}

// ListSalarios retorna todos os salários registrados de um funcionário
func (s *SalarioService) ListSalarios(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Salario, error) {
	if err := s.authService.Authorize(ctx, claims, "salario:list"); err != nil {
		return nil, err
	}
	return s.repo.GetSalariosByFuncionarioID(funcionarioID)
}

// AtualizarSalario altera dados de um salário registrado existente
func (s *SalarioService) AtualizarSalario(ctx context.Context, claims Claims, sEntity *entity.Salario) error {
	if err := s.authService.Authorize(ctx, claims, "salario:update"); err != nil {
		return err
	}
	if err := s.repo.Update(sEntity); err != nil {
		return fmt.Errorf("erro ao atualizar salário: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Salário registrado atualizado id=%d valor=%.2f", sEntity.ID, sEntity.Valor),
	})
	return nil
}

// DeletarSalario remove um salário registrado
func (s *SalarioService) DeletarSalario(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "salario:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar salário: %w", err)
	}

	// Log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5, // DELETAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Salário registrado deletado id=%d", id),
	})
	return nil
}
