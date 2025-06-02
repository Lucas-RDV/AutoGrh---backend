package entity

import "time"

type Salario struct {
	Id            int64
	FuncionarioID int64
	Inicio        time.Time
	Fim           *time.Time
	Valor         float64
}

func NewSalario(funcionarioID int64, inicio time.Time, valor float64) *Salario {
	return &Salario{
		FuncionarioID: funcionarioID,
		Inicio:        inicio,
		Valor:         valor,
	}
}
