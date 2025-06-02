package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateFerias Cria um novo período de férias
func CreateFerias(f *entity.Ferias) error {
	query := `INSERT INTO ferias (funcionarioID, dias, inicio, vencimento, vencido, valor)
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, f.FuncionarioID, f.Dias, f.Inicio, f.Vencimento, f.Vencido, f.Valor)
	if err != nil {
		return fmt.Errorf("erro ao inserir ferias: %w", err)
	}

	f.Id, err = result.LastInsertId()
	return err
}

// GetFeriasByID Busca férias por ID, incluindo os descansos
func GetFeriasByID(id int64) (*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor
              FROM ferias WHERE feriasID = ?`

	row := DB.QueryRow(query, id)

	var f entity.Ferias
	err := row.Scan(&f.Id, &f.FuncionarioID, &f.Dias, &f.Inicio, &f.Vencimento, &f.Vencido, &f.Valor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar ferias: %w", err)
	}

	// Carrega os descansos
	descansos, err := GetDescansosByFeriasID(f.Id)
	if err != nil {
		log.Printf("erro ao carregar descansos: %v", err)
	}
	f.Descansos = []entity.Descanso{}
	for _, d := range descansos {
		f.Descansos = append(f.Descansos, *d)
	}

	return &f, nil
}

// UpdateFerias Atualiza férias
func UpdateFerias(f *entity.Ferias) error {
	query := `UPDATE ferias SET dias = ?, inicio = ?, vencimento = ?, vencido = ?, valor = ?
              WHERE feriasID = ?`

	_, err := DB.Exec(query, f.Dias, f.Inicio, f.Vencimento, f.Vencido, f.Valor, f.Id)
	return err
}

// DeleteFerias Deleta período de férias
func DeleteFerias(id int64) error {
	query := `DELETE FROM ferias WHERE feriasID = ?`
	_, err := DB.Exec(query, id)
	return err
}

// GetFeriasByFuncionarioID Lista todas as férias de um funcionário
func GetFeriasByFuncionarioID(funcionarioID int64) ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor
              FROM ferias WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ferias por funcionario: %w", err)
	}
	defer rows.Close()

	var feriasList []*entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		var inicioStr, vencimentoStr string

		err := rows.Scan(&f.Id, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr, &f.Vencido, &f.Valor)
		if err != nil {
			log.Printf("erro ao ler ferias: %v", err)
			continue
		}

		// Parse das datas
		f.Inicio, err = time.Parse("2006-01-02", inicioStr)
		if err != nil {
			log.Printf("erro ao converter data de inicio: %v", err)
			continue
		}
		f.Vencimento, err = time.Parse("2006-01-02", vencimentoStr)
		if err != nil {
			log.Printf("erro ao converter data de vencimento: %v", err)
			continue
		}

		// Carrega descansos
		descansos, _ := GetDescansosByFeriasID(f.Id)
		for _, d := range descansos {
			f.Descansos = append(f.Descansos, *d)
		}

		feriasList = append(feriasList, &f)
	}
	return feriasList, nil
}

func ListFerias() ([]entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor FROM ferias`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar férias: %w", err)
	}
	defer rows.Close()

	var lista []entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		var inicioStr, vencimentoStr string

		err := rows.Scan(&f.Id, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr, &f.Vencido, &f.Valor)
		if err != nil {
			log.Printf("erro ao ler férias: %v", err)
			continue
		}

		f.Inicio, err = time.Parse("2006-01-02", inicioStr)
		if err != nil {
			log.Printf("erro ao converter data inicio: %v", err)
			continue
		}
		f.Vencimento, err = time.Parse("2006-01-02", vencimentoStr)
		if err != nil {
			log.Printf("erro ao converter data vencimento: %v", err)
			continue
		}

		lista = append(lista, f)
	}
	return lista, nil
}
