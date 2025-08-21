package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateFolha cria uma nova folha de pagamento
func CreateFolha(f *entity.FolhaPagamentos) error {
	query := `INSERT INTO folha_pagamento (data) VALUES (?)`

	result, err := DB.Exec(query, f.Data)
	if err != nil {
		return fmt.Errorf("erro ao criar folha: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da folha criada: %w", err)
	}
	f.ID = id
	return nil
}

// GetFolhaByID busca uma folha de pagamento e seus pagamentos
func GetFolhaByID(id int64) (*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento WHERE folhaID = ?`
	row := DB.QueryRow(query, id)

	var f entity.FolhaPagamentos
	var dataStr string
	if err := row.Scan(&f.ID, &dataStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar folha: %w", err)
	}

	var err error
	f.Data, err = dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data da folha: %w", err)
	}

	// Buscar pagamentos relacionados e somar o valor total
	pagamentos, perr := GetPagamentosByFolhaID(f.ID)
	if perr != nil {
		log.Printf("erro ao carregar pagamentos da folha %d: %v", f.ID, perr)
	} else {
		f.Pagamentos = pagamentos
		for _, p := range pagamentos {
			f.Valor += p.Valor
		}
	}

	return &f, nil
}

// ListFolhas lista todas as folhas de pagamento (com total calculado)
func ListFolhas() ([]*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento ORDER BY data DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar folhas: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListFolhas: %v", cerr)
		}
	}()

	var folhas []*entity.FolhaPagamentos
	for rows.Next() {
		var f entity.FolhaPagamentos
		var dataStr string
		if err := rows.Scan(&f.ID, &dataStr); err != nil {
			log.Printf("erro ao ler folha: %v", err)
			continue
		}

		parsed, perr := dateStringToTime.DateStringToTime(dataStr)
		if perr != nil {
			log.Printf("erro ao converter data da folha: %v", perr)
			continue
		}
		f.Data = parsed

		pagamentos, gerr := GetPagamentosByFolhaID(f.ID)
		if gerr != nil {
			log.Printf("erro ao carregar pagamentos da folha %d: %v", f.ID, gerr)
		} else {
			f.Pagamentos = pagamentos
			for _, p := range pagamentos {
				f.Valor += p.Valor
			}
		}

		folhas = append(folhas, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar folhas: %w", err)
	}
	return folhas, nil
}
