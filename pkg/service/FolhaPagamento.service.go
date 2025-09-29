package service

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"context"
	"fmt"
	"time"
)

type FolhaPagamentoRepository interface {
	Create(f *entity.FolhaPagamentos) error
	GetByID(id int64) (*entity.FolhaPagamentos, error)
	GetByMesAnoTipo(mes, ano int, tipo string) (*entity.FolhaPagamentos, error)
	Update(f *entity.FolhaPagamentos) error
	Delete(id int64) error
	List() ([]entity.FolhaPagamentos, error)
	MarcarComoPaga(id int64) error
}

type FolhaPagamentoService struct {
	repo        FolhaPagamentoRepository
	authService *AuthService
	logRepo     LogRepository
}

func NewFolhaPagamentoService(repo FolhaPagamentoRepository, auth *AuthService, logRepo LogRepository) *FolhaPagamentoService {
	return &FolhaPagamentoService{
		repo:        repo,
		authService: auth,
		logRepo:     logRepo,
	}
}

func (s *FolhaPagamentoService) CriarFolhaSalario(ctx context.Context, claims Claims, mes, ano int) (*entity.FolhaPagamentos, error) {
	if err := s.authService.Authorize(ctx, claims, "folha:create"); err != nil {
		return nil, err
	}

	folha := &entity.FolhaPagamentos{
		Mes:         mes,
		Ano:         ano,
		Tipo:        "SALARIO",
		DataGeracao: time.Now(),
		Pago:        false,
	}

	if err := s.repo.Create(folha); err != nil {
		return nil, fmt.Errorf("erro ao criar folha de salário: %w", err)
	}

	funcionarios, err := repository.ListFuncionariosAtivos()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar funcionários: %w", err)
	}

	var total float64

	for _, f := range funcionarios {
		salarioReal, err := repository.GetSalarioRealAtual(f.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar salário real: %w", err)
		}
		if salarioReal == nil {
			continue
		}

		faltas, err := repository.GetTotalFaltasByFuncionarioMesAno(f.ID, mes, ano)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar faltas: %w", err)
		}

		vales, err := repository.GetValesByFuncionarioMesAno(f.ID, mes, ano)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter vales: %w", err)
		}
		var totalVales float64
		for _, v := range vales {
			if v.Aprovado && v.Ativo && v.Pago {
				totalVales += v.Valor
			}
		}

		salarioBase := salarioReal.Valor
		descontoFaltas := (salarioBase / 30) * float64(faltas)
		valorFinal := salarioBase - descontoFaltas - totalVales

		pag := entity.NewPagamento(f.ID, folha.ID, valorFinal)
		if err := repository.CreatePagamento(pag); err != nil {
			return nil, fmt.Errorf("erro ao criar pagamento: %w", err)
		}

		total += valorFinal
	}

	folha.ValorTotal = total
	if err := s.repo.Update(folha); err != nil {
		return nil, fmt.Errorf("erro ao atualizar total da folha: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Criou folha de salário ID=%d", folha.ID),
	})

	return folha, nil
}

func (s *FolhaPagamentoService) CriarFolhaVale(ctx context.Context, claims Claims, mes, ano int) (*entity.FolhaPagamentos, error) {
	if err := s.authService.Authorize(ctx, claims, "folha:create"); err != nil {
		return nil, err
	}

	folha := &entity.FolhaPagamentos{
		Mes:         mes,
		Ano:         ano,
		Tipo:        "VALE",
		DataGeracao: time.Now(),
		Pago:        false,
	}

	if err := s.repo.Create(folha); err != nil {
		return nil, fmt.Errorf("erro ao criar folha de vale: %w", err)
	}

	vales, err := repository.ListValesAprovadosNaoPagos()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar vales: %w", err)
	}

	var total float64
	for _, v := range vales {
		pag := entity.NewPagamento(v.FuncionarioID, folha.ID, v.Valor)
		if err := repository.CreatePagamento(pag); err != nil {
			return nil, fmt.Errorf("erro ao criar pagamento: %w", err)
		}
		total += v.Valor
	}

	folha.ValorTotal = total
	if err := s.repo.Update(folha); err != nil {
		return nil, fmt.Errorf("erro ao atualizar total da folha: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Criou folha de vale ID=%d", folha.ID),
	})

	return folha, nil
}

// dentro de FolhaPagamento.service.go

func (s *FolhaPagamentoService) RecalcularFolha(ctx context.Context, claims Claims, folhaID int64) error {
	if err := s.authService.Authorize(ctx, claims, "folha:update"); err != nil {
		return err
	}

	folha, err := s.repo.GetByID(folhaID)
	if err != nil {
		return err
	}
	if folha == nil {
		return fmt.Errorf("folha %d não encontrada", folhaID)
	}

	// 🔹 Agora não limpamos mais os pagamentos!
	switch folha.Tipo {
	case "SALARIO":
		if err := s.rebuildPagamentosSalario(ctx, claims, folha); err != nil {
			return err
		}
	case "VALE":
		// mantém a lógica atual de VALE (puxa todos os aprovados/não pagos)
		return s.RecalcularFolhaVale(ctx, claims, folhaID)
	default:
		return fmt.Errorf("tipo de folha desconhecido: %s", folha.Tipo)
	}

	// log
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Recalculou folha ID=%d", folha.ID),
	})
	return nil
}

func (s *FolhaPagamentoService) rebuildPagamentosSalario(
	ctx context.Context,
	claims Claims,
	folha *entity.FolhaPagamentos,
) error {
	funcionarios, err := repository.ListFuncionariosAtivos()
	if err != nil {
		return fmt.Errorf("erro ao listar funcionários: %w", err)
	}

	// Buscar pagamentos já existentes da folha
	existentes, err := repository.GetPagamentosByFolhaID(folha.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar pagamentos existentes: %w", err)
	}
	mapPag := make(map[int64]*entity.Pagamento)
	for i := range existentes {
		mapPag[existentes[i].FuncionarioID] = &existentes[i]
	}

	var total float64
	for _, f := range funcionarios {
		salarioReal, err := repository.GetSalarioRealAtual(f.ID)
		if err != nil {
			return fmt.Errorf("erro ao buscar salário real: %w", err)
		}
		if salarioReal == nil {
			continue
		}

		faltas, err := repository.GetTotalFaltasByFuncionarioMesAno(f.ID, folha.Mes, folha.Ano)
		if err != nil {
			return fmt.Errorf("erro ao buscar faltas: %w", err)
		}

		vales, err := repository.GetValesByFuncionarioMesAno(f.ID, folha.Mes, folha.Ano)
		if err != nil {
			return fmt.Errorf("erro ao obter vales: %w", err)
		}
		var totalVales float64
		for _, v := range vales {
			if v.Aprovado && v.Ativo && v.Pago {
				totalVales += v.Valor
			}
		}

		// cálculo automático
		salarioBase := salarioReal.Valor
		descontoFaltas := (salarioBase / 30) * float64(faltas)

		if pag, ok := mapPag[f.ID]; ok {
			pag.SalarioBase = salarioBase
			pag.DescontoVales = totalVales
			pag.RecalcularValorFinal(descontoFaltas)

			if err := repository.UpdatePagamento(pag); err != nil {
				return fmt.Errorf("erro ao atualizar pagamento: %w", err)
			}
			total += pag.ValorFinal
		} else {
			p := entity.NewPagamento(f.ID, folha.ID, salarioBase)
			p.DescontoVales = totalVales
			p.RecalcularValorFinal(descontoFaltas)

			if err := repository.CreatePagamento(p); err != nil {
				return fmt.Errorf("erro ao criar pagamento: %w", err)
			}
			total += p.ValorFinal
		}
	}

	// Atualizar total da folha
	folha.ValorTotal = total
	if err := s.repo.Update(folha); err != nil {
		return fmt.Errorf("erro ao atualizar total da folha: %w", err)
	}
	return nil
}

func (s *FolhaPagamentoService) FecharFolha(ctx context.Context, claims Claims, folhaID int64) error {
	if err := s.authService.Authorize(ctx, claims, "folha:update"); err != nil {
		return err
	}

	folha, err := s.repo.GetByID(folhaID)
	if err != nil {
		return err
	}
	if folha == nil {
		return fmt.Errorf("folha %d não encontrada", folhaID)
	}

	if folha.Tipo == "VALE" {
		if err := repository.MarcarTodosValesComoPagos(); err != nil {
			return fmt.Errorf("erro ao marcar vales como pagos: %w", err)
		}
	}
	if err := repository.MarcarPagamentosDaFolhaComoPagos(folha.ID); err != nil {
		return fmt.Errorf("erro ao marcar pagamentos da folha como pagos: %w", err)
	}

	if err := s.repo.MarcarComoPaga(folha.ID); err != nil {
		return fmt.Errorf("erro ao marcar folha como paga: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Fechou folha ID=%d", folha.ID),
	})

	return nil
}

func (s *FolhaPagamentoService) ListarFolhas(ctx context.Context, claims Claims) ([]entity.FolhaPagamentos, error) {
	if err := s.authService.Authorize(ctx, claims, "folha:list"); err != nil {
		return nil, err
	}
	return s.repo.List()
}

func (s *FolhaPagamentoService) BuscarFolha(ctx context.Context, claims Claims, folhaID int64) (*entity.FolhaPagamentos, error) {
	if err := s.authService.Authorize(ctx, claims, "folha:read"); err != nil {
		return nil, err
	}
	return s.repo.GetByID(folhaID)
}

func (s *FolhaPagamentoService) BuscarFolhaPorMesAnoTipo(ctx context.Context, claims Claims, mes, ano int, tipo string) (*entity.FolhaPagamentos, error) {
	if err := s.authService.Authorize(ctx, claims, "folha:read"); err != nil {
		return nil, err
	}
	return s.repo.GetByMesAnoTipo(mes, ano, tipo)
}

func (s *FolhaPagamentoService) ExcluirFolha(ctx context.Context, claims Claims, folhaID int64) error {
	if err := s.authService.Authorize(ctx, claims, "folha:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(folhaID); err != nil {
		return err
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  5,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Excluiu folha ID=%d", folhaID),
	})
	return nil
}

// RecalcularFolhaVale refaz os pagamentos de uma folha do tipo VALE
func (s *FolhaPagamentoService) RecalcularFolhaVale(ctx context.Context, claims Claims, folhaID int64) error {
	if err := s.authService.Authorize(ctx, claims, "folha:update"); err != nil {
		return err
	}

	// Carrega a folha e valida tipo
	folha, err := s.repo.GetByID(folhaID)
	if err != nil {
		return fmt.Errorf("erro ao buscar folha: %w", err)
	}
	if folha == nil {
		return fmt.Errorf("folha %d não encontrada", folhaID)
	}
	if folha.Tipo != "VALE" {
		return fmt.Errorf("folha %d não é do tipo VALE", folhaID)
	}

	// Remove pagamentos antigos da folha
	if err := repository.DeletePagamentosByFolhaID(folhaID); err != nil {
		return fmt.Errorf("erro ao limpar pagamentos antigos: %w", err)
	}

	// Recria pagamentos a partir dos vales aprovados e não pagos (sem filtro de data, por design)
	vales, err := repository.ListValesAprovadosNaoPagos()
	if err != nil {
		return fmt.Errorf("erro ao listar vales aprovados e não pagos: %w", err)
	}

	var total float64
	for _, v := range vales {
		p := entity.NewPagamento(v.FuncionarioID, folha.ID, v.Valor)
		if err := repository.CreatePagamento(p); err != nil {
			return fmt.Errorf("erro ao criar pagamento do vale (valeID=%d): %w", v.ID, err)
		}
		total += v.Valor
	}

	// Atualiza total da folha
	folha.ValorTotal = total
	if err := s.repo.Update(folha); err != nil {
		return fmt.Errorf("erro ao atualizar total da folha: %w", err)
	}

	// Log no padrão dos demais services
	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    s.authService.clock(),
		Detalhe:   fmt.Sprintf("Recalculou folha de vale ID=%d (total=%.2f)", folha.ID, total),
	})

	return nil
}
