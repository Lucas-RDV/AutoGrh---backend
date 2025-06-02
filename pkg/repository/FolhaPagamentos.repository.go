package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
)

// Cria uma nova folha de pagamento
func CreateFolha(f *entity.FolhaPagamentos) error {
	query := `INSERT INTO folha_pagamento (data) VALUES (?)`

	result, err := DB.Exec(query, f.Data)
	if err != nil {
		return fmt.Errorf("erro ao criar folha: %w", err)
	}

	f.Id, err = result.LastInsertId()
	return err
}

// Busca uma folha de pagamento e seus pagamentos
func GetFolhaByID(id int64) (*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento WHERE folhaID = ?`
	row := DB.QueryRow(query, id)

	var f entity.FolhaPagamentos
	err := row.Scan(&f.Id, &f.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar folha: %w", err)
	}

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
func ListFolhas() ([]*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, data FROM folha_pagamento ORDER BY data DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar folhas: %w", err)
	}
	defer rows.Close()

	var folhas []*entity.FolhaPagamentos
	for rows.Next() {
		var f entity.FolhaPagamentos
		err := rows.Scan(&f.Id, &f.Data)
		if err != nil {
			log.Println("erro ao ler folha:", err)
			continue
		}

		// Buscar e somar pagamentos
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

// Auxiliar: Buscar pagamentos de uma folha
func GetPagamentosByFolhaID(folhaID int64) ([]entity.Pagamento, error) {
	query := `SELECT p.pagamentoID, p.funcionarioID, p.folhaID, p.tipoID, t.tipo, p.valor, p.data
              FROM pagamento p
              JOIN tipo_pagamento t ON p.tipoID = t.tipoID
              WHERE p.folhaID = ?`

	rows, err := DB.Query(query, folhaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos por folha: %w", err)
	}
	defer rows.Close()

	var pagamentos []entity.Pagamento
	for rows.Next() {
		var p entity.Pagamento
		err := rows.Scan(&p.Id, &p.FuncionarioId, &p.FolhaId, &p.TipoId, &p.Tipo, &p.Valor, &p.Data)
		if err != nil {
			log.Println("erro ao ler pagamento:", err)
			continue
		}
		pagamentos = append(pagamentos, p)
	}
	return pagamentos, nil
}
