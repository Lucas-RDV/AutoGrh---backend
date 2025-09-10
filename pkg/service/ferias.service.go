package service

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"context"
	"fmt"
	"time"
)

// FeriasRepository define as operações de acesso ao banco
type FeriasRepository interface {
	Create(ctx context.Context, f *entity.Ferias) error
	GetFeriasByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Ferias, error)
	GetByID(ctx context.Context, id int64) (*entity.Ferias, error)
	Update(ctx context.Context, f *entity.Ferias) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*entity.Ferias, error)
}

type SaldoFeriasDTO struct {
	DiasRestantes int     `json:"dias_restantes"`
	Valor         float64 `json:"valor"`
	Terco         float64 `json:"terco"`
	Total         float64 `json:"total"`
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

// CRUD

func (s *FeriasService) CriarFerias(ctx context.Context, claims Claims, funcionarioID int64, dias int, valor float64, inicio time.Time) (*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:create"); err != nil {
		return nil, err
	}

	f := entity.NewFerias(funcionarioID, inicio, dias)
	f.Valor = valor
	f.Terco = valor / 3.0
	f.TercoPago = false
	f.Vencido = false

	if err := s.repo.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("erro ao criar férias: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3, // CRIAR
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias criadas funcionarioID=%d dias=%d inicio=%s", funcionarioID, dias, inicio.Format("2006-01-02")),
	})

	return f, nil
}

func (s *FeriasService) GetFeriasByID(ctx context.Context, claims Claims, id int64) (*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:read"); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *FeriasService) GetFeriasByFuncionarioID(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:list"); err != nil {
		return nil, err
	}
	return s.repo.GetFeriasByFuncionarioID(ctx, funcionarioID)
}

func (s *FeriasService) ListFerias(ctx context.Context, claims Claims) ([]*entity.Ferias, error) {
	if err := s.authService.Authorize(ctx, claims, "ferias:list"); err != nil {
		return nil, err
	}
	return s.repo.List(ctx)
}

func (s *FeriasService) AtualizarFerias(ctx context.Context, claims Claims, f *entity.Ferias) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:update"); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, f); err != nil {
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

func (s *FeriasService) DeletarFerias(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
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

// Regras de negócio

// MarcarComoVencidas define que o período expirou
func (s *FeriasService) MarcarComoVencidas(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:update"); err != nil {
		return err
	}
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar férias: %w", err)
	}
	if f == nil {
		return fmt.Errorf("férias não encontradas")
	}
	f.Vencido = true
	if err := s.repo.Update(ctx, f); err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Férias vencidas id=%d", id),
	})
	return nil
}

// MarcarTercoComoPago define que o terço já foi pago
func (s *FeriasService) MarcarTercoComoPago(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "ferias:update"); err != nil {
		return err
	}
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar férias: %w", err)
	}
	if f == nil {
		return fmt.Errorf("férias não encontradas")
	}
	f.TercoPago = true
	if err := s.repo.Update(ctx, f); err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Terço pago de férias id=%d", id),
	})
	return nil
}

// CalcularSaldo retorna os dias e valores restantes das férias
func (s *FeriasService) CalcularSaldo(ctx context.Context, claims Claims, f *entity.Ferias) (*SaldoFeriasDTO, error) {
	//  Verifica permissão
	if err := s.authService.Authorize(ctx, claims, "ferias:read"); err != nil {
		return nil, err
	}

	//  Buscar salário real atual do funcionário
	salarioReal, err := repository.GetSalarioRealAtual(f.FuncionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salário real atual: %w", err)
	}
	if salarioReal == nil {
		return nil, fmt.Errorf("nenhum salário real encontrado para funcionarioID=%d", f.FuncionarioID)
	}

	//  Calcular saldo
	diasRestantes := f.DiasRestantes()
	valorDias := (salarioReal.Valor / 30.0) * float64(diasRestantes)
	terco := f.Terco

	var total float64
	if !f.TercoPago {
		total = valorDias + terco
	} else {
		total = valorDias
	}

	dto := &SaldoFeriasDTO{
		DiasRestantes: diasRestantes,
		Valor:         valorDias,
		Terco:         terco,
		Total:         total,
	}

	return dto, nil
}
