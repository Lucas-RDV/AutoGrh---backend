package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateVale cria um novo vale
func CreateVale(v *entity.Vale) error {
	query := `INSERT INTO vale (funcionarioID, valor, data, aprovado, pago) VALUES (?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, v.FuncionarioID, v.Valor, v.Data, v.Aprovado, v.Pago)
	if err != nil {
		return fmt.Errorf("erro ao criar vale: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do vale criado: %w", err)
	}
	v.ID = id
	return nil
}

// GetValeByID busca um vale por ID
func GetValeByID(id int64) (*entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago FROM vale WHERE valeID = ?`
	row := DB.QueryRow(query, id)

	var v entity.Vale
	var dataStr string
	if err := row.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar vale: %w", err)
	}

	parsed, err := dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data do vale: %w", err)
	}
	v.Data = parsed

	return &v, nil
}

// GetValesByFuncionarioID lista todos os vales de um funcionário
func GetValesByFuncionarioID(funcionarioID int64) ([]entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago FROM vale WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vales: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetValesByFuncionarioID: %v", cerr)
		}
	}()

	var vales []entity.Vale
	for rows.Next() {
		var v entity.Vale
		var dataStr string
		if err := rows.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago); err != nil {
			log.Printf("erro ao ler vale: %v", err)
			continue
		}

		parsed, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do vale: %v", err)
			continue
		}
		v.Data = parsed

		vales = append(vales, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar vales: %w", err)
	}
	return vales, nil
}

// UpdateVale atualiza um vale existente
func UpdateVale(v *entity.Vale) error {
	query := `UPDATE vale SET valor = ?, data = ?, aprovado = ?, pago = ? WHERE valeID = ?`
	_, err := DB.Exec(query, v.Valor, v.Data, v.Aprovado, v.Pago, v.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar vale: %w", err)
	}
	return nil
}

// DeleteVale remove um vale por ID
func DeleteVale(id int64) error {
	query := `DELETE FROM vale WHERE valeID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar vale: %w", err)
	}
	return nil
}
