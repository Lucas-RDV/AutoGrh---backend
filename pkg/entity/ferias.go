package entity

import "time"

// Ferias representa o direito a um determinado período de descanso de um funcionário.
// Inclui os descansos efetivamente gozados e permite calcular dias restantes e valores.
type Ferias struct {
	ID            int64      `json:"id"`
	FuncionarioID int64      `json:"funcionario_id"`
	Dias          int        `json:"dias"`
	Inicio        time.Time  `json:"inicio"`
	Vencimento    time.Time  `json:"vencimento"`
	Vencido       bool       `json:"vencido"`
	Valor         float64    `json:"valor"`
	Descansos     []Descanso `json:"descansos,omitempty"`
	Pago          bool       `json:"pago"`
	Terco         float64    `json:"terco"`
	TercoPago     bool       `json:"tercoPago"`
}

// NewFerias cria uma nova instância de Ferias com vencimento um ano após a data de início.
// Os dias devem ser informados já considerando a regra da CLT e faltas injustificadas.
func NewFerias(funcionarioID int64, inicio time.Time, dias int) *Ferias {
	return &Ferias{
		FuncionarioID: funcionarioID,
		Dias:          dias,
		Inicio:        inicio,
		Vencimento:    inicio.AddDate(1, 0, 0),
		Vencido:       false,
		Pago:          false,
		TercoPago:     false,
		Descansos:     []Descanso{},
	}
}

// DiasUtilizados retorna a soma total dos dias utilizados em descansos.
func (f *Ferias) DiasUtilizados() int {
	total := 0
	for _, d := range f.Descansos {
		total += d.DuracaoEmDias()
	}
	return total
}

// DiasRestantes calcula quantos dias de férias ainda estão disponíveis.
func (f *Ferias) DiasRestantes() int {
	return f.Dias - f.DiasUtilizados()
}

// CalcularValor calcula o valor das férias e o adicional de 1/3.
// Retorna: valor base, valor do terço, valor total.
func (f *Ferias) CalcularValor(salario float64) (float64, float64, float64) {
	valor := (salario / 30.0) * float64(f.Dias)
	terco := salario / 3.0
	total := valor + terco

	f.Valor = valor
	f.Terco = terco
	return valor, terco, total
}
