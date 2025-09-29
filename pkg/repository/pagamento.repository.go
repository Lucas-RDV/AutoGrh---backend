package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
)

// CreatePagamento insere um novo pagamento no banco
func CreatePagamento(p *entity.Pagamento) error {
	query := `INSERT INTO pagamento 
		(funcionarioID, folhaID, salarioBase, adicional, descontoINSS, salarioFamilia, descontoVales, valorFinal, pago)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		p.FuncionarioID,
		p.FolhaID,
		p.SalarioBase,
		p.Adicional,
		p.DescontoINSS,
		p.SalarioFamilia,
		p.DescontoVales,
		p.ValorFinal,
		p.Pago,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir pagamento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do pagamento: %w", err)
	}
	p.ID = id
	return nil
}

// UpdatePagamento atualiza os dados de um pagamento existente
func UpdatePagamento(p *entity.Pagamento) error {
	query := `UPDATE pagamento 
		SET salarioBase = ?, adicional = ?, descontoINSS = ?, salarioFamilia = ?, descontoVales = ?, valorFinal = ?, pago = ?
		WHERE pagamentoID = ?`

	_, err := DB.Exec(query,
		p.SalarioBase,
		p.Adicional,
		p.DescontoINSS,
		p.SalarioFamilia,
		p.DescontoVales,
		p.ValorFinal,
		p.Pago,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pagamento: %w", err)
	}
	return nil
}

// GetPagamentoByID retorna um pagamento pelo ID
func GetPagamentoByID(id int64) (*entity.Pagamento, error) {
	query := `SELECT pagamentoID, funcionarioID, folhaID, salarioBase, adicional, descontoINSS, salarioFamilia, descontoVales, valorFinal, pago
			  FROM pagamento WHERE pagamentoID = ?`

	var p entity.Pagamento
	err := DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.FuncionarioID,
		&p.FolhaID,
		&p.SalarioBase,
		&p.Adicional,
		&p.DescontoINSS,
		&p.SalarioFamilia,
		&p.DescontoVales,
		&p.ValorFinal,
		&p.Pago,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pagamento: %w", err)
	}
	return &p, nil
}

// GetPagamentosByFolhaID retorna todos os pagamentos de uma folha
func GetPagamentosByFolhaID(folhaID int64) ([]entity.Pagamento, error) {
	query := `SELECT pagamentoID, funcionarioID, folhaID, salarioBase, adicional, descontoINSS, salarioFamilia, descontoVales, valorFinal, pago
			  FROM pagamento WHERE folhaID = ?`

	rows, err := DB.Query(query, folhaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos da folha %d: %w", folhaID, err)
	}
	defer rows.Close()

	var pagamentos []entity.Pagamento
	for rows.Next() {
		var p entity.Pagamento
		if err := rows.Scan(
			&p.ID,
			&p.FuncionarioID,
			&p.FolhaID,
			&p.SalarioBase,
			&p.Adicional,
			&p.DescontoINSS,
			&p.SalarioFamilia,
			&p.DescontoVales,
			&p.ValorFinal,
			&p.Pago,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler pagamento: %w", err)
		}
		pagamentos = append(pagamentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pagamentos: %w", err)
	}

	return pagamentos, nil
}

// ListPagamentosByFuncionarioID lista os pagamentos de um funcionário
func ListPagamentosByFuncionarioID(funcionarioID int64) ([]entity.Pagamento, error) {
	query := `SELECT pagamentoID, funcionarioID, folhaID, salarioBase, adicional, descontoINSS, salarioFamilia, descontoVales, valorFinal, pago
			  FROM pagamento WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pagamentos do funcionário %d: %w", funcionarioID, err)
	}
	defer rows.Close()

	var pagamentos []entity.Pagamento
	for rows.Next() {
		var p entity.Pagamento
		if err := rows.Scan(
			&p.ID,
			&p.FuncionarioID,
			&p.FolhaID,
			&p.SalarioBase,
			&p.Adicional,
			&p.DescontoINSS,
			&p.SalarioFamilia,
			&p.DescontoVales,
			&p.ValorFinal,
			&p.Pago,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler pagamento: %w", err)
		}
		pagamentos = append(pagamentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pagamentos: %w", err)
	}

	return pagamentos, nil
}

// DeletePagamentosByFolhaID remove todos os pagamentos de uma folha
func DeletePagamentosByFolhaID(folhaID int64) error {
	query := `DELETE FROM pagamento WHERE folhaID = ?`
	_, err := DB.Exec(query, folhaID)
	if err != nil {
		return fmt.Errorf("erro ao deletar pagamentos da folha %d: %w", folhaID, err)
	}
	return nil
}

// MarcarPagamentosDaFolhaComoPagos marca todos os pagamentos de uma folha como pagos
func MarcarPagamentosDaFolhaComoPagos(folhaID int64) error {
	query := `UPDATE pagamento SET pago = 1 WHERE folhaID = ?`
	_, err := DB.Exec(query, folhaID)
	if err != nil {
		return fmt.Errorf("erro ao marcar pagamentos da folha %d como pagos: %w", folhaID, err)
	}
	return nil
}
