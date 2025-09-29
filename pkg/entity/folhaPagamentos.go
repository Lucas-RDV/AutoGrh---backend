package entity

import "time"

// FolhaPagamentos representa uma folha mensal de pagamento ou de vale
type FolhaPagamentos struct {
	ID          int64     `json:"id"`
	Mes         int       `json:"mes"`         // mês de referência (1-12)
	Ano         int       `json:"ano"`         // ano de referência
	Tipo        string    `json:"tipo"`        // "SALARIO" ou "VALE"
	DataGeracao time.Time `json:"dataGeracao"` // quando a folha foi criada
	ValorTotal  float64   `json:"valorTotal"`  // somatório dos pagamentos da folha
	Pago        bool      `json:"pago"`        // indica se a folha foi fechada/paga
}

// NewFolhaPagamentos cria uma nova folha com valor inicial zerado.
func NewFolhaPagamentos(mes int, ano int, tipo string) *FolhaPagamentos {
	return &FolhaPagamentos{
		Mes:         mes,
		Ano:         ano,
		Tipo:        tipo,
		DataGeracao: time.Now(),
		ValorTotal:  0,
		Pago:        false,
	}
}
