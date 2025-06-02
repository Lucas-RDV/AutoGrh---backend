package entity

import "time"

type Salario struct {
	Id    int64
	Valor float64
	Ano   time.Time
}

func newSalario(Valor float64, Ano time.Time) *Salario {
	s := new(Salario)
	s.Valor = Valor
	s.Ano = Ano
	return s
}
