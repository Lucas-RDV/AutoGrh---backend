package entity

import "time"

type Pagamento struct {
	Id            int64
	FuncionarioId int64
	FolhaId       int64
	TipoId        int64
	Tipo          string
	Data          time.Time
	Valor         float64
}

func NewPagamento(TipoId int64, Data time.Time, Valor float64) *Pagamento {
	d := new(Pagamento)
	d.TipoId = TipoId
	d.Data = Data
	d.Valor = Valor
	return d
}
