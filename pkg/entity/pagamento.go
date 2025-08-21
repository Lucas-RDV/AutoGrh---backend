package entity

import "time"

// Pagamento representa um pagamento realizado a um funcionário
// Pode estar vinculado a uma folha de pagamento e tem um tipo identificado

type Pagamento struct {
	ID            int64     `json:"id"`
	FuncionarioID int64     `json:"funcionario_id"`
	FolhaID       int64     `json:"folha_id"`
	TipoID        int64     `json:"tipo_id"`
	Tipo          string    `json:"tipo"`
	Data          time.Time `json:"data"`
	Valor         float64   `json:"valor"`
}

// NewPagamento cria uma nova instância de Pagamento com tipo, data e valor
func NewPagamento(tipoID int64, data time.Time, valor float64) *Pagamento {
	return &Pagamento{
		TipoID: tipoID,
		Data:   data,
		Valor:  valor,
	}
}
