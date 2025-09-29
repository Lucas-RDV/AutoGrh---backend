package entity

type Pagamento struct {
	ID             int64   `json:"id"`
	FuncionarioID  int64   `json:"funcionarioId"`
	FolhaID        int64   `json:"folhaId"`
	SalarioBase    float64 `json:"salarioBase"`
	Adicional      float64 `json:"adicional"`
	DescontoINSS   float64 `json:"descontoINSS"`
	SalarioFamilia float64 `json:"salarioFamilia"`
	DescontoVales  float64 `json:"descontoVales"`
	ValorFinal     float64 `json:"valorFinal"`
	Pago           bool    `json:"pago"`
}

func NewPagamento(funcionarioID, folhaID int64, salarioBase float64) *Pagamento {
	return &Pagamento{
		FuncionarioID:  funcionarioID,
		FolhaID:        folhaID,
		SalarioBase:    salarioBase,
		Adicional:      0,
		DescontoINSS:   0,
		SalarioFamilia: 0,
		DescontoVales:  0,
		ValorFinal:     salarioBase,
		Pago:           false,
	}
}

func (p *Pagamento) RecalcularValorFinal(descontoFaltas float64) {
	p.ValorFinal = p.SalarioBase +
		p.Adicional +
		p.SalarioFamilia -
		p.DescontoINSS -
		p.DescontoVales -
		descontoFaltas
}
