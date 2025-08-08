package Entity

import "time"

// Vale representa um adiantamento salarial solicitado por um funcionário
// Inclui dados sobre aprovação, pagamento e data do pedido

type Vale struct {
	ID            int64     `json:"id"`
	FuncionarioID int64     `json:"funcionario_id"`
	Valor         float64   `json:"valor"`
	Data          time.Time `json:"data"`
	Aprovado      bool      `json:"aprovado"`
	Pago          bool      `json:"pago"`
}

// NewVale cria uma nova instância de Vale com aprovação e pagamento desabilitados
func NewVale(funcionarioID int64, valor float64, data time.Time) *Vale {
	return &Vale{
		FuncionarioID: funcionarioID,
		Valor:         valor,
		Data:          data,
		Aprovado:      false,
		Pago:          false,
	}
}
