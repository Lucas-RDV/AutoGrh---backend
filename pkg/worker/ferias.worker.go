package worker

import (
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	"context"
	"fmt"
	"time"
)

type FeriasWorker struct {
	feriasSvc      *service.FeriasService
	descansoSvc    *service.DescansoService
	salarioRealSvc *service.SalarioRealService
	funcionarioSvc *service.FuncionarioService
	faltaSvc       *service.FaltaService
	claims         service.Claims
}

func NewFeriasWorker(
	feriasSvc *service.FeriasService,
	descansoSvc *service.DescansoService,
	salarioRealSvc *service.SalarioRealService,
	funcionarioSvc *service.FuncionarioService,
	faltaSvc *service.FaltaService,
) *FeriasWorker {
	return &FeriasWorker{
		feriasSvc:      feriasSvc,
		descansoSvc:    descansoSvc,
		salarioRealSvc: salarioRealSvc,
		funcionarioSvc: funcionarioSvc,
		faltaSvc:       faltaSvc,
		claims:         middleware.SystemClaims(),
	}
}

func (w *FeriasWorker) Start() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		w.run(ctx)
		cancel()

		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		fmt.Printf("[Worker Férias] Próxima execução agendada para: %v\n", next)

		time.Sleep(time.Until(next))

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Minute)
		w.run(ctx2)
		cancel2()

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			w.run(ctx)
			cancel()
		}
	}()
}

func (w *FeriasWorker) run(ctx context.Context) {
	fmt.Println("[Worker Férias] Iniciando ciclo automático...")

	// Garante períodos
	if err := w.garantirParaFuncionariosAtivos(ctx); err != nil {
		fmt.Println("[Worker Férias] Erro ao garantir férias:", err)
	}

	// Marca descansos como pagos quando a data 'fim' já passou (e já aprovados)
	if err := w.pagarDescansosComFimAnteriorAHoje(ctx); err != nil {
		fmt.Println("[Worker Férias] Erro ao pagar descansos vencidos:", err)
	}

	fmt.Println("[Worker Férias] Ciclo concluído.")
}

func (w *FeriasWorker) garantirParaFuncionariosAtivos(ctx context.Context) error {
	funcionarios, err := w.funcionarioSvc.ListFuncionariosAtivos(ctx, w.claims)
	if err != nil {
		return fmt.Errorf("erro ao listar funcionários ativos: %w", err)
	}
	for _, f := range funcionarios {
		if _, err := w.feriasSvc.GarantirFeriasAteHoje(ctx, w.claims, f.ID); err != nil {
			fmt.Printf("[Worker Férias] funcionarioID=%d erro: %v\n", f.ID, err)
		}
	}
	return nil
}

func (w *FeriasWorker) pagarDescansosComFimAnteriorAHoje(ctx context.Context) error {
	list, err := w.descansoSvc.ListarAprovados(ctx, w.claims)
	if err != nil {
		return err
	}
	hoje := time.Now()
	for _, d := range list {
		if !d.Pago && d.Fim.Before(hoje) {
			// paga o descanso
			if err := w.descansoSvc.MarcarComoPago(ctx, w.claims, d.ID); err != nil {
				fmt.Printf("[Worker Férias] Falha ao pagar descanso ID=%d: %v\n", d.ID, err)
				continue
			}
			// após pagar, verifica e fecha férias se elegível
			if f, ferr := repository.GetFeriasByID(d.FeriasID); ferr == nil && f != nil {
				if f.Dias == 0 && f.TercoPago && !f.Pago {
					if err := repository.MarcarFeriasComoPagas(f.ID); err != nil {
						fmt.Printf("[Worker Férias] Falha ao marcar férias %d como pagas: %v\n", f.ID, err)
					}
				}
			}
		}
	}
	return nil
}
