package repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"fmt"
	"log"
)

// CreatePagamento cria um novo pagamento vinculado a funcionário, tipo e (opcionalmente) folha
func CreatePagamento(p *Entity.Pagamento) error {
	query := `INSERT INTO pagamento (funcionarioID, folhaID, tipoID, valor, data)
              VALUES (?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, p.FuncionarioID, p.FolhaID, p.TipoID, p.Valor, p.Data)
	if err != nil {
		return fmt.Errorf("erro ao inserir pagamento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do pagamento inserido: %w", err)
	}
	p.ID = id
	return nil
}

// GetPagamentosByFuncionarioID busca pagamentos por funcionário (mais recentes primeiro)
func GetPagamentosByFuncionarioID(funcionarioID int64) ([]Entity.Pagamento, error) {
	query := `SELECT pagamentoID, funcionarioID, folhaID, tipoID, valor, data
			  FROM pagamento
			  WHERE funcionarioID = ?
			  ORDER BY data DESC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetPagamentosByFuncionarioID: %v", cerr)
		}
	}()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		if err := rows.Scan(&p.ID, &p.FuncionarioID, &p.FolhaID, &p.TipoID, &p.Valor, &dataStr); err != nil {
			log.Printf("erro ao ler pagamento: %v", err)
			continue
		}

		parsed, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do pagamento: %v", err)
			continue
		}
		p.Data = parsed

		pagamentos = append(pagamentos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pagamentos por funcionário: %w", err)
	}
	return pagamentos, nil
}

// UpdatePagamento atualiza um pagamento
func UpdatePagamento(p *Entity.Pagamento) error {
	query := `UPDATE pagamento SET folhaID = ?, tipoID = ?, valor = ?, data = ?
              WHERE pagamentoID = ?`

	_, err := DB.Exec(query, p.FolhaID, p.TipoID, p.Valor, p.Data, p.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pagamento: %w", err)
	}
	return nil
}

// DeletePagamento remove um pagamento por ID
func DeletePagamento(id int64) error {
	query := `DELETE FROM pagamento WHERE pagamentoID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar pagamento: %w", err)
	}
	return nil
}

// ListPagamentos retorna todos os pagamentos (com o tipo resolvido), mais recentes primeiro
func ListPagamentos() ([]Entity.Pagamento, error) {
	query := `SELECT p.pagamentoID, p.funcionarioID, p.folhaID, p.tipoID, t.tipo, p.data, p.valor
              FROM pagamento p
              JOIN tipo_pagamento t ON p.tipoID = t.tipoID
              ORDER BY p.data DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListPagamentos: %v", cerr)
		}
	}()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		if err := rows.Scan(&p.ID, &p.FuncionarioID, &p.FolhaID, &p.TipoID, &p.Tipo, &dataStr, &p.Valor); err != nil {
			log.Printf("erro ao ler pagamento: %v", err)
			continue
		}

		parsed, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do pagamento: %v", err)
			continue
		}
		p.Data = parsed

		pagamentos = append(pagamentos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pagamentos: %w", err)
	}
	return pagamentos, nil
}

// GetPagamentosByFolhaID busca pagamentos de uma folha específica
func GetPagamentosByFolhaID(folhaID int64) ([]Entity.Pagamento, error) {
	query := `SELECT p.pagamentoID, p.funcionarioID, p.folhaID, p.tipoID, t.tipo, p.valor, p.data
              FROM pagamento p
              JOIN tipo_pagamento t ON p.tipoID = t.tipoID
              WHERE p.folhaID = ?`

	rows, err := DB.Query(query, folhaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos por folha: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetPagamentosByFolhaID: %v", cerr)
		}
	}()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		if err := rows.Scan(&p.ID, &p.FuncionarioID, &p.FolhaID, &p.TipoID, &p.Tipo, &p.Valor, &dataStr); err != nil {
			log.Printf("erro ao ler pagamento: %v", err)
			continue
		}

		parsed, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do pagamento: %v", err)
			continue
		}
		p.Data = parsed

		pagamentos = append(pagamentos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pagamentos por folha: %w", err)
	}
	return pagamentos, nil
}
