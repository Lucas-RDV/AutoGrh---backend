package worker

import (
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"context"
	"fmt"
	"time"
)

// Worker responsável por rotinas automáticas da folha de pagamento
type FolhaWorker struct {
	folhaSvc *service.FolhaPagamentoService
	claims   service.Claims
}

// Construtor
func NewFolhaWorker(folhaSvc *service.FolhaPagamentoService) *FolhaWorker {
	return &FolhaWorker{
		folhaSvc: folhaSvc,
		claims:   middleware.SystemClaims(),
	}
}

func (w *FolhaWorker) Start() {
	go func() {
		// Executa imediatamente ao iniciar
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		w.run(ctx)
		cancel()

		// Calcula próximo horário alvo (05:00)
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		fmt.Printf("[Worker Folha] Próxima execução agendada para: %v\n", next)

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

func (w *FolhaWorker) run(ctx context.Context) {
	fmt.Println("[Worker Folha] Iniciando ciclo automático...")

	if err := w.gerarOuRecalcularFolhaSalario(ctx); err != nil {
		fmt.Println("[Worker Folha] Erro ao processar folha de salário:", err)
	}

	fmt.Println("[Worker Folha] Ciclo concluído.")
}

// gerarOuRecalcularFolhaSalario cria ou recalcula a folha de salário do mês atual
func (w *FolhaWorker) gerarOuRecalcularFolhaSalario(ctx context.Context) error {
	now := time.Now()
	mes := int(now.Month())
	ano := now.Year()

	// Busca se já existe folha de salário para o mês/ano atuais
	folha, err := w.folhaSvc.BuscarFolhaPorMesAnoTipo(ctx, w.claims, mes, ano, "SALARIO")
	if err != nil {
		return fmt.Errorf("erro ao buscar folha existente: %w", err)
	}

	if folha == nil {
		// Não existe: cria nova
		if _, err := w.folhaSvc.CriarFolhaSalario(ctx, w.claims, mes, ano); err != nil {
			return fmt.Errorf("erro ao criar folha salário %02d/%d: %w", mes, ano, err)
		}
		fmt.Printf("[Worker Folha] Folha salário %02d/%d criada com sucesso\n", mes, ano)
	} else if !folha.Pago {
		// Existe e não está paga: recalcula
		if err := w.folhaSvc.RecalcularFolha(ctx, w.claims, folha.ID); err != nil {
			return fmt.Errorf("erro ao recalcular folha salário %02d/%d: %w", mes, ano, err)
		}
		fmt.Printf("[Worker Folha] Folha salário %02d/%d recalculada com sucesso\n", mes, ano)
	} else {
		// Já paga: ignora
		fmt.Printf("[Worker Folha] Folha salário %02d/%d já está paga, nenhuma ação necessária\n", mes, ano)
	}

	return nil
}
