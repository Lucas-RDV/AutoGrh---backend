package Repository

import (
	"AutoGRH/pkg/Entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Cria uma nova folha de pagamento
func CreateFolha(f *Entity.FolhaPagamentos) error {
	query := `INSERT INTO folha_pagamento (data) VALUES (?)`

	result, err := DB.Exec(query, f.Data.Format("2006-01-02"))
	if err != nil {
		return fmt.Errorf("erro ao criar folha: %w", err)
	}

	f.Id, err = result.LastInsertId()
	return err
}

// Busca uma folha de pagamento e seus pagamentos
func GetFolhaByID(id int64) (*Entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento WHERE folhaID = ?`
	row := DB.QueryRow(query, id)

	var f Entity.FolhaPagamentos
	var dataStr string
	err := row.Scan(&f.Id, &dataStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar folha: %w", err)
	}

	parsed, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data: %w", err)
	}
	f.Data = parsed

	// Buscar pagamentos relacionados
	pagamentos, err := GetPagamentosByFolhaID(f.Id)
	if err == nil {
		f.Pagamentos = pagamentos
		for _, p := range pagamentos {
			f.Valor += p.Valor
		}
	}

	return &f, nil
}

// Lista todas as folhas de pagamento
func ListFolhas() ([]*Entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento ORDER BY data DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar folhas: %w", err)
	}
	defer rows.Close()

	var folhas []*Entity.FolhaPagamentos
	for rows.Next() {
		var f Entity.FolhaPagamentos
		var dataStr string
		err := rows.Scan(&f.Id, &dataStr)
		if err != nil {
			log.Println("erro ao ler folha:", err)
			continue
		}

		parsed, err := time.Parse("2006-01-02", dataStr)
		if err != nil {
			log.Println("erro ao converter data da folha:", err)
			continue
		}
		f.Data = parsed

		pagamentos, err := GetPagamentosByFolhaID(f.Id)
		if err == nil {
			f.Pagamentos = pagamentos
			for _, p := range pagamentos {
				f.Valor += p.Valor
			}
		}

		folhas = append(folhas, &f)
	}
	return folhas, nil
}
