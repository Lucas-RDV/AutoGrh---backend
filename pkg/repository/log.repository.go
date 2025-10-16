package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type LogView struct {
	ID        int64     `json:"id"`
	UsuarioID int64     `json:"usuario_id"`
	EventoID  int64     `json:"evento_id"`
	Evento    string    `json:"evento"`  // nome do evento (tabela 'evento', coluna 'tipo')
	Message   string    `json:"message"` // vem de l.action
	Data      time.Time `json:"data"`
}

func ListAllLogsView(limit int) ([]*LogView, error) {
	if DB == nil {
		return nil, errors.New("DB não inicializado")
	}
	if limit <= 0 {
		limit = 200
	}

	const q = `
		SELECT
			l.logID,
			l.usuarioID,
			l.eventoID,
			e.tipo      AS evento_nome,
			l.action    AS message,
			l.data
		FROM log l
		JOIN evento e ON e.eventoID = l.eventoID
		ORDER BY l.data DESC, l.logID DESC
		LIMIT ?
	`

	rows, err := DB.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("ListAllLogsView: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListAllLogsView: %v", cerr)
		}
	}()

	var out []*LogView
	for rows.Next() {
		var lv LogView
		var dtStr string

		if err := rows.Scan(
			&lv.ID,
			&lv.UsuarioID,
			&lv.EventoID,
			&lv.Evento,  // e.tipo
			&lv.Message, // l.action
			&dtStr,      // l.data como string
		); err != nil {
			log.Printf("erro ao ler log (view): %v", err)
			continue
		}

		t, err := dateStringToTime.DateStringToTime(dtStr)
		if err != nil {
			log.Printf("erro ao converter data do log (view): %v", err)
			continue
		}
		lv.Data = t

		out = append(out, &lv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar logs (view): %w", err)
	}
	return out, nil
}

// CreateLog cria um novo log no banco
func CreateLog(l *entity.Log) error {
	if DB == nil {
		return errors.New("DB não inicializado")
	}

	// Normaliza a data para o fuso do servidor (ou um fuso fixo, se preferir)
	t := l.Data
	if t.IsZero() {
		t = time.Now()
	}
	// Se quiser cravar o fuso da aplicação:
	// loc, _ := time.LoadLocation("America/Campo_Grande")
	// t = t.In(loc)
	t = t.In(time.Local)

	// Grava como string "YYYY-MM-DD HH:MM:SS" para evitar ambiguidade de driver
	const layout = "2006-01-02 15:04:05"
	_, err := DB.Exec(`
		INSERT INTO log (usuarioID, eventoID, data, action)
		VALUES (?, ?, ?, ?)`,
		l.UsuarioID, l.EventoID, t.Format(layout), l.Message,
	)
	if err != nil {
		return fmt.Errorf("CreateLog: %w", err)
	}
	return nil
}

// GetLogByID busca um log por ID
func GetLogByID(id int64) (*entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log WHERE logID = ?`
	row := DB.QueryRow(query, id)

	var l entity.Log
	var dtStr string
	if err := row.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar log: %w", err)
	}
	t, err := dateStringToTime.DateStringToTime(dtStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data do log: %w", err)
	}
	l.Data = t
	return &l, nil
}

// GetLogsByUsuarioID lista todos os logs de um usuário (mais recentes primeiro)
func GetLogsByUsuarioID(usuarioID int64) ([]*entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log WHERE usuarioID = ? ORDER BY data DESC`

	rows, err := DB.Query(query, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar logs do usuário: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetLogsByUsuarioID: %v", cerr)
		}
	}()

	var logs []*entity.Log
	for rows.Next() {
		var l entity.Log
		var dtStr string
		if err := rows.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}
		t, err := dateStringToTime.DateStringToTime(dtStr)
		if err != nil {
			log.Printf("erro ao converter data do log: %v", err)
			continue
		}
		l.Data = t
		logs = append(logs, &l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar logs do usuário: %w", err)
	}
	return logs, nil
}

// ListAllLogs lista todos os logs do sistema (com limite obrigatório > 0)
func ListAllLogs(limit int) ([]*entity.Log, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit deve ser maior que zero")
	}

	query := `SELECT logID, usuarioID, eventoID, data, action FROM log ORDER BY data DESC LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar logs: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListAllLogs: %v", cerr)
		}
	}()

	var logs []*entity.Log
	for rows.Next() {
		var l entity.Log
		var dtStr string
		if err := rows.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}
		t, err := dateStringToTime.DateStringToTime(dtStr)
		if err != nil {
			log.Printf("erro ao converter data do log: %v", err)
			continue
		}
		l.Data = t
		logs = append(logs, &l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar logs: %w", err)
	}
	return logs, nil
}
