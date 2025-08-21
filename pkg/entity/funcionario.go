package entity

import "time"

// Funcionario representa um colaborador da empresa
// Inclui dados pessoais, histórico de pagamento e relações funcionais
// Todas as listas são iniciadas vazias no construtor

type Funcionario struct {
	ID                int64       `json:"id"`
	Nome              string      `json:"nome"`
	RG                string      `json:"rg"`
	CPF               string      `json:"cpf"`
	PIS               string      `json:"pis"`
	CTPF              string      `json:"ctpf"`
	Endereco          string      `json:"endereco"`
	Contato           string      `json:"contato"`
	ContatoEmergencia string      `json:"contato_emergencia"`
	Nascimento        time.Time   `json:"nascimento"`
	Admissao          time.Time   `json:"admissao"`
	Demissao          *time.Time  `json:"demissao,omitempty"`
	Cargo             string      `json:"cargo"`
	SalarioInicial    float64     `json:"salario_inicial"`
	Salarios          []Salario   `json:"salarios,omitempty"`
	Documentos        []Documento `json:"documentos,omitempty"`
	Ferias            []Ferias    `json:"ferias,omitempty"`
	Faltas            []Falta     `json:"faltas,omitempty"`
	FeriasDisponiveis int         `json:"ferias_disponiveis"`
	Pagamentos        []Pagamento `json:"pagamentos,omitempty"`
	Vales             []Vale      `json:"vales,omitempty"`
}

// NewFuncionario cria uma nova instância de Funcionario com listas vazias e demissão nula
func NewFuncionario(
	nome, rg, cpf, pis, ctpf, endereco, contato, contatoEmergencia, cargo string,
	nascimento, admissao time.Time,
	salarioInicial float64,
) *Funcionario {
	return &Funcionario{
		Nome:              nome,
		RG:                rg,
		CPF:               cpf,
		PIS:               pis,
		CTPF:              ctpf,
		Endereco:          endereco,
		Contato:           contato,
		ContatoEmergencia: contatoEmergencia,
		Nascimento:        nascimento,
		Admissao:          admissao,
		Cargo:             cargo,
		SalarioInicial:    salarioInicial,
		Demissao:          nil,
		Salarios:          []Salario{},
		Documentos:        []Documento{},
		Ferias:            []Ferias{},
		Faltas:            []Falta{},
		Pagamentos:        []Pagamento{},
		Vales:             []Vale{},
	}
}
