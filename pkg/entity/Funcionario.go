package entity

import "time"

type Funcionario struct {
	Id                int64
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
	Demissao          *time.Time
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
	f.CTPF = CTPF
	f.Endereco = Endereco
	f.Contato = Contato
	f.ContatoEmergencia = ContatoEmergencia
	f.Cargo = Cargo
	f.Nascimento = Nascimento
	f.Admissao = Admissao
	f.Demissao = nil
	f.SalarioInicial = SalarioInicial
	f.Salarios = make([]Salario, 0)
	f.Documentos = make([]Documento, 0)
	f.Ferias = make([]Ferias, 0)
	f.Faltas = make([]Falta, 0)
	f.Pagamentos = make([]Pagamento, 0)
	f.Vales = make([]Vale, 0)
	return f
}
