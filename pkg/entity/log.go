package entity

import "time"

// Log representa um registro de ação ou evento gerado por um usuário no sistema
// Inclui a referência ao usuário, o tipo de evento, mensagem e data de ocorrência

type Log struct {
	ID        int64     `json:"id"`
	UsuarioID int64     `json:"usuario_id"`
	EventoID  int64     `json:"evento_id"`
	Message   string    `json:"message"`
	Data      time.Time `json:"data"`
}

// NewLog cria uma nova instância de Log com data atual
func NewLog(usuarioID, eventoID int64, message string) *Log {
	return &Log{
		UsuarioID: usuarioID,
		EventoID:  eventoID,
		Message:   message,
		Data:      time.Now(),
	}
}
