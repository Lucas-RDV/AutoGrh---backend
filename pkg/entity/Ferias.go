package entity

import "time"

type Ferias struct {
	Id         string
	Dias       int
	Inicio     time.Time
	Vencimento time.Time
	Vencido    bool
	Valor      float64
	Descansos  []Descanso
}

func newFerias(Inicio time.Time) *Ferias {
	f := new(Ferias)
	f.Inicio = Inicio
	f.Vencimento = Inicio.AddDate(1, 0, 0)
	return f
}
