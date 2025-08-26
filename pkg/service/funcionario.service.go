package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"strings"
)

// FuncionarioRepository interface usada pelo service
type FuncionarioRepository interface {
	Create(ctx context.Context, f *entity.Funcionario) error
	GetByID(ctx context.Context, id int64) (*entity.Funcionario, error)
	Update(ctx context.Context, f *entity.Funcionario) error
	Delete(ctx context.Context, id int64) error
	ListAtivos(ctx context.Context) ([]*entity.Funcionario, error)
	ListInativos(ctx context.Context) ([]*entity.Funcionario, error)
	ListTodos(ctx context.Context) ([]*entity.Funcionario, error)
}

type FuncionarioService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        FuncionarioRepository
}

func NewFuncionarioService(auth *AuthService, logRepo LogRepository, repo FuncionarioRepository) *FuncionarioService {
	return &FuncionarioService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CreateFuncionario cria um novo funcionário
func (s *FuncionarioService) CreateFuncionario(ctx context.Context, claims Claims, f *entity.Funcionario) error {
	f.Cargo = strings.TrimSpace(f.Cargo)

	if f.PessoaID <= 0 {
		return fmt.Errorf("PessoaID inválido ou não informado")
	}
	if f.Cargo == "" {
		return fmt.Errorf("cargo não pode ser vazio")
	}
	if f.SalarioInicial < 0 {
		return fmt.Errorf("salário inicial não pode ser negativo")
	}
	if f.FeriasDisponiveis < 0 {
		return fmt.Errorf("dias de férias disponíveis não podem ser negativos")
	}
	if f.Admissao.IsZero() {
		return fmt.Errorf("data de admissão inválida")
	}
	if f.Nascimento.IsZero() {
		return fmt.Errorf("data de nascimento inválida")
	}

	if err := s.repo.Create(ctx, f); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Criou funcionário ID=%d", f.ID),
	})

	return nil
}

// GetFuncionarioByID busca um funcionário pelo ID
func (s *FuncionarioService) GetFuncionarioByID(ctx context.Context, claims Claims, id int64) (*entity.Funcionario, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateFuncionario atualiza os dados de um funcionário
func (s *FuncionarioService) UpdateFuncionario(ctx context.Context, claims Claims, f *entity.Funcionario) error {
	f.Cargo = strings.TrimSpace(f.Cargo)

	if f.ID <= 0 {
		return fmt.Errorf("ID do funcionário inválido")
	}
	if f.PessoaID <= 0 {
		return fmt.Errorf("PessoaID inválido")
	}
	if f.Cargo == "" {
		return fmt.Errorf("cargo não pode ser vazio")
	}
	if f.SalarioInicial < 0 {
		return fmt.Errorf("salário inicial não pode ser negativo")
	}
	if f.FeriasDisponiveis < 0 {
		return fmt.Errorf("dias de férias disponíveis não podem ser negativos")
	}
	if f.Admissao.IsZero() {
		return fmt.Errorf("data de admissão inválida")
	}
	if f.Nascimento.IsZero() {
		return fmt.Errorf("data de nascimento inválido")
	}

	if err := s.repo.Update(ctx, f); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Atualizou funcionário ID=%d", f.ID),
	})

	return nil
}

// DeleteFuncionario remove um funcionário
func (s *FuncionarioService) DeleteFuncionario(ctx context.Context, claims Claims, id int64) error {
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
		Detalhe:   fmt.Sprintf("Deletou funcionário ID=%d", id),
	})

	return nil
}

// ListFuncionariosAtivos lista todos os funcionários ativos
func (s *FuncionarioService) ListFuncionariosAtivos(ctx context.Context, claims Claims) ([]*entity.Funcionario, error) {
	return s.repo.ListAtivos(ctx)
}

// ListFuncionariosInativos lista todos os funcionários inativos
func (s *FuncionarioService) ListFuncionariosInativos(ctx context.Context, claims Claims) ([]*entity.Funcionario, error) {
	return s.repo.ListInativos(ctx)
}

// ListTodosFuncionarios lista todos os funcionários, ativos e inativos
func (s *FuncionarioService) ListTodosFuncionarios(ctx context.Context, claims Claims) ([]*entity.Funcionario, error) {
	return s.repo.ListTodos(ctx)
}
