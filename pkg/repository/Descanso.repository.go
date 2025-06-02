package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateDescanso Cria um descanso vinculado a um período de férias
func CreateDescanso(d *entity.Descanso) error {
	query := `INSERT INTO descanso (feriasID, inicio, fim, valor, pago, aprovado)
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, d.FeriasID, d.Inicio, d.Fim, d.Valor, d.Pago, d.Aprovado)
	if err != nil {
		return fmt.Errorf("erro ao inserir descanso: %w", err)
	}

	d.Id, err = result.LastInsertId()
	return err
}

// GetDescansoByID Busca um descanso por ID
func GetDescansoByID(id int64) (*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
              FROM descanso WHERE descansoID = ?`

	row := DB.QueryRow(query, id)

	var d entity.Descanso
	var inicioStr, fimStr string

	err := row.Scan(&d.Id, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar descanso: %w", err)
	}
	d.Inicio, _ = time.Parse("2006-01-02", inicioStr)
	d.Fim, _ = time.Parse("2006-01-02", fimStr)

	return &d, nil
}

// GetDescansosByFeriasID Busca todos os descansos de um período de férias
func GetDescansosByFeriasID(feriasID int64) ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
              FROM descanso WHERE feriasID = ?`

	rows, err := DB.Query(query, feriasID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar descansos: %w", err)
	}
	defer rows.Close()

	var descansos []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		err := rows.Scan(&d.Id, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
		if err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		d.Inicio, _ = time.Parse("2006-01-02", inicioStr)
		d.Fim, _ = time.Parse("2006-01-02", fimStr)

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// ListDescansos Lista todos os descansos
func ListDescansos() ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado FROM descanso`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar descansos: %w", err)
	}
	defer rows.Close()

	var descansos []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		err := rows.Scan(&d.Id, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado)
		if err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		d.Inicio, _ = time.Parse("2006-01-02", inicioStr)
		d.Fim, _ = time.Parse("2006-01-02", fimStr)

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// UpdateDescanso Atualiza um descanso
func UpdateDescanso(d *entity.Descanso) error {
	query := `UPDATE descanso SET inicio = ?, fim = ?, valor = ?, pago = ?, aprovado = ? 
              WHERE descansoID = ?`

	_, err := DB.Exec(query, d.Inicio, d.Fim, d.Valor, d.Pago, d.Aprovado, d.Id)
	return err
}

// DeleteDescanso Deleta um descanso
func DeleteDescanso(id int64) error {
	query := `DELETE FROM descanso WHERE descansoID = ?`
	_, err := DB.Exec(query, id)
	return err
}
