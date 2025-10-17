package worker

import (
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"context"
	"fmt"
	"time"
)

type AvisosWorker struct {
	avisoSvc *service.AvisoService
	claims   service.Claims
}

func NewAvisosWorker(avisoSvc *service.AvisoService) *AvisosWorker {
	return &AvisosWorker{
		avisoSvc: avisoSvc,
		claims:   middleware.SystemClaims(),
	}
}

func (w *AvisosWorker) Start() {
	go func() {
		// Executa imediatamente ao subir
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		_ = w.avisoSvc.RunDaily(ctx)
		cancel()

		// Agenda para 05:10 diariamente (depois dos outros workers)
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 5, 10, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		fmt.Printf("[Worker Avisos] Próxima execução agendada para: %v\n", next)

		time.Sleep(time.Until(next))

		// Primeira execução no horário-alvo
		ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Minute)
		_ = w.avisoSvc.RunDaily(ctx2)
		cancel2()

		// Ticker diário
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			_ = w.avisoSvc.RunDaily(ctx)
			cancel()
		}
	}()
}
