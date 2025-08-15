package Entity

import "time"

// Salario representa o valor do salário de um funcionário em um determinado período
// Pode conter uma data de fim para indicar alterações salariais ao longo do tempo

type Salario struct {
	ID            int64      `json:"id"`
	FuncionarioID int64      `json:"funcionario_id"`
	Inicio        time.Time  `json:"inicio"`
	Fim           *time.Time `json:"fim,omitempty"`
	Valor         float64    `json:"valor"`
}

// NewSalario cria uma nova instância de Salario com início e valor definidos
func NewSalario(funcionarioID int64, inicio time.Time, valor float64) *Salario {
	return &Salario{
		FuncionarioID: funcionarioID,
		Inicio:        inicio,
		Valor:         valor,
	}
}
