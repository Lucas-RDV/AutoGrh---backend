package Entity

import "time"

// FolhaPagamentos representa uma folha de pagamento mensal contendo vários pagamentos
// Armazena a data da folha, valor total e os pagamentos associados

type FolhaPagamentos struct {
	ID         int64       `json:"id"`
	Data       time.Time   `json:"data"`
	Valor      float64     `json:"valor"`
	Pagamentos []Pagamento `json:"pagamentos,omitempty"`
}

// NewFolhaPagamentos cria uma nova folha de pagamento com valor inicial zerado
func NewFolhaPagamentos(data time.Time) *FolhaPagamentos {
	return &FolhaPagamentos{
		Data:  data,
		Valor: 0,
	}
}
