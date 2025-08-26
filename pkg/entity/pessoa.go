package entity

// Pessoa representa uma pessoa física única, independente do vínculo empregatício

type Pessoa struct {
	ID                int64  `json:"id"`
	Nome              string `json:"nome"`
	CPF               string `json:"cpf"`
	RG                string `json:"rg"`
	Endereco          string `json:"endereco"`
	Contato           string `json:"contato"`
	ContatoEmergencia string `json:"contato_emergencia"`
}

func NewPessoa(nome, cpf, rg, endereco, contato, contatoEmergencia string) *Pessoa {
	return &Pessoa{
		Nome:              nome,
		CPF:               cpf,
		RG:                rg,
		Endereco:          endereco,
		Contato:           contato,
		ContatoEmergencia: contatoEmergencia,
	}
}
