package entity

import "time"

type Ferias struct {
	Id            int64
	FuncionarioID int64
	Dias          int
	Inicio        time.Time
	Vencimento    time.Time
	Vencido       bool
	Valor         float64
	Descansos     []Descanso
}

func NewFerias(Inicio time.Time) *Ferias {
	f := new(Ferias)
	f.Inicio = Inicio
	f.Vencimento = Inicio.AddDate(1, 0, 0)
	return f
}

func (f *Ferias) DiasUtilizados() int {
	total := 0
	for _, d := range f.Descansos {
		total += d.DuracaoEmDias()
	}
	return total
}

func (f *Ferias) DiasRestantes() int {
	return f.Dias - f.DiasUtilizados()
}
