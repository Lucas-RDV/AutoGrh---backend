package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Cria um novo log no banco
func CreateLog(l *entity.Log) error {
	query := `INSERT INTO log (usuarioID, eventoID, data, action)
			  VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, l.UsuarioId, l.EventoId, l.Data, l.Message)
	if err != nil {
		return fmt.Errorf("erro ao inserir log: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		l.Id = id
	}
	return err
}

// Busca um log por ID
func GetLogByID(id int64) (*entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log WHERE logID = ?`
	row := DB.QueryRow(query, id)

	var l entity.Log
	var dataStr string
	err := row.Scan(&l.Id, &l.UsuarioId, &l.EventoId, &dataStr, &l.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar log: %w", err)
	}

	l.Data, err = time.Parse("2006-01-02 15:04:05", dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data do log: %w", err)
	}

	return &l, nil
}

// Lista todos os logs de um usuário
func GetLogsByUsuarioID(usuarioID int64) ([]*entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log WHERE usuarioID = ? ORDER BY data DESC`

	rows, err := DB.Query(query, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar logs do usuário: %w", err)
	}
	defer rows.Close()

	var logs []*entity.Log
	for rows.Next() {
		var l entity.Log
		var dataStr string
		err := rows.Scan(&l.Id, &l.UsuarioId, &l.EventoId, &dataStr, &l.Message)
		if err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}

		l.Data, err = time.Parse("2006-01-02 15:04:05", dataStr)
		if err != nil {
			log.Printf("erro ao converter data do log: %v", err)
			continue
		}

		logs = append(logs, &l)
	}
	return logs, nil
}

// Lista todos os logs do sistema (com limite opcional)
func ListAllLogs(limit int) ([]*entity.Log, error) {
	query := `SELECT logID, usuarioID, eventoID, data, action FROM log ORDER BY data DESC LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar logs: %w", err)
	}
	defer rows.Close()

	var logs []*entity.Log
	for rows.Next() {
		var l entity.Log
		var dataStr string
		err := rows.Scan(&l.Id, &l.UsuarioId, &l.EventoId, &dataStr, &l.Message)
		if err != nil {
			log.Printf("erro ao ler log: %v", err)
			continue
		}

		l.Data, err = time.Parse("2006-01-02 15:04:05", dataStr)
		if err != nil {
			log.Printf("erro ao converter data do log: %v", err)
			continue
		}

		logs = append(logs, &l)
	}
	return logs, nil
}
