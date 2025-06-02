package repository

import (
	"AutoGRH/pkg/entity"
	"fmt"
)

// Cria um salário para um funcionário
func CreateSalario(funcionarioID int64, s *entity.Salario) error {
	query := `INSERT INTO salario (funcionarioID, valor, ano)
              VALUES (?, ?, ?)`

	result, err := DB.Exec(query, funcionarioID, s.Valor, s.Ano)
	if err != nil {
		return fmt.Errorf("erro ao inserir salário: %w", err)
	}

	s.Id, err = result.LastInsertId()
	return err
}

// Retorna todos os salários de um funcionário
func GetSalariosByFuncionarioID(funcionarioID int64) ([]entity.Salario, error) {
	query := `SELECT salarioID, valor, ano FROM salario WHERE funcionarioID = ? ORDER BY ano ASC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salários: %w", err)
	}
	defer rows.Close()

	var salarios []entity.Salario
	for rows.Next() {
		var s entity.Salario
		err := rows.Scan(&s.Id, &s.Valor, &s.Ano)
		if err != nil {
			continue
		}
		salarios = append(salarios, s)
	}
	return salarios, nil
}

// Atualiza um salário
func UpdateSalario(s *entity.Salario) error {
	query := `UPDATE salario SET valor = ?, ano = ? WHERE salarioID = ?`
	_, err := DB.Exec(query, s.Valor, s.Ano, s.Id)
	return err
}

// Deleta um salário
func DeleteSalario(id int64) error {
	query := `DELETE FROM salario WHERE salarioID = ?`
	_, err := DB.Exec(query, id)
	return err
}
