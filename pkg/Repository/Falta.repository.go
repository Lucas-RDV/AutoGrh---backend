package Repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/DateStringToTime"
	"fmt"
	"log"
)

// CreateFalta cria um registro de falta
func CreateFalta(f *Entity.Falta) error {
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
func GetFaltasByFuncionarioID(funcionarioID int64) ([]Entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faltas: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetFaltasByFuncionarioID: %v", cerr)
		}
	}()

	var faltas []Entity.Falta
	for rows.Next() {
		var f Entity.Falta
		var mesStr string

		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &mesStr); err != nil {
			log.Printf("erro ao ler falta: %v", err)
			continue
		}

		parsed, err := DateStringToTime.DateStringToTime(mesStr)
		if err != nil {
			log.Printf("erro ao converter data da falta: %v", err)
			continue
		}
		f.Mes = parsed

		faltas = append(faltas, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar faltas: %w", err)
	}
	return faltas, nil
}

// UpdateFalta atualiza um registro de falta
func UpdateFalta(f *Entity.Falta) error {
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

// ListFaltas lista todas as faltas
func ListFaltas() ([]Entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListFaltas: %v", cerr)
		}
	}()

	var faltas []Entity.Falta
	for rows.Next() {
		var f Entity.Falta
		var mesStr string

		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &mesStr); err != nil {
			log.Printf("erro ao ler falta: %v", err)
			continue
		}

		parsed, err := DateStringToTime.DateStringToTime(mesStr)
		if err != nil {
			log.Printf("erro ao converter data da falta: %v", err)
			continue
		}
		f.Mes = parsed

		faltas = append(faltas, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar faltas: %w", err)
	}
	return faltas, nil
}
