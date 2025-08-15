package Entity

import "time"

// Ferias representa o direito a um determinado período de descanso de um funcionário
// Inclui os descansos efetivamente gozados e permite calcular dias restantes

type Ferias struct {
	ID            int64      `json:"id"`
	FuncionarioID int64      `json:"funcionario_id"`
	Dias          int        `json:"dias"`
	Inicio        time.Time  `json:"inicio"`
	Vencimento    time.Time  `json:"vencimento"`
	Vencido       bool       `json:"vencido"`
	Valor         float64    `json:"valor"`
	Descansos     []Descanso `json:"descansos,omitempty"`
}

// NewFerias cria uma nova instância de Ferias com vencimento de um ano após a data de início
func NewFerias(inicio time.Time) *Ferias {
	return &Ferias{
		Inicio:     inicio,
		Vencimento: inicio.AddDate(1, 0, 0),
	}
}

// DiasUtilizados retorna a soma total dos dias utilizados em descansos
func (f *Ferias) DiasUtilizados() int {
	total := 0
	for _, d := range f.Descansos {
		total += d.DuracaoEmDias()
	}
	return total
}

// DiasRestantes calcula quantos dias de férias ainda estão disponíveis
func (f *Ferias) DiasRestantes() int {
	return f.Dias - f.DiasUtilizados()
}
