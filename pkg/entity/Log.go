package entity

import "time"

type Log struct {
	Id        int64
	UsuarioId int64
	EventoId  int
	Message   string
	Data      time.Time
}

func NewLog(usuarioId int64, eventoId int, message string) *Log {
	l := new(Log)
	l.Message = message
	l.UsuarioId = usuarioId
	l.EventoId = eventoId
	l.Data = time.Now()
	return l
}
