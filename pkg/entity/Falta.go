package entity

import "time"

type Falta struct {
	Id         string
	Quantidade int
	mes        time.Time
}

func newFalta(qtd int, mes time.Time) *Falta {
	f := new(Falta)
	f.Quantidade = qtd
	f.mes = mes
	return f
}
