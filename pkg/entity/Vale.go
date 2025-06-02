package entity

import "time"

type Vale struct {
	Id            int64
	FuncionarioId int64
	Valor         float64
	Data          time.Time
	Aprovado      bool
	Pago          bool
}

func NewVale(valor float64, Data time.Time) *Vale {
	v := new(Vale)
	v.Valor = valor
	v.Data = Data
	v.Aprovado = false
	v.Pago = false
	return v
}
