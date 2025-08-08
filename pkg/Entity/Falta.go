package Entity

import "time"

// Falta representa a quantidade de faltas de um funcionário em um determinado mês
// Usado para cálculo de descontos, controle de presença e geração da folha

type Falta struct {
	ID            int64     `json:"id"`
	FuncionarioID int64     `json:"funcionario_id"`
	Quantidade    int       `json:"quantidade"`
	Mes           time.Time `json:"mes"`
}

// NewFalta cria uma nova instância de Falta para um funcionário em um determinado mês
func NewFalta(qtd int, mes time.Time, funcionarioID int64) *Falta {
	return &Falta{
		Quantidade:    qtd,
		Mes:           mes,
		FuncionarioID: funcionarioID,
	}
}
