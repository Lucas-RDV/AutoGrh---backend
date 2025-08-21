package entity

// Documento representa um arquivo/documento associado a um funcionário
// Pode ser um comprovante, contrato, atestado, etc.
type Documento struct {
	ID            int64  `json:"id"`
	Doc           []byte `json:"doc"`
	FuncionarioID int64  `json:"funcionario_id"`
}

// NewDocumento cria uma nova instância de Documento vinculado a um funcionário
func NewDocumento(doc []byte, funcionarioID int64) *Documento {
	return &Documento{
		Doc:           doc,
		FuncionarioID: funcionarioID,
	}
}
