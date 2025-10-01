package entity

// Aviso representa um aviso gerado automaticamente pelo sistema
// para alertar sobre prazos próximos de vencimento (ex.: férias não pagas).
type Aviso struct {
	ID         int64  `json:"id"`
	Tipo       string `json:"tipo"`      // ferias, descanso
	Descricao  string `json:"descricao"` // mensagem descritiva
	DataEvento string `json:"dataEvento"`
	Ativo      bool   `json:"ativo"`
}

// NewAviso cria uma nova instância de Aviso já ativa.
// Exemplo de uso: NewAviso("ferias", "Férias do João vencem em 10/10/2025", "2025-10-10")
func NewAviso(tipo, descricao, dataEvento string) *Aviso {
	return &Aviso{
		Tipo:       tipo,
		Descricao:  descricao,
		DataEvento: dataEvento,
		Ativo:      true,
	}
}
