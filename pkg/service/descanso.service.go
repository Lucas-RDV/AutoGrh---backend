package service

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"context"
	"fmt"
	"time"
)

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

func NewDescansoService(auth *AuthService, logRepo LogRepository, repo DescansoRepository) *DescansoService {
	return &DescansoService{authService: auth, logRepo: logRepo, repo: repo}
}

func diasReferenciaFerias(f *entity.Ferias) int {
	if f == nil {
		return 30
	}
	if f.Terco > 0 && f.Valor > 0 {
		dias := int((10.0 * f.Valor / f.Terco) + 0.5)
		if dias > 0 {
			return dias
		}
	}
	if f.Dias > 0 {
		return f.Dias
	}
	return 30
}

func (s *DescansoService) consumirDiasDescansoUnico(feriasBaseID int64, dias int) error {
	if dias <= 0 {
		return nil
	}
	feriasBase, err := repository.GetFeriasByID(feriasBaseID)
	if err != nil {
		return fmt.Errorf("erro ao buscar férias base: %w", err)
	}
	if feriasBase == nil {
		return fmt.Errorf("férias base não encontradas")
	}

	periodos, err := repository.GetFeriasNaoPagasComSaldo(feriasBase.FuncionarioID)
	if err != nil {
		return fmt.Errorf("erro ao listar períodos de férias: %w", err)
	}
	if len(periodos) == 0 {
		return fmt.Errorf("não há períodos de férias em aberto")
	}

	inicio := -1
	for i, f := range periodos {
		if f != nil && f.ID == feriasBaseID {
			inicio = i
			break
		}
	}
	if inicio < 0 {
		return fmt.Errorf("férias base id=%d não está em aberto para consumo", feriasBaseID)
	}

	restantes := dias
	for i := inicio; i < len(periodos) && restantes > 0; i++ {
		f := periodos[i]
		if f == nil || f.Pago || f.Dias <= 0 {
			continue
		}
		consome := f.Dias
		if consome > restantes {
			consome = restantes
		}
		if err := repository.ConsumirDiasFerias(f.ID, consome); err != nil {
			return fmt.Errorf("erro ao consumir dias de férias id=%d: %w", f.ID, err)
		}
		restantes -= consome
	}

	if restantes > 0 {
		var alvo *entity.Ferias
		for i := len(periodos) - 1; i >= 0; i-- {
			if periodos[i] != nil && !periodos[i].Pago {
				alvo = periodos[i]
				break
			}
		}
		if alvo == nil {
			return fmt.Errorf("não há período de férias válido para lançar saldo negativo")
		}
		if err := repository.ConsumirDiasFerias(alvo.ID, restantes); err != nil {
			return fmt.Errorf("erro ao lançar saldo negativo de férias: %w", err)
		}
	}

	return nil
}

func (s *DescansoService) CreateDescanso(ctx context.Context, claims Claims, d *entity.Descanso) error {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return err
	}
	if d.Inicio.After(d.Fim) {
		return fmt.Errorf("data inicial não pode ser depois da final")
	}
	ferias, err := repository.GetFeriasByID(d.FeriasID)
	if err != nil {
		return fmt.Errorf("erro ao buscar férias vinculadas: %w", err)
	}
	if ferias == nil {
		return fmt.Errorf("férias não encontradas para ID=%d", d.FeriasID)
	}
	diasDescanso := d.DuracaoEmDias()
	if diasDescanso <= 0 {
		return fmt.Errorf("duração do descanso inválida")
	}
	diasRef := diasReferenciaFerias(ferias)
	valorBasePorDia := ferias.Valor / float64(diasRef)
	var tercoPorDia float64
	if ferias.Terco > 0 {
		tercoPorDia = ferias.Terco / float64(diasRef)
	} else {
		tercoPorDia = valorBasePorDia / 3.0
	}
	d.Valor = (valorBasePorDia + tercoPorDia) * float64(diasDescanso)
	d.Aprovado = false
	d.Pago = false

	if err := s.repo.Create(d); err != nil {
		return fmt.Errorf("erro ao criar descanso: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3,
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso criado ID=%d FeriasID=%d Dias=%d Valor=%.2f", d.ID, d.FeriasID, diasDescanso, d.Valor),
	})
	return nil
}

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
	if descanso.Aprovado {
		return fmt.Errorf("descanso já está aprovado")
	}
	// Consome os dias mantendo a solicitação única, distribuindo por períodos quando necessário.
	if err := s.consumirDiasDescansoUnico(descanso.FeriasID, descanso.DuracaoEmDias()); err != nil {
		return err
	}

	// Aprova somente após consumo bem sucedido, evitando descanso aprovado sem débito de saldo.
	descanso.Aprovado = true
	if err := s.repo.Update(descanso); err != nil {
		return fmt.Errorf("erro ao aprovar descanso: %w", err)
	}

	// Consome os dias mantendo a solicitação única, distribuindo por períodos quando necessário.
	if err := s.consumirDiasDescansoUnico(descanso.FeriasID, descanso.DuracaoEmDias()); err != nil {
		return err
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso aprovado ID=%d (consumiu %d dia(s))", id, descanso.DuracaoEmDias()),
	})
	return nil
}

func (s *DescansoService) MarcarComoPago(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return err
	}
	descanso, err := s.repo.GetDescansoByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar descanso: %w", err)
	}
	if descanso == nil {
		return fmt.Errorf("descanso não encontrado")
	}
	if descanso.Pago {
		return fmt.Errorf("descanso já está pago")
	}
	descanso.Pago = true
	if err := s.repo.Update(descanso); err != nil {
		return fmt.Errorf("erro ao marcar descanso como pago: %w", err)
	}

	// Tenta fechar a férias se não há mais saldo e terço já pago
	if f, ferr := repository.GetFeriasByID(descanso.FeriasID); ferr == nil && f != nil {
		if f.Dias == 0 && f.TercoPago && !f.Pago {
			_ = repository.MarcarFeriasComoPagas(f.ID) // seta pago=true, tercoPago=true
		}
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  4,
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso pago ID=%d", id),
	})
	return nil
}

func (s *DescansoService) DesmarcarPago(ctx context.Context, claims Claims, id int64) error {
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
	if !descanso.Pago {
		return fmt.Errorf("descanso já está não-pago")
	}
	descanso.Pago = false
	if err := s.repo.Update(descanso); err != nil {
		return fmt.Errorf("erro ao desmarcar pagamento do descanso: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{EventoID: 4, UsuarioID: &claims.UserID, Quando: time.Now(), Detalhe: fmt.Sprintf("Descanso desmarcado como pago ID=%d", id)})
	return nil
}

func (s *DescansoService) ListarPorFerias(ctx context.Context, claims Claims, feriasID int64) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosByFeriasID(feriasID)
}

func (s *DescansoService) ListarPorFuncionario(ctx context.Context, claims Claims, funcionarioID int64) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosByFuncionarioID(funcionarioID)
}

func (s *DescansoService) ListarAprovados(ctx context.Context, claims Claims) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosAprovados()
}

func (s *DescansoService) ListarPendentes(ctx context.Context, claims Claims) ([]*entity.Descanso, error) {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return nil, err
	}
	return s.repo.GetDescansosPendentes()
}

func (s *DescansoService) DeleteDescanso(ctx context.Context, claims Claims, id int64) error {
	if err := s.authService.Authorize(ctx, claims, "descanso:delete"); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar descanso: %w", err)
	}
	_, _ = s.logRepo.Create(ctx, LogEntry{EventoID: 5, UsuarioID: &claims.UserID, Quando: time.Now(), Detalhe: fmt.Sprintf("Descanso deletado ID=%d", id)})
	return nil
}

// Cria uma única solicitação de descanso; o consumo dos dias ocorre apenas na aprovação.
func (s *DescansoService) CreateDescansoAuto(ctx context.Context, claims Claims, funcionarioID int64, inicio, fim time.Time) error {
	if err := s.authService.Authorize(ctx, claims, ""); err != nil {
		return err
	}
	if fim.Before(inicio) {
		return fmt.Errorf("data final não pode ser antes da inicial")
	}
	totalDias := int(fim.Sub(inicio).Hours()/24) + 1
	if totalDias <= 0 {
		return fmt.Errorf("duração do descanso inválida")
	}

	periodos, err := repository.GetFeriasNaoPagasComSaldo(funcionarioID)
	if err != nil {
		return fmt.Errorf("erro ao listar períodos de férias: %w", err)
	}
	if len(periodos) == 0 {
		return fmt.Errorf("não há períodos disponíveis para solicitação")
	}

	base := periodos[0]
	for _, f := range periodos {
		if f != nil && !f.Pago && f.Dias > 0 {
			base = f
			break
		}
	}
	if base == nil || base.Pago {
		return fmt.Errorf("não foi possível definir período base para a solicitação")
	}

	diasRef := diasReferenciaFerias(base)
	valorBaseDia := base.Valor / float64(diasRef)
	tercoDia := base.Terco / float64(diasRef)

	d := &entity.Descanso{
		FeriasID: base.ID,
		Inicio:   inicio,
		Fim:      fim,
		Valor:    (valorBaseDia + tercoDia) * float64(totalDias),
		Aprovado: false,
		Pago:     false,
	}
	if err := s.repo.Create(d); err != nil {
		return fmt.Errorf("erro ao criar solicitação de descanso: %w", err)
	}

	_, _ = s.logRepo.Create(ctx, LogEntry{
		EventoID:  3,
		UsuarioID: &claims.UserID,
		Quando:    time.Now(),
		Detalhe:   fmt.Sprintf("Descanso único criado ID=%d FeriasBaseID=%d Dias=%d", d.ID, base.ID, totalDias),
	})
	return nil
}
