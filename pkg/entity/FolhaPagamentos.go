package entity

import "time"

type FolhaPagamentos struct {
	Id         string
	Tipo       string
	Data       time.Time
	Valor      float64
	Pagamentos []Pagamento
}

func NewFolhaPagamentos(Tipo string, Data time.Time) *FolhaPagamentos {
	f := new(FolhaPagamentos)
	f.Tipo = Tipo
	f.Data = Data
	f.Valor = 0
	return f
}
