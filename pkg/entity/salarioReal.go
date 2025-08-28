package entity

import "time"

// SalarioReal representa o salário líquido efetivamente recebido por um funcionário.
// Diferente do Salario (registrado em carteira), este é o valor manual usado para
// cálculos de folha, férias e descontos.
// Mantém histórico de alterações por período (inicio/fim).
type SalarioReal struct {
	ID            int64      `json:"id"`
	FuncionarioID int64      `json:"funcionario_id"`
	Inicio        time.Time  `json:"inicio"`
	Fim           *time.Time `json:"fim,omitempty"`
	Valor         float64    `json:"valor"`
}

// NewSalarioReal cria uma nova instância de SalarioReal, com Fim = nil (ativo).
func NewSalarioReal(funcionarioID int64, inicio time.Time, valor float64) *SalarioReal {
	return &SalarioReal{
		FuncionarioID: funcionarioID,
		Inicio:        inicio,
		Valor:         valor,
		Fim:           nil,
	}
}
