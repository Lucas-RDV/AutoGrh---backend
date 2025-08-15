package repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateDescanso cria um descanso vinculado a um período de férias
func CreateDescanso(d *Entity.Descanso) error {
	query := `INSERT INTO descanso (feriasID, inicio, fim, valor, pago, aprovado)
              VALUES (?, ?, ?, ?, ?, ?)`

	res, err := DB.Exec(query, d.FeriasID, d.Inicio, d.Fim, d.Valor, d.Pago, d.Aprovado)
	if err != nil {
		return fmt.Errorf("erro ao inserir descanso: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do descanso inserido: %w", err)
	}
	d.ID = id
	return nil
}

// GetDescansoByID busca um descanso por ID
func GetDescansoByID(id int64) (*Entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
	          FROM descanso WHERE descansoID = ?`

	row := DB.QueryRow(query, id)

	var d Entity.Descanso
	var inicioStr, fimStr string

	err := row.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar descanso: %w", err)
	}
	d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de início: %w", err)
	}
	d.Fim, err = dateStringToTime.DateStringToTime(fimStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de fim: %w", err)
	}

	return &d, nil
}

// GetDescansosByFeriasID busca todos os descansos de um período de férias
func GetDescansosByFeriasID(feriasID int64) ([]*Entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
	          FROM descanso WHERE feriasID = ?`

	rows, err := DB.Query(query, feriasID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar descansos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetDescansosByFeriasID: %v", cerr)
		}
	}()

	var descansos []*Entity.Descanso
	for rows.Next() {
		var d Entity.Descanso
		var inicioStr, fimStr string

		err := rows.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
		if err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr)
		if err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		d.Fim, err = dateStringToTime.DateStringToTime(fimStr)
		if err != nil {
			log.Printf("erro ao converter data de fim: %v", err)
			continue
		}

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// ListDescansos lista todos os descansos
func ListDescansos() ([]*Entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado FROM descanso`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar descansos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListDescansos: %v", cerr)
		}
	}()

	var descansos []*Entity.Descanso
	for rows.Next() {
		var d Entity.Descanso
		var inicioStr, fimStr string

		err := rows.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
		if err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr)
		if err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		d.Fim, err = dateStringToTime.DateStringToTime(fimStr)
		if err != nil {
			log.Printf("erro ao converter data de fim: %v", err)
			continue
		}

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// UpdateDescanso atualiza um descanso
func UpdateDescanso(d *Entity.Descanso) error {
	query := `UPDATE descanso SET inicio = ?, fim = ?, valor = ?, pago = ?, aprovado = ? 
	          WHERE descansoID = ?`

	_, err := DB.Exec(query, d.Inicio, d.Fim, d.Valor, d.Pago, d.Aprovado, d.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar descanso: %w", err)
	}
	return nil
}

// DeleteDescanso deleta um descanso
func DeleteDescanso(id int64) error {
	query := `DELETE FROM descanso WHERE descansoID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar descanso: %w", err)
	}
	return nil
}
