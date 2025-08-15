package repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateFerias cria um novo período de férias
func CreateFerias(f *Entity.Ferias) error {
	query := `INSERT INTO ferias (funcionarioID, dias, inicio, vencimento, vencido, valor)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, f.FuncionarioID, f.Dias, f.Inicio, f.Vencimento, f.Vencido, f.Valor)
	if err != nil {
		return fmt.Errorf("erro ao inserir férias: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID de férias inseridas: %w", err)
	}
	f.ID = id
	return nil
}

// GetFeriasByID busca férias por ID, incluindo os descansos
func GetFeriasByID(id int64) (*Entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor
	          FROM ferias WHERE feriasID = ?`

	row := DB.QueryRow(query, id)

	var f Entity.Ferias
	var inicioStr, vencimentoStr string
	if err := row.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr, &f.Vencido, &f.Valor); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar férias: %w", err)
	}

	var err error
	f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de início: %w", err)
	}
	f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de vencimento: %w", err)
	}

	// Carrega os descansos
	descansos, err := GetDescansosByFeriasID(f.ID)
	if err != nil {
		log.Printf("erro ao carregar descansos: %v", err)
	} else {
		f.Descansos = make([]Entity.Descanso, 0, len(descansos))
		for _, d := range descansos {
			f.Descansos = append(f.Descansos, *d)
		}
	}

	return &f, nil
}

// UpdateFerias atualiza um período de férias
func UpdateFerias(f *Entity.Ferias) error {
	query := `UPDATE ferias SET dias = ?, inicio = ?, vencimento = ?, vencido = ?, valor = ?
	          WHERE feriasID = ?`

	_, err := DB.Exec(query, f.Dias, f.Inicio, f.Vencimento, f.Vencido, f.Valor, f.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}
	return nil
}

// DeleteFerias remove um período de férias por ID
func DeleteFerias(id int64) error {
	query := `DELETE FROM ferias WHERE feriasID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar férias: %w", err)
	}
	return nil
}

// GetFeriasByFuncionarioID lista todas as férias de um funcionário (com descansos)
func GetFeriasByFuncionarioID(funcionarioID int64) ([]*Entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor
	          FROM ferias WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar férias por funcionário: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetFeriasByFuncionarioID: %v", cerr)
		}
	}()

	var lista []*Entity.Ferias
	for rows.Next() {
		var f Entity.Ferias
		var inicioStr, vencimentoStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr, &f.Vencido, &f.Valor); err != nil {
			log.Printf("erro ao ler férias: %v", err)
			continue
		}

		if f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr); err != nil {
			log.Printf("erro ao converter data de vencimento: %v", err)
			continue
		}

		// Carrega descansos para cada registro
		descansos, derr := GetDescansosByFeriasID(f.ID)
		if derr != nil {
			log.Printf("erro ao carregar descansos: %v", derr)
		} else {
			for _, d := range descansos {
				f.Descansos = append(f.Descansos, *d)
			}
		}

		lista = append(lista, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar férias: %w", err)
	}
	return lista, nil
}

// ListFerias lista todos os registros de férias (sem carregar descansos)
func ListFerias() ([]Entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor FROM ferias`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar férias: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListFerias: %v", cerr)
		}
	}()

	var lista []Entity.Ferias
	for rows.Next() {
		var f Entity.Ferias
		var inicioStr, vencimentoStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr, &f.Vencido, &f.Valor); err != nil {
			log.Printf("erro ao ler férias: %v", err)
			continue
		}

		if f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr); err != nil {
			log.Printf("erro ao converter data de vencimento: %v", err)
			continue
		}

		lista = append(lista, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar férias: %w", err)
	}
	return lista, nil
}
