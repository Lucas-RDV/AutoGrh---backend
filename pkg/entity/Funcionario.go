package entity

import "time"

type Funcionario struct {
	Id                string
	Nome              string
	RG                string
	CPF               string
	PIS               string
	CTPF              string
	Endereco          string
	Contato           string
	ContatoEmergencia string
	Nascimento        time.Time
	Admissao          time.Time
	Demissao          time.Time
	Cargo             string
	SalarioInicial    float64
	Salarios          []Salario
	Documentos        []Documento
	Ferias            []Ferias
	Faltas            []Falta
	FeriasDisponiveis int
	Pagamentos        []Pagamento
	Vales             []Vale
}

func NewFuncionario(
	Nome, RG, CPF, PIS, CTPF, Endereco, Contato, ContatoEmergencia, Cargo string,
	Nascimento, Admissao time.Time,
	SalarioInicial float64,
) *Funcionario {
	f := new(Funcionario)
	f.Nome = Nome
	f.RG = RG
	f.CPF = CPF
	f.PIS = PIS
	f.CPF = CPF
	f.Endereco = Endereco
	f.Contato = Contato
	f.ContatoEmergencia = ContatoEmergencia
	f.Cargo = Cargo
	f.Nascimento = Nascimento
	f.Admissao = Admissao
	f.SalarioInicial = SalarioInicial
	f.Salarios = make([]Salario, 0)
	f.Ferias = make([]Ferias, 0)
	return f
}
