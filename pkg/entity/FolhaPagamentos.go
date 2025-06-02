package entity

import "time"

type FolhaPagamentos struct {
	Id         int64
	Data       time.Time
	Valor      float64
	Pagamentos []Pagamento
}

func NewFolhaPagamentos(Data time.Time) *FolhaPagamentos {
	f := new(FolhaPagamentos)
	f.Data = Data
	f.Valor = 0
	return f
}
