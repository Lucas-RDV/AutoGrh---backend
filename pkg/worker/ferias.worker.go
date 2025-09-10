package worker

import (
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	"context"
	"fmt"
	"time"
)

// Worker responsável por rotinas automáticas de férias
type FeriasWorker struct {
	feriasSvc      *service.FeriasService
	descansoSvc    *service.DescansoService
	salarioRealSvc *service.SalarioRealService
	funcionarioSvc *service.FuncionarioService
	faltaSvc       *service.FaltaService
	claims         service.Claims
}

// Construtor
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
		// Executa imediatamente
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		w.run(ctx)
		cancel()

		// Calcula próximo horário alvo (05:00)
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		fmt.Printf("[Worker Férias] Próxima execução agendada para: %v\n", next)

		// Dorme até 05:00
		time.Sleep(time.Until(next))

		// Executa às 05:00
		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Minute)
		w.run(ctx2)
		cancel2()

		// Inicia ticker diário
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			w.run(ctx)
			cancel()
		}
	}()
}

func (w *FeriasWorker) run(ctx context.Context) {
	fmt.Println("[Worker Férias] Iniciando ciclo automático...")

	if err := w.gerarFeriasAnuais(ctx); err != nil {
		fmt.Println("[Worker Férias] Erro ao gerar férias:", err)
	}
	if err := w.marcarFeriasVencidas(ctx); err != nil {
		fmt.Println("[Worker Férias] Erro ao marcar vencidas:", err)
	}

	fmt.Println("[Worker Férias] Ciclo concluído.")
}

// Inicia o worker em uma goroutine
func (w *FeriasWorker) gerarFeriasAnuais(ctx context.Context) error {
	fmt.Println("[Worker Férias] Gerando férias anuais automaticamente...")

	funcionarios, err := w.funcionarioSvc.ListFuncionariosAtivos(ctx, w.claims)
	if err != nil {
		return fmt.Errorf("erro ao listar funcionários ativos: %w", err)
	}

	for _, f := range funcionarios {
		// Definir a data-base: última férias (mais recente) ou admissão
		feriasDoFunc, err := w.feriasSvc.GetFeriasByFuncionarioID(ctx, w.claims, f.ID)
		if err != nil {
			fmt.Printf("[Worker Férias] erro ao buscar férias do funcionário %d: %v\n", f.ID, err)
			continue
		}

		var base time.Time
		if len(feriasDoFunc) == 0 {
			base = f.Admissao
		} else {
			ultima := feriasDoFunc[0]
			for _, ff := range feriasDoFunc[1:] {
				if ff.Inicio.After(ultima.Inicio) {
					ultima = ff
				}
			}
			base = ultima.Inicio
		}

		// Já completou 12 meses desde a base?
		if time.Since(base) < 365*24*time.Hour {
			continue
		}

		// Período aquisitivo: [base, base+12m)
		inicioPeriodo := base
		fimPeriodo := base.AddDate(1, 0, 0)

		// Buscar faltas do funcionário e somar Quantidade dentro do período
		faltas, err := w.faltaSvc.GetFaltasByFuncionarioID(ctx, w.claims, f.ID)
		if err != nil {
			fmt.Printf("[Worker Férias] erro ao buscar faltas do funcionário %d: %v\n", f.ID, err)
			continue
		}

		totalFaltasPeriodo := 0
		for _, fal := range faltas {
			// Falta é mensal: considerar somente meses dentro do período aquisitivo
			if !fal.Mes.Before(inicioPeriodo) && fal.Mes.Before(fimPeriodo) {
				totalFaltasPeriodo += fal.Quantidade
			}
		}

		// Determinar dias de direito conforme total de faltas
		dias := 30
		switch {
		case totalFaltasPeriodo >= 33:
			dias = 0
		case totalFaltasPeriodo >= 24:
			dias = 12
		case totalFaltasPeriodo >= 15:
			dias = 18
		case totalFaltasPeriodo >= 6:
			dias = 24
		default:
			dias = 30
		}
		if dias == 0 {
			fmt.Printf("[Worker Férias] funcionarioID=%d perdeu direito (faltas=%d)\n", f.ID, totalFaltasPeriodo)
			continue
		}

		// Salário real atual
		salario, err := w.salarioRealSvc.GetSalarioRealAtual(ctx, w.claims, f.ID)
		if err != nil || salario == nil {
			fmt.Printf("[Worker Férias] salário real ausente/erro para funcionarioID=%d\n", f.ID)
			continue
		}

		// Calcular valor (dias + 1/3)
		valorBase := (salario.Valor / 30.0) * float64(dias)
		valorTerco := valorBase / 3.0
		total := valorBase + valorTerco

		// Criar férias (início = agora; ajuste conforme sua regra)
		inicio := time.Now()
		if _, err := w.feriasSvc.CriarFerias(ctx, w.claims, f.ID, dias, total, inicio); err != nil {
			fmt.Printf("[Worker Férias] erro ao criar férias do funcionarioID=%d: %v\n", f.ID, err)
			continue
		}

		fmt.Printf("[Worker Férias] Férias criadas: funcionarioID=%d dias=%d faltas=%d valor=%.2f\n",
			f.ID, dias, totalFaltasPeriodo, total)
	}

	return nil
}

func (w *FeriasWorker) marcarFeriasVencidas(ctx context.Context) error {
	fmt.Println("[Worker Férias] Marcando férias vencidas automaticamente...")

	// Buscar todas as férias
	feriasList, err := repository.ListFerias()
	if err != nil {
		return fmt.Errorf("erro ao listar férias: %w", err)
	}

	for _, f := range feriasList {
		if f.Vencido {
			continue
		}
		if f.Vencimento.Before(time.Now()) {
			f.Vencido = true
			if err := repository.UpdateFerias(f); err != nil {
				fmt.Printf("[Worker Férias] erro ao atualizar férias ID=%d: %v\n", f.ID, err)
				continue
			}
			fmt.Printf("[Worker Férias] Férias ID=%d do funcionarioID=%d marcadas como vencidas\n",
				f.ID, f.FuncionarioID)
		}
	}

	return nil
}
