package service

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"context"
	"fmt"
	"time"
)

type AvisoService struct {
	auth *AuthService
}

func NewAvisoService(auth *AuthService) *AvisoService {
	return &AvisoService{auth: auth}
}

func (s *AvisoService) List(ctx context.Context, claims Claims) ([]entity.Aviso, error) {
	if err := s.auth.Authorize(ctx, claims, ""); err != nil {
		return nil, err
	}
	return repository.ListAvisos()
}

func (s *AvisoService) RunDaily(ctx context.Context) error {
	// failsafe: desativar avisos muito antigos
	_ = repository.DeactivateAvisosVencidosAntes(time.Now().AddDate(-1, 0, 0))

	now := time.Now()
	horizonte := now.AddDate(0, 0, 60) // 60 dias

	// ===== Férias =====
	if ferias, err := repository.ListFerias(); err == nil {
		for _, f := range ferias {
			nome, _ := repository.GetFuncionarioNomeByID(f.FuncionarioID) // se der erro, cai no fallback abaixo

			// Se já pagas, limpa quaisquer avisos relacionados
			if f.Pago {
				_ = repository.DeleteAvisoByTypeAndRef("FERIAS_VENCENDO", f.ID)
				_ = repository.DeleteAvisoByTypeAndRef("FERIAS_VENCIDAS", f.ID)
				continue
			}

			if f.Vencido {
				// FERIAS_VENCIDAS
				ref := f.ID
				label := fmt.Sprintf("Férias de %s", firstOrID(nome, f.FuncionarioID))
				msg := fmt.Sprintf("%s estão vencidas desde %s.",
					label, f.Vencimento.Format("02/01/2006"))
				_ = repository.CreateAviso(&entity.Aviso{
					Tipo:         "FERIAS_VENCIDAS",
					Mensagem:     msg,
					ReferenciaID: &ref,
					CriadoEm:     time.Now(),
					Ativo:        true,
				})
				// trocamos eventual "vencendo" por "vencidas"
				_ = repository.DeleteAvisoByTypeAndRef("FERIAS_VENCENDO", f.ID)
				continue
			}

			// A 60 dias do vencimento:
			if f.Vencimento.After(now) && (f.Vencimento.Equal(horizonte) || f.Vencimento.Before(horizonte)) {
				ref := f.ID
				label := fmt.Sprintf("Férias de %s", firstOrID(nome, f.FuncionarioID))
				msg := fmt.Sprintf("%s vencem em %s.",
					label, f.Vencimento.Format("02/01/2006"))
				_ = repository.CreateAviso(&entity.Aviso{
					Tipo:         "FERIAS_VENCENDO",
					Mensagem:     msg,
					ReferenciaID: &ref,
					CriadoEm:     time.Now(),
					Ativo:        true,
				})
			}
		}
	}

	// ===== Vales aprovados e não pagos =====
	if aprovados, err := repository.ListValesAprovadosNaoPagos(); err == nil {
		for _, v := range aprovados {
			nome, _ := repository.GetFuncionarioNomeByID(v.FuncionarioID)
			ref := v.ID
			msg := fmt.Sprintf("Vale pendente de pagamento (%s, R$ %.2f em %s).",
				firstOrID(nome, v.FuncionarioID), v.Valor, v.Data.Format("02/01/2006"))
			_ = repository.CreateAviso(&entity.Aviso{
				Tipo:         "VALE_PENDENTE",
				Mensagem:     msg,
				ReferenciaID: &ref,
				CriadoEm:     time.Now(),
				Ativo:        true,
			})
		}
	}
	// Vales pagos limpam avisos pendentes
	_ = repository.DeleteAvisosByType("VALE_PENDENTE")

	// ===== Descansos pendentes =====
	if pend, err := repository.GetDescansosPendentes(); err == nil {
		for _, d := range pend {
			// descanso só tem ferias_id; buscamos nome via férias
			nome, _ := repository.GetFuncionarioNomeByFeriasID(d.FeriasID)
			ref := d.ID
			msg := fmt.Sprintf("Descanso pendente de aprovação (%s, %s a %s).",
				firstOrID(nome, 0), d.Inicio.Format("02/01/2006"), d.Fim.Format("02/01/2006"))
			_ = repository.CreateAviso(&entity.Aviso{
				Tipo:         "DESCANSO_PENDENTE",
				Mensagem:     msg,
				ReferenciaID: &ref,
				CriadoEm:     time.Now(),
				Ativo:        true,
			})
		}
	}
	// Descansos aprovados limpam os pendentes
	_ = repository.DeleteAvisosByType("DESCANSO_PENDENTE")

	return nil
}

// firstOrID: se nome estiver vazio, retorna "funcionário <ID>"
func firstOrID(nome string, id int64) string {
	if nome != "" {
		return nome
	}
	if id > 0 {
		return fmt.Sprintf("funcionário %d", id)
	}
	return "funcionário"
}
