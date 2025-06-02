package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Cria um salário para um funcionário
func CreateSalario(s *entity.Salario) error {
	query := `INSERT INTO salario (funcionarioID, inicio, valor) VALUES (?, ?, ?)`

	result, err := DB.Exec(query, s.FuncionarioID, s.Inicio.Format("2006-01-02"), s.Valor)
	if err != nil {
		return fmt.Errorf("erro ao inserir salário: %w", err)
	}

	s.Id, err = result.LastInsertId()
	return err
}

// Retorna todos os salários de um funcionário
func GetSalariosByFuncionarioID(funcionarioID int64) ([]entity.Salario, error) {
	query := `SELECT salarioID, funcionarioID, inicio, fim, valor FROM salario WHERE funcionarioID = ? ORDER BY inicio ASC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salários: %w", err)
	}
	defer rows.Close()

	var salarios []entity.Salario
	for rows.Next() {
		var s entity.Salario
		var inicioStr, fimStr sql.NullString
		err := rows.Scan(&s.Id, &s.FuncionarioID, &inicioStr, &fimStr, &s.Valor)
		if err != nil {
			log.Printf("erro ao ler salário: %v", err)
			continue
		}

		s.Inicio, err = time.Parse("2006-01-02", inicioStr.String)
		if err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}

		if fimStr.Valid {
			fimParsed, err := time.Parse("2006-01-02", fimStr.String)
			if err == nil {
				s.Fim = &fimParsed
			}
		}

		salarios = append(salarios, s)
	}
	return salarios, nil
}

// Atualiza um salário
func UpdateSalario(s *entity.Salario) error {
	query := `UPDATE salario SET valor = ?, inicio = ?, fim = ? WHERE salarioID = ?`
	var fimStr interface{}
	if s.Fim != nil {
		fimStr = s.Fim.Format("2006-01-02")
	} else {
		fimStr = nil
	}
	_, err := DB.Exec(query, s.Valor, s.Inicio.Format("2006-01-02"), fimStr, s.Id)
	return err
}

// Deleta um salário
func DeleteSalario(id int64) error {
	query := `DELETE FROM salario WHERE salarioID = ?`
	_, err := DB.Exec(query, id)
	return err
}
