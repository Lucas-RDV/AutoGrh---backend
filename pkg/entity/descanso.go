package entity

import "time"

// Descanso representa um período de férias efetivamente gozado por um funcionário
// Está relacionado a uma requisição de férias (FeriasID) e inclui status de aprovação e pagamento

type Descanso struct {
	ID       int64     `json:"id"`
	Inicio   time.Time `json:"inicio"`
	Fim      time.Time `json:"fim"`
	Valor    float64   `json:"valor"`
	Aprovado bool      `json:"aprovado"`
	Pago     bool      `json:"pago"`
	FeriasID int64     `json:"ferias_id"`
}

// NewDescanso cria uma nova instância de Descanso não aprovado nem pago
func NewDescanso(inicio, fim time.Time, feriasID int64) *Descanso {
	return &Descanso{
		Inicio:   inicio,
		Fim:      fim,
		Aprovado: false,
		Pago:     false,
		FeriasID: feriasID,
	}
}

// DuracaoEmDias retorna o número de dias do descanso, incluindo o primeiro e último dias
func (d *Descanso) DuracaoEmDias() int {
	return int(d.Fim.Sub(d.Inicio).Hours()/24) + 1
}
