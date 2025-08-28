package entity

import "time"

// Funcionario representa um vínculo contratual com uma pessoa
// Dados pessoais são referenciados via PessoaID; este modelo armazena dados contratuais

// Funcionario
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

	SalarioRegistradoAtual *Salario     `json:"salario_registrado_atual,omitempty"`
	SalarioRealAtual       *SalarioReal `json:"salario_real_atual,omitempty"`

	SalariosRegistrados []*Salario     `json:"salarios_registrados,omitempty"`
	SalariosReais       []*SalarioReal `json:"salarios_reais,omitempty"`

	Documentos []Documento `json:"documentos,omitempty"`
	Ferias     []Ferias    `json:"ferias,omitempty"`
	Faltas     []*Falta    `json:"faltas,omitempty"`
	Pagamentos []Pagamento `json:"pagamentos,omitempty"`
	Vales      []Vale      `json:"vales,omitempty"`
}
