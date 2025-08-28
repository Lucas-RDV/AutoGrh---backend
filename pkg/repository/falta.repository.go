package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
)

// CreateFalta cria um registro de falta
func CreateFalta(f *entity.Falta) error {
	query := `INSERT INTO falta (funcionarioID, quantidade, data) VALUES (?, ?, ?)`

	result, err := DB.Exec(query, f.FuncionarioID, f.Quantidade, f.Mes)
	if err != nil {
		return fmt.Errorf("erro ao inserir falta: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da falta inserida: %w", err)
	}
	f.ID = id
	return nil
}

// GetFaltasByFuncionarioID busca todas as faltas de um funcionário
func GetFaltasByFuncionarioID(funcionarioID int64) ([]*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas por funcionário: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Falta
	for rows.Next() {
		var f entity.Falta
		var dataStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
			return nil, fmt.Errorf("erro ao ler falta: %w", err)
		}

		f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data: %w", err)
		}

		lista = append(lista, &f)
	}
	return lista, nil
}

// GetFaltaByID retorna uma falta pelo ID
func GetFaltaByID(id int64) (*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE faltaID = ?`
	row := DB.QueryRow(query, id)

	var f entity.Falta
	var dataStr string
	if err := row.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar falta: %w", err)
	}

	var err error
	f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data: %w", err)
	}

	return &f, nil
}

// UpdateFalta atualiza um registro de falta
func UpdateFalta(f *entity.Falta) error {
	query := `UPDATE falta SET quantidade = ?, data = ? WHERE faltaID = ?`
	_, err := DB.Exec(query, f.Quantidade, f.Mes, f.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar falta: %w", err)
	}
	return nil
}

// DeleteFalta remove uma falta por ID
func DeleteFalta(id int64) error {
	query := `DELETE FROM falta WHERE faltaID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar falta: %w", err)
	}
	return nil
}

// ListAllFaltas retorna todas as faltas cadastradas
func ListAllFaltas() ([]*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Falta
	for rows.Next() {
		var f entity.Falta
		var dataStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
			return nil, fmt.Errorf("erro ao ler falta: %w", err)
		}

		f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data: %w", err)
		}

		lista = append(lista, &f)
	}
	return lista, nil
}
