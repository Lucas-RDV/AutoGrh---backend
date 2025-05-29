package entity

import "time"

type Salario struct {
	Id    string
	Valor float64
	Ano   time.Time
}

func newSalario(Valor float64, Ano time.Time) *Salario {
	s := new(Salario)
	s.Valor = Valor
	s.Ano = Ano
	return s
}
