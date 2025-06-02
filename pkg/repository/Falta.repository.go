package repository

import (
	"AutoGRH/pkg/entity"
	"fmt"
	"log"
	"time"
)

// Cria uma falta
func CreateFalta(f *entity.Falta) error {
	query := `INSERT INTO falta (funcionarioID, quantidade, data) VALUES (?, ?, ?)`

	result, err := DB.Exec(query, f.FuncionarioId, f.Quantidade, f.Mes)
	if err != nil {
		return fmt.Errorf("erro ao inserir falta: %w", err)
	}

	f.Id, err = result.LastInsertId()
	return err
}

// Busca todas as faltas de um funcionário
func GetFaltasByFuncionarioID(funcionarioId int64) ([]entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioId)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar faltas: %w", err)
	}
	defer rows.Close()

	var faltas []entity.Falta
	for rows.Next() {
		var f entity.Falta
		var mesStr string

		err := rows.Scan(&f.Id, &f.FuncionarioId, &f.Quantidade, &mesStr)
		if err != nil {
			log.Printf("erro ao ler falta: %v", err)
			continue
		}

		parsed, err := time.Parse("2006-01-02", mesStr)
		if err != nil {
			log.Printf("erro ao converter data: %v", err)
			continue
		}
		f.Mes = parsed

		faltas = append(faltas, f)
	}
	return faltas, nil
}

// Atualiza uma falta
func UpdateFalta(f *entity.Falta) error {
	query := `UPDATE falta SET quantidade = ?, data = ? WHERE faltaID = ?`
	_, err := DB.Exec(query, f.Quantidade, f.Mes, f.Id)
	return err
}

// Deleta uma falta
func DeleteFalta(id int64) error {
	query := `DELETE FROM falta WHERE faltaID = ?`
	_, err := DB.Exec(query, id)
	return err
}

func ListFaltas() ([]entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas: %w", err)
	}
	defer rows.Close()

	var faltas []entity.Falta
	for rows.Next() {
		var f entity.Falta
		var mesStr string

		err := rows.Scan(&f.Id, &f.FuncionarioId, &f.Quantidade, &mesStr)
		if err != nil {
			log.Printf("erro ao ler falta: %v", err)
			continue
		}

		parsed, err := time.Parse("2006-01-02", mesStr)
		if err != nil {
			log.Printf("erro ao converter data: %v", err)
			continue
		}
		f.Mes = parsed

		faltas = append(faltas, f)
	}
	return faltas, nil
}
