package entity

import "time"

type Log struct {
	Id      string
	Message string
	Data    time.Time
}

func NewLog(Message string) *Log {
	l := new(Log)
	l.Message = Message
	l.Data = time.Now()
	return l
}
