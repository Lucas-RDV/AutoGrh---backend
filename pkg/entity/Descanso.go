package entity

import "time"

type Descanso struct {
	Id       string
	Inicio   time.Time
	Fim      time.Time
	Duracao  int
	Valor    float64
	Aprovado bool
	Pago     bool
}

func newDescanso(Inicio time.Time, Duracao int) *Descanso {
	d := new(Descanso)
	d.Inicio = Inicio
	d.Fim = Inicio.AddDate(0, 0, Duracao)
	d.Aprovado = false
	d.Pago = false
	return d
}
