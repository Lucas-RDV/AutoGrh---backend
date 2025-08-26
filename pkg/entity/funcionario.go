package entity

import "time"

// Funcionario representa um vínculo contratual com uma pessoa
// Dados pessoais são referenciados via PessoaID; este modelo armazena dados contratuais

type Funcionario struct {
	ID                int64      `json:"id"`
	PessoaID          int64      `json:"pessoa_id"`
	PIS               string     `json:"pis"`
	CTPF              string     `json:"ctpf"`
	Nascimento        time.Time  `json:"nascimento"`
	Admissao          time.Time  `json:"admissao"`
	Demissao          *time.Time `json:"demissao,omitempty"`
	Cargo             string     `json:"cargo"`
	SalarioInicial    float64    `json:"salario_inicial"`
	FeriasDisponiveis int        `json:"ferias_disponiveis"`
	Ativo             bool       `json:"ativo"`

	Salarios   []Salario   `json:"salarios,omitempty"`
	Documentos []Documento `json:"documentos,omitempty"`
	Ferias     []Ferias    `json:"ferias,omitempty"`
	Faltas     []Falta     `json:"faltas,omitempty"`
	Pagamentos []Pagamento `json:"pagamentos,omitempty"`
	Vales      []Vale      `json:"vales,omitempty"`
}

// NewFuncionario cria uma nova instância de Funcionario com listas vazias e ativo = true
func NewFuncionario(
	pessoaID int64,
	pis, ctpf, cargo string,
	nascimento, admissao time.Time,
	salarioInicial float64,
) *Funcionario {
	return &Funcionario{
		PessoaID:          pessoaID,
		PIS:               pis,
		CTPF:              ctpf,
		Nascimento:        nascimento,
		Admissao:          admissao,
		Cargo:             cargo,
		SalarioInicial:    salarioInicial,
		Demissao:          nil,
		FeriasDisponiveis: 0,
		Ativo:             true,
		Salarios:          []Salario{},
		Documentos:        []Documento{},
		Ferias:            []Ferias{},
		Faltas:            []Falta{},
		Pagamentos:        []Pagamento{},
		Vales:             []Vale{},
	}
}
