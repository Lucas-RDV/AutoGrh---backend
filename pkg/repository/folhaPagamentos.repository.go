package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/timeToDateString"
	"database/sql"
	"fmt"
)

// CreateFolhaPagamento insere uma nova folha no banco
func CreateFolhaPagamento(f *entity.FolhaPagamentos) error {
	query := `INSERT INTO folha_pagamento (mes, ano, tipo, dataGeracao, valorTotal, pago)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		f.Mes,
		f.Ano,
		f.Tipo,
		timeToDateString.TimeToDateString(f.DataGeracao),
		f.ValorTotal,
		f.Pago,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir folha: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da folha inserida: %w", err)
	}
	f.ID = id
	return nil
}

// UpdateFolhaPagamento atualiza os dados de uma folha existente
func UpdateFolhaPagamento(f *entity.FolhaPagamentos) error {
	query := `UPDATE folha_pagamento
	          SET mes = ?, ano = ?, tipo = ?, dataGeracao = ?, valorTotal = ?, pago = ?
	          WHERE folhaID = ?`

	_, err := DB.Exec(query,
		f.Mes,
		f.Ano,
		f.Tipo,
		timeToDateString.TimeToDateString(f.DataGeracao),
		f.ValorTotal,
		f.Pago,
		f.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar folha: %w", err)
	}
	return nil
}

// GetFolhaPagamentoByID busca uma folha pelo ID
func GetFolhaPagamentoByID(id int64) (*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, mes, ano, tipo, dataGeracao, valorTotal, pago
	          FROM folha_pagamento WHERE folhaID = ?`

	row := DB.QueryRow(query, id)

	var f entity.FolhaPagamentos
	var dataStr string

	err := row.Scan(&f.ID, &f.Mes, &f.Ano, &f.Tipo, &dataStr, &f.ValorTotal, &f.Pago)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar folha: %w", err)
	}

	t, err := dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data da folha: %w", err)
	}
	f.DataGeracao = t

	return &f, nil
}

// ListFolhasPagamentos retorna todas as folhas registradas
func ListFolhasPagamentos() ([]entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, mes, ano, tipo, dataGeracao, valorTotal, pago
	          FROM folha_pagamento ORDER BY ano DESC, mes DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar folhas: %w", err)
	}
	defer rows.Close()

	var folhas []entity.FolhaPagamentos
	for rows.Next() {
		var f entity.FolhaPagamentos
		var dataStr string

		if err := rows.Scan(&f.ID, &f.Mes, &f.Ano, &f.Tipo, &dataStr, &f.ValorTotal, &f.Pago); err != nil {
			return nil, fmt.Errorf("erro ao ler folha: %w", err)
		}

		t, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data da folha: %w", err)
		}
		f.DataGeracao = t

		folhas = append(folhas, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar folhas: %w", err)
	}

	return folhas, nil
}

// MarcarFolhaComoPaga atualiza a folha para paga = true
func MarcarFolhaComoPaga(folhaID int64) error {
	query := `UPDATE folha_pagamento SET pago = TRUE WHERE folhaID = ?`
	_, err := DB.Exec(query, folhaID)
	if err != nil {
		return fmt.Errorf("erro ao marcar folha %d como paga: %w", folhaID, err)
	}
	return nil
}

// GetFolhaByMesAnoTipo busca uma folha pelo mês, ano e tipo (ex.: SALARIO, VALE)
func GetFolhaByMesAnoTipo(mes, ano int, tipo string) (*entity.FolhaPagamentos, error) {
	query := `SELECT folhaID, mes, ano, tipo, dataGeracao, valorTotal, pago
			  FROM folha_pagamento WHERE mes = ? AND ano = ? AND tipo = ? LIMIT 1`

	row := DB.QueryRow(query, mes, ano, tipo)

	var f entity.FolhaPagamentos
	var dataGeracaoStr string
	if err := row.Scan(&f.ID, &f.Mes, &f.Ano, &f.Tipo, &dataGeracaoStr, &f.ValorTotal, &f.Pago); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar folha por mes/ano/tipo: %w", err)
	}

	t, err := dateStringToTime.DateStringToTime(dataGeracaoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dataGeracao da folha: %w", err)
	}
	f.DataGeracao = t

	return &f, nil
}

// DeleteFolhaPagamento exclui uma folha de pagamento permanentemente
func DeleteFolhaPagamento(id int64) error {
	query := `DELETE FROM folha_pagamento WHERE folhaID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir folha de pagamento: %w", err)
	}
	return nil
}
