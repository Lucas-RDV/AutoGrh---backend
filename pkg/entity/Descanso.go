package entity

import "time"

type Descanso struct {
	Id       int64
	Inicio   time.Time
	Fim      time.Time
	Valor    float64
	Aprovado bool
	Pago     bool
	FeriasID int64
}

func NewDescanso(Inicio time.Time, Fim time.Time, FeriasID int64) *Descanso {
	d := new(Descanso)
	d.Inicio = Inicio
	d.Fim = Fim
	d.Aprovado = false
	d.Pago = false
	d.FeriasID = FeriasID
	return d
}

func (d *Descanso) DuracaoEmDias() int {
	return int(d.Fim.Sub(d.Inicio).Hours()/24) + 1
}
