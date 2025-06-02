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

func NewVale(funcId int64, valor float64, Data time.Time) *Vale {
	v := new(Vale)
	v.Valor = valor
	v.Data = Data
	v.FuncionarioId = funcId
	v.Aprovado = false
	v.Pago = false
	return v
}
