package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"strings"
)

// PessoaRepository interface usada pelo service
// Isso permite injetar diferentes implementações (ex.: MySQL, mock para testes)
type PessoaRepository interface {
	Create(ctx context.Context, p *entity.Pessoa) error
	GetByID(ctx context.Context, id int64) (*entity.Pessoa, error)
	GetByCPF(ctx context.Context, cpf string) (*entity.Pessoa, error)
	Update(ctx context.Context, p *entity.Pessoa) error
	Delete(ctx context.Context, id int64) error
	ExistsByCPF(ctx context.Context, cpf string) (bool, error)
	ExistsByRG(ctx context.Context, rg string) (bool, error)
	SearchByNome(ctx context.Context, nome string) ([]*entity.Pessoa, error)
	ListAll(ctx context.Context) ([]*entity.Pessoa, error)
}

type PessoaService struct {
	authService *AuthService
	logRepo     LogRepository
	repo        PessoaRepository
}

func NewPessoaService(auth *AuthService, logRepo LogRepository, repo PessoaRepository) *PessoaService {
	return &PessoaService{
		authService: auth,
		logRepo:     logRepo,
		repo:        repo,
	}
}

// CreatePessoa cria uma nova pessoa
func (s *PessoaService) CreatePessoa(ctx context.Context, claims Claims, p *entity.Pessoa) error {
	p.Nome = strings.TrimSpace(p.Nome)
	p.CPF = strings.TrimSpace(p.CPF)
	p.RG = strings.TrimSpace(p.RG)

	if p.Nome == "" {
		return fmt.Errorf("nome não pode ser vazio")
	}
	if p.CPF == "" {
		return fmt.Errorf("CPF não pode ser vazio")
	}
	if p.RG == "" {
		return fmt.Errorf("RG não pode ser vazio")
	}

	exists, err := s.repo.ExistsByCPF(ctx, p.CPF)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("já existe uma pessoa com este CPF")
	}

	exists, err = s.repo.ExistsByRG(ctx, p.RG)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("já existe uma pessoa com este RG")
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Criou pessoa ID=%d", p.ID),
	})

	return nil
}

// GetPessoaByID retorna uma pessoa pelo ID
func (s *PessoaService) GetPessoaByID(ctx context.Context, claims Claims, id int64) (*entity.Pessoa, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}
	return s.repo.GetByID(ctx, id)
}

// UpdatePessoa atualiza os dados de uma pessoa
func (s *PessoaService) UpdatePessoa(ctx context.Context, claims Claims, p *entity.Pessoa) error {
	if p.ID <= 0 {
		return fmt.Errorf("ID inválido")
	}
	p.Nome = strings.TrimSpace(p.Nome)
	p.CPF = strings.TrimSpace(p.CPF)
	p.RG = strings.TrimSpace(p.RG)

	if p.Nome == "" || p.CPF == "" || p.RG == "" {
		return fmt.Errorf("nome, CPF e RG não podem ser vazios")
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4, // ATUALIZAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Atualizou pessoa ID=%d", p.ID),
	})

	return nil
}

// DeletePessoa remove uma pessoa
func (s *PessoaService) DeletePessoa(ctx context.Context, claims Claims, id int64) error {
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
		Detalhe:   fmt.Sprintf("Deletou pessoa ID=%d", id),
	})

	return nil
}

// ListPessoas retorna todas as pessoas
func (s *PessoaService) ListPessoas(ctx context.Context, claims Claims) ([]*entity.Pessoa, error) {
	return s.repo.ListAll(ctx)
}
