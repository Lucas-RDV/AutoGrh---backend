package entity

import "time"

type Falta struct {
	Id            int64
	FuncionarioId int64
	Quantidade    int
	Mes           time.Time
}

func newFalta(qtd int, mes time.Time, funcionarioId int64) *Falta {
	f := new(Falta)
	f.Quantidade = qtd
	f.Mes = mes
	f.FuncionarioId = funcionarioId
	return f
}
