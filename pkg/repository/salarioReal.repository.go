package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/ptrToNullTime"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateSalarioReal insere um novo salário real no banco de dados
func CreateSalarioReal(s *entity.SalarioReal) error {
	query := `INSERT INTO salario_real (funcionarioID, inicio, fim, valor) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query,
		s.FuncionarioID,
		s.Inicio,
		ptrToNullTime.PtrToNullTime(s.Fim),
		s.Valor,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir salário real: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do novo salário real: %w", err)
	}
	s.ID = id
	return nil
}

// GetSalariosReaisByFuncionarioID retorna o histórico de salários reais de um funcionário
func GetSalariosReaisByFuncionarioID(funcionarioID int64) ([]*entity.SalarioReal, error) {
	query := `
		SELECT salarioRealID, funcionarioID, inicio, fim, valor
		FROM salario_real
		WHERE funcionarioID = ?
		ORDER BY inicio DESC
	`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar salários reais: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetSalariosReaisByFuncionarioID: %v", cerr)
		}
	}()

	var lista []*entity.SalarioReal
	for rows.Next() {
		var s entity.SalarioReal
		var inicioStr string
		var fimStr sql.NullString

		if err := rows.Scan(&s.ID, &s.FuncionarioID, &inicioStr, &fimStr, &s.Valor); err != nil {
			return nil, fmt.Errorf("erro ao ler salário real: %w", err)
		}

		tInicio, err := dateStringToTime.DateStringToTime(inicioStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter inicio: %w", err)
		}
		s.Inicio = tInicio

		if fimStr.Valid && fimStr.String != "" && fimStr.String != "0000-00-00" && fimStr.String != "0000-00-00 00:00:00" {
			tFim, err := dateStringToTime.DateStringToTime(fimStr.String)
			if err != nil {
				return nil, fmt.Errorf("erro ao converter fim: %w", err)
			}
			s.Fim = &tFim
		} else {
			s.Fim = nil
		}

		lista = append(lista, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar salários reais: %w", err)
	}
	return lista, nil
}

// GetSalarioRealAtual retorna o salário real ativo de um funcionário (sem fim definido)
func GetSalarioRealAtual(funcionarioID int64) (*entity.SalarioReal, error) {
	const q = `
		SELECT salarioRealID, funcionarioID, inicio, fim, valor
		FROM salario_real
		WHERE funcionarioID = ?
		  AND fim IS NULL 
		ORDER BY inicio DESC
		LIMIT 1`
	row := DB.QueryRow(q, funcionarioID)

	var s entity.SalarioReal
	var inicioStr string
	var fimStr sql.NullString

	if err := row.Scan(&s.ID, &s.FuncionarioID, &inicioStr, &fimStr, &s.Valor); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar salário real atual: %w", err)
	}

	tInicio, err := dateStringToTime.DateStringToTime(inicioStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter inicio: %w", err)
	}
	s.Inicio = tInicio

	if fimStr.Valid && fimStr.String != "" && fimStr.String != "0000-00-00 00:00:00" {
		tFim, err := dateStringToTime.DateStringToTime(fimStr.String)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter fim: %w", err)
		}
		s.Fim = &tFim
	} else {
		s.Fim = nil
	}
	return &s, nil
}

// UpdateSalarioReal atualiza um salário real existente
func UpdateSalarioReal(s *entity.SalarioReal) error {
	query := `UPDATE salario_real SET inicio = ?, fim = ?, valor = ? WHERE salarioRealID = ?`
	_, err := DB.Exec(query,
		s.Inicio,
		ptrToNullTime.PtrToNullTime(s.Fim),
		s.Valor,
		s.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar salário real: %w", err)
	}
	return nil
}

// DeleteSalarioReal remove um salário real do banco
func DeleteSalarioReal(id int64) error {
	query := `DELETE FROM salario_real WHERE salarioRealID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar salário real: %w", err)
	}
	return nil
}

func EncerrarSalarioReal(id int64, fim time.Time) error {
	nt := ptrToNullTime.PtrToNullTime(&fim)
	_, err := DB.Exec(`UPDATE salario_real SET fim=? WHERE salarioRealID=?`, nt, id)
	return err
}
