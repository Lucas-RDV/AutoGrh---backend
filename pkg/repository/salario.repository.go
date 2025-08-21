package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/ptrToNullTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateSalario cria um salário para um funcionário
func CreateSalario(s *entity.Salario) error {
	query := `INSERT INTO salario (funcionarioID, inicio, valor) VALUES (?, ?, ?)`

	result, err := DB.Exec(query, s.FuncionarioID, s.Inicio, s.Valor)
	if err != nil {
		return fmt.Errorf("erro ao inserir salário: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do salário: %w", err)
	}
	s.ID = id
	return nil
}

// GetSalariosByFuncionarioID retorna todos os salários de um funcionário (ordenados por início)
func GetSalariosByFuncionarioID(funcionarioID int64) ([]entity.Salario, error) {
	query := `SELECT salarioID, funcionarioID, inicio, fim, valor
	          FROM salario WHERE funcionarioID = ? ORDER BY inicio ASC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salários: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetSalariosByFuncionarioID: %v", cerr)
		}
	}()

	var salarios []entity.Salario
	for rows.Next() {
		var s entity.Salario
		var inicioStr string
		var fimStr sql.NullString

		if err := rows.Scan(&s.ID, &s.FuncionarioID, &inicioStr, &fimStr, &s.Valor); err != nil {
			log.Printf("erro ao ler salário: %v", err)
			continue
		}

		parsedInicio, err := dateStringToTime.DateStringToTime(inicioStr)
		if err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		s.Inicio = parsedInicio

		if fimStr.Valid {
			parsedFim, err := dateStringToTime.DateStringToTime(fimStr.String)
			if err == nil {
				s.Fim = &parsedFim
			}
		}

		salarios = append(salarios, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar salários: %w", err)
	}
	return salarios, nil
}

// UpdateSalario atualiza um salário
func UpdateSalario(s *entity.Salario) error {
	query := `UPDATE salario SET valor = ?, inicio = ?, fim = ? WHERE salarioID = ?`
	_, err := DB.Exec(query, s.Valor, s.Inicio, ptrToNullTime.PtrToNullTime(s.Fim), s.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar salário: %w", err)
	}
	return nil
}

// DeleteSalario remove um salário
func DeleteSalario(id int64) error {
	query := `DELETE FROM salario WHERE salarioID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar salário: %w", err)
	}
	return nil
}
