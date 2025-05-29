package entity

import "time"

type Pagamento struct {
	Id    string
	Tipo  string
	Data  time.Time
	Valor float64
}

func NewPagamento(Tipo string, Data time.Time, Valor float64) *Pagamento {
	d := new(Pagamento)
	d.Tipo = Tipo
	d.Data = Data
	d.Valor = Valor
	return d
}
