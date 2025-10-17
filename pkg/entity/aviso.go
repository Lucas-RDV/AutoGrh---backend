package entity

import "time"

type Aviso struct {
	ID           int64     `json:"id"`
	Tipo         string    `json:"tipo"` // "FERIAS_VENCENDO" | "FERIAS_VENCIDAS" | "VALE_PENDENTE" | "DESCANSO_PENDENTE"
	Mensagem     string    `json:"mensagem"`
	ReferenciaID *int64    `json:"referencia_id"` // opcional (ex.: feriasID, valeID, descansoID)
	CriadoEm     time.Time `json:"criado_em"`
	Ativo        bool      `json:"ativo"`
}
