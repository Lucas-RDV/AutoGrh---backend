package Repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/DateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateLog cria um novo log no banco
func CreateLog(l *Entity.Log) error {
	query := `INSERT INTO log (usuarioID, eventoID, data, action) VALUES (?, ?, ?, ?)`

	// l.Data é time.Time; o driver MySQL mapeia para TIMESTAMP corretamente
	result, err := DB.Exec(query, l.UsuarioID, l.EventoID, l.Data, l.Message)
	if err != nil {
		return fmt.Errorf("erro ao inserir log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do novo log: %w", err)
	}
	l.ID = id
	return nil
}

// GetLogByID busca um log por ID
func GetLogByID(id int64) (*Entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log WHERE logID = ?`
	row := DB.QueryRow(query, id)

	var l Entity.Log
	var dtStr string
	if err := row.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar log: %w", err)
	}
	t, err := DateStringToTime.DateStringToTime(dtStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data do log: %w", err)
	}
	l.Data = t
	return &l, nil
}

// GetLogsByUsuarioID lista todos os logs de um usuário (mais recentes primeiro)
func GetLogsByUsuarioID(usuarioID int64) ([]*Entity.Log, error) {
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

	var logs []*Entity.Log
	for rows.Next() {
		var l Entity.Log
		var dtStr string
		if err := rows.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}
		t, err := DateStringToTime.DateStringToTime(dtStr)
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
func ListAllLogs(limit int) ([]*Entity.Log, error) {
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

	var logs []*Entity.Log
	for rows.Next() {
		var l Entity.Log
		var dtStr string
		if err := rows.Scan(&l.ID, &l.UsuarioID, &l.EventoID, &dtStr, &l.Message); err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}
		t, err := DateStringToTime.DateStringToTime(dtStr)
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
