package Repository

import (
	"AutoGRH/pkg/Entity"
	"fmt"
	"log"
	"time"
)

// Cria um novo pagamento vinculado a funcionário, tipo e folha (opcional)
func CreatePagamento(p *Entity.Pagamento) error {
	query := `INSERT INTO pagamento (funcionarioID, folhaID, tipoID, valor, data)
              VALUES (?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, p.FuncionarioId, p.FolhaId, p.TipoId, p.Valor, p.Data)
	if err != nil {
		return fmt.Errorf("erro ao inserir pagamento: %w", err)
	}

	p.Id, err = result.LastInsertId()
	return err
}

// Busca pagamentos por funcionário
func GetPagamentosByFuncionarioID(funcionarioID int64) ([]Entity.Pagamento, error) {
	query := `SELECT pagamentoID, funcionarioID, folhaID, tipoID, valor, data
			  FROM pagamento
			  WHERE funcionarioID = ?
			  ORDER BY data DESC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos: %w", err)
	}
	defer rows.Close()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		err := rows.Scan(&p.Id, &p.FuncionarioId, &p.FolhaId, &p.TipoId, &p.Valor, &dataStr)
		if err != nil {
			log.Printf("erro ao ler pagamento: %v", err)
			continue
		}

		p.Data, err = time.Parse("2006-01-02", dataStr)
		if err != nil {
			log.Printf("erro ao converter data: %v", err)
			continue
		}

		pagamentos = append(pagamentos, p)
	}
	return pagamentos, nil
}

// Atualiza um pagamento
func UpdatePagamento(p *Entity.Pagamento) error {
	query := `UPDATE pagamento SET folhaID = ?, tipoID = ?, valor = ?, data = ?
              WHERE pagamentoID = ?`

	_, err := DB.Exec(query, p.FolhaId, p.TipoId, p.Valor, p.Data, p.Id)
	return err
}

// Deleta um pagamento
func DeletePagamento(id int64) error {
	query := `DELETE FROM pagamento WHERE pagamentoID = ?`
	_, err := DB.Exec(query, id)
	return err
}

// ListPagamentos retorna todos os pagamentos registrados no sistema
func ListPagamentos() ([]Entity.Pagamento, error) {
	query := `SELECT p.pagamentoID, p.funcionarioID, p.folhaID, p.tipoID, t.tipo, p.data, p.valor
              FROM pagamento p
              JOIN tipo_pagamento t ON p.tipoID = t.tipoID
              ORDER BY p.data DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos: %w", err)
	}
	defer rows.Close()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		err := rows.Scan(&p.Id, &p.FuncionarioId, &p.FolhaId, &p.TipoId, &p.Tipo, &dataStr, &p.Valor)
		if err != nil {
			log.Printf("erro ao ler pagamento: %v", err)
			continue
		}

		// Corrigido: formato apenas com data
		p.Data, err = time.Parse("2006-01-02", dataStr)
		if err != nil {
			log.Printf("erro ao converter data do pagamento: %v", err)
			continue
		}

		pagamentos = append(pagamentos, p)
	}

	return pagamentos, nil
}

// Buscar pagamentos de uma folha
func GetPagamentosByFolhaID(folhaID int64) ([]Entity.Pagamento, error) {
	query := `SELECT p.pagamentoID, p.funcionarioID, p.folhaID, p.tipoID, t.tipo, p.valor, p.data
              FROM pagamento p
              JOIN tipo_pagamento t ON p.tipoID = t.tipoID
              WHERE p.folhaID = ?`

	rows, err := DB.Query(query, folhaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos por folha: %w", err)
	}
	defer rows.Close()

	var pagamentos []Entity.Pagamento
	for rows.Next() {
		var p Entity.Pagamento
		var dataStr string

		err := rows.Scan(&p.Id, &p.FuncionarioId, &p.FolhaId, &p.TipoId, &p.Tipo, &p.Valor, &dataStr)
		if err != nil {
			log.Println("erro ao ler pagamento:", err)
			continue
		}

		parsed, err := time.Parse("2006-01-02", dataStr)
		if err != nil {
			log.Println("erro ao converter data do pagamento:", err)
			continue
		}
		p.Data = parsed

		pagamentos = append(pagamentos, p)
	}
	return pagamentos, nil
}
