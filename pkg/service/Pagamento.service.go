package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
)

type PagamentoRepository interface {
	GetPagamentoByID(id int64) (*entity.Pagamento, error)
	Update(p *entity.Pagamento) error
	ListPagamentosByFuncionarioID(funcionarioID int64) ([]entity.Pagamento, error)
}

type PagamentoService struct {
	repo    PagamentoRepository
	auth    *AuthService
	logRepo LogRepository
}

func NewPagamentoService(repo PagamentoRepository, auth *AuthService, logRepo LogRepository) *PagamentoService {
	return &PagamentoService{
		repo:    repo,
		auth:    auth,
		logRepo: logRepo,
	}
}

// BuscarPagamento retorna um pagamento específico
func (s *PagamentoService) BuscarPagamento(ctx context.Context, claims Claims, pagamentoID int64) (*entity.Pagamento, error) {
	if err := s.auth.Authorize(ctx, claims, "pagamento:read"); err != nil {
		return nil, err
	}

	p, err := s.repo.GetPagamentoByID(pagamentoID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("pagamento %d não encontrado", pagamentoID)
	}

	return p, nil
}

// AtualizarPagamento permite ajustes manuais em adicional, descontoINSS e salarioFamilia
func (s *PagamentoService) AtualizarPagamento(
	ctx context.Context, claims Claims, pagamentoID int64,
	adicional, descontoINSS, salarioFamilia float64,
) error {
	if err := s.auth.Authorize(ctx, claims, "pagamento:update"); err != nil {
		return err
	}

	p, err := s.repo.GetPagamentoByID(pagamentoID)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("pagamento %d não encontrado", pagamentoID)
	}

	// aplicar ajustes manuais
	p.Adicional = adicional
	p.DescontoINSS = descontoINSS
	p.SalarioFamilia = salarioFamilia

	// 🔹 manter desconto de vales já armazenado
	descontoFaltas := 0.0 // só recalculado no rebuild da folha
	p.RecalcularValorFinal(descontoFaltas)

	if err := s.repo.Update(p); err != nil {
		return fmt.Errorf("erro ao atualizar pagamento: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Pagamento %d atualizado manualmente", p.ID),
	})

	return nil
}

// ListarPagamentosFuncionario retorna todos os pagamentos de um funcionário
func (s *PagamentoService) ListarPagamentosFuncionario(ctx context.Context, claims Claims, funcionarioID int64) ([]entity.Pagamento, error) {
	if err := s.auth.Authorize(ctx, claims, "pagamento:list"); err != nil {
		return nil, err
	}
	return s.repo.ListPagamentosByFuncionarioID(funcionarioID)
}

// MarcarPagamentoComoPago marca um pagamento individual como pago (somente Admin)
func (s *PagamentoService) MarcarPagamentoComoPago(ctx context.Context, claims Claims, pagamentoID int64) error {
	if err := s.auth.Authorize(ctx, claims, "pagamento:update"); err != nil {
		return err
	}

	p, err := s.repo.GetPagamentoByID(pagamentoID)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("pagamento %d não encontrado", pagamentoID)
	}

	p.Pago = true
	if err := s.repo.Update(p); err != nil {
		return fmt.Errorf("erro ao marcar pagamento como pago: %w", err)
	}

	// log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Pagamento %d marcado como pago", p.ID),
	})

	return nil
}
