package service

import (
	"AutoGRH/pkg/entity"
	"context"
	"fmt"
	"time"
)

type ValeRepository interface {
	Create(v *entity.Vale) error
	GetByID(id int64) (*entity.Vale, error)
	GetValesByFuncionarioID(funcionarioID int64) ([]entity.Vale, error)
	Update(v *entity.Vale) error
	SoftDelete(id int64) error
	Delete(id int64) error
	ListPendentes() ([]entity.Vale, error)
	ListAprovadosNaoPagos() ([]entity.Vale, error)
}

type ValeService struct {
	repo    ValeRepository
	auth    *AuthService
	logRepo LogRepository
}

func NewValeService(repo ValeRepository, auth *AuthService, logRepo LogRepository) *ValeService {
	return &ValeService{
		repo:    repo,
		auth:    auth,
		logRepo: logRepo,
	}
}

func (s *ValeService) CriarVale(ctx context.Context, claims Claims, funcionarioID int64, valor float64, data time.Time) (*entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:create"); err != nil {
		return nil, err
	}

	v := &entity.Vale{
		FuncionarioID: funcionarioID,
		Valor:         valor,
		Data:          data,
		Aprovado:      false,
		Pago:          false,
		Ativo:         true,
	}

	if err := s.repo.Create(v); err != nil {
		return nil, fmt.Errorf("erro ao criar vale: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Criou vale ID=%d", v.ID),
	})

	return v, nil
}

func (s *ValeService) GetVale(ctx context.Context, claims Claims, id int64) (*entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:read"); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *ValeService) ListarVales(ctx context.Context, claims Claims) ([]entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:list"); err != nil {
		return nil, err
	}
	pendentes, err := s.repo.ListPendentes()
	if err != nil {
		return nil, err
	}
	aprovadosNaoPagos, err := s.repo.ListAprovadosNaoPagos()
	if err != nil {
		return nil, err
	}
	return append(pendentes, aprovadosNaoPagos...), nil
}

func (s *ValeService) ListarValesFuncionario(ctx context.Context, claims Claims, funcionarioID int64) ([]entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:list"); err != nil {
		return nil, err
	}
	return s.repo.GetValesByFuncionarioID(funcionarioID)
}

func (s *ValeService) ListarValesPendentes(ctx context.Context, claims Claims) ([]entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:list"); err != nil {
		return nil, err
	}
	return s.repo.ListPendentes()
}

func (s *ValeService) ListarValesAprovadosNaoPagos(ctx context.Context, claims Claims) ([]entity.Vale, error) {
	if err := s.auth.Authorize(ctx, claims, "vale:list"); err != nil {
		return nil, err
	}
	return s.repo.ListAprovadosNaoPagos()
}

func (s *ValeService) AtualizarVale(ctx context.Context, claims Claims, v *entity.Vale) error {
	if err := s.auth.Authorize(ctx, claims, "vale:update"); err != nil {
		return err
	}
	if err := s.repo.Update(v); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Atualizou vale ID=%d", v.ID),
	})
	return nil
}

func (s *ValeService) SoftDeleteVale(ctx context.Context, claims Claims, id int64) error {
	if err := s.auth.Authorize(ctx, claims, "vale:delete"); err != nil {
		return err
	}
	if err := s.repo.SoftDelete(id); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Soft delete vale ID=%d", id),
	})
	return nil
}

func (s *ValeService) AprovarVale(ctx context.Context, claims Claims, id int64) error {
	if err := s.auth.Authorize(ctx, claims, "vale:update"); err != nil {
		return err
	}

	vale, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if vale == nil {
		return fmt.Errorf("vale %d não encontrado", id)
	}
	vale.Aprovado = true
	if err := s.repo.Update(vale); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Aprovou vale ID=%d", id),
	})
	return nil
}

func (s *ValeService) MarcarValeComoPago(ctx context.Context, claims Claims, id int64) error {
	if err := s.auth.Authorize(ctx, claims, "vale:update"); err != nil {
		return err
	}

	vale, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if vale == nil {
		return fmt.Errorf("vale %d não encontrado", id)
	}
	vale.Pago = true
	if err := s.repo.Update(vale); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Marcou como pago vale ID=%d", id),
	})
	return nil
}

func (s *ValeService) DeleteVale(ctx context.Context, claims Claims, id int64) error {
	if err := s.auth.Authorize(ctx, claims, "vale:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5,
		UsuarioID: &claims.UserID,
		Quando:    s.auth.clock(),
		Detalhe:   fmt.Sprintf("Excluiu vale ID=%d", id),
	})
	return nil
}
