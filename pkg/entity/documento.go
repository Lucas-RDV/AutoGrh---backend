package entity

// Documento representa um arquivo/documento associado a um funcionário
// Pode ser um comprovante, contrato, atestado, etc.
type Documento struct {
	ID            int64  `json:"id"`
	FuncionarioID int64  `json:"funcionarioID"`
	Caminho       string `json:"caminho"`
}

// NewDocumento cria uma nova instância de Documento vinculado a um funcionário
func NewDocumento(caminho string, funcionarioID int64) *Documento {
	return &Documento{
		Caminho:       caminho,
		FuncionarioID: funcionarioID,
	}
}
