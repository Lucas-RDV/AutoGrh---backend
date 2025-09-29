package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/timeToDateString"
	"database/sql"
	"fmt"
	"log"
)

// CreateVale cria um novo vale (inicia como ativo = true, aprovado = false, pago = false)
func CreateVale(v *entity.Vale) error {
	query := `INSERT INTO vale (funcionarioID, valor, data, aprovado, pago, ativo)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, v.FuncionarioID, v.Valor, timeToDateString.TimeToDateString(v.Data), v.Aprovado, v.Pago, v.Ativo)
	if err != nil {
		return fmt.Errorf("erro ao inserir vale: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do vale inserido: %w", err)
	}
	v.ID = id
	return nil
}

// GetValeByID busca um vale pelo ID
func GetValeByID(id int64) (*entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago, ativo
			  FROM vale WHERE valeID = ?`
	row := DB.QueryRow(query, id)

	var v entity.Vale
	var dataStr string
	if err := row.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago, &v.Ativo); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar vale: %w", err)
	}

	t, err := dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data do vale: %w", err)
	}
	v.Data = t
	return &v, nil
}

// GetValesByFuncionarioID lista todos os vales de um funcionário
func GetValesByFuncionarioID(funcionarioID int64) ([]entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago, ativo
			  FROM vale WHERE funcionarioID = ? ORDER BY data DESC`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vales do funcionario: %w", err)
	}
	defer rows.Close()

	var vales []entity.Vale
	for rows.Next() {
		var v entity.Vale
		var dataStr string
		if err := rows.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago, &v.Ativo); err != nil {
			log.Printf("erro ao ler vale: %v", err)
			continue
		}
		t, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do vale: %v", err)
			continue
		}
		v.Data = t
		vales = append(vales, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar vales: %w", err)
	}
	return vales, nil
}

// UpdateVale atualiza os campos de um vale existente
func UpdateVale(v *entity.Vale) error {
	query := `UPDATE vale SET funcionarioID = ?, valor = ?, data = ?, aprovado = ?, pago = ?, ativo = ?
	          WHERE valeID = ?`

	_, err := DB.Exec(query, v.FuncionarioID, v.Valor, timeToDateString.TimeToDateString(v.Data), v.Aprovado, v.Pago, v.Ativo, v.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar vale: %w", err)
	}
	return nil
}

// SoftDeleteVale realiza exclusão lógica de um vale (ativo = false)
func SoftDeleteVale(id int64) error {
	query := `UPDATE vale SET ativo = FALSE WHERE valeID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao inativar vale: %w", err)
	}
	return nil
}

// DeleteVale remove permanentemente um vale do banco (exclusão física)
func DeleteVale(id int64) error {
	query := `DELETE FROM vale WHERE valeID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir vale permanentemente: %w", err)
	}
	return nil
}

// ListValesPendentes retorna todos os vales ativos que aguardam aprovação
func ListValesPendentes() ([]entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago, ativo
			  FROM vale WHERE ativo = TRUE AND aprovado = FALSE AND pago = FALSE`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar vales pendentes: %w", err)
	}
	defer rows.Close()

	var vales []entity.Vale
	for rows.Next() {
		var v entity.Vale
		var dataStr string
		if err := rows.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago, &v.Ativo); err != nil {
			log.Printf("erro ao ler vale pendente: %v", err)
			continue
		}
		t, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do vale pendente: %v", err)
			continue
		}
		v.Data = t
		vales = append(vales, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar vales pendentes: %w", err)
	}
	return vales, nil
}

// ListValesAprovadosNaoPagos retorna todos os vales ativos aprovados mas ainda não pagos
func ListValesAprovadosNaoPagos() ([]entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago, ativo
			  FROM vale WHERE ativo = TRUE AND aprovado = TRUE AND pago = FALSE`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar vales aprovados não pagos: %w", err)
	}
	defer rows.Close()

	var vales []entity.Vale
	for rows.Next() {
		var v entity.Vale
		var dataStr string
		if err := rows.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago, &v.Ativo); err != nil {
			log.Printf("erro ao ler vale aprovado: %v", err)
			continue
		}
		t, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			log.Printf("erro ao converter data do vale aprovado: %v", err)
			continue
		}
		v.Data = t
		vales = append(vales, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar vales aprovados não pagos: %w", err)
	}
	return vales, nil
}

// GetValesByFuncionarioMesAno retorna os vales de um funcionário em um mês/ano específico
func GetValesByFuncionarioMesAno(funcionarioID int64, mes int, ano int) ([]entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago, ativo
			  FROM vale
			  WHERE funcionarioID = ?
			    AND MONTH(data) = ?
			    AND YEAR(data) = ?
				AND ativo = 1`

	rows, err := DB.Query(query, funcionarioID, mes, ano)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vales do funcionário %d em %02d/%d: %w",
			funcionarioID, mes, ano, err)
	}
	defer rows.Close()

	var vales []entity.Vale
	for rows.Next() {
		var v entity.Vale
		var dataStr string

		if err := rows.Scan(&v.ID, &v.FuncionarioID, &v.Valor, &dataStr, &v.Aprovado, &v.Pago, &v.Ativo); err != nil {
			return nil, fmt.Errorf("erro ao ler vale: %w", err)
		}

		t, err := dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data do vale: %w", err)
		}
		v.Data = t

		vales = append(vales, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar vales: %w", err)
	}

	return vales, nil
}

func MarcarValesComoPagos(mes int, ano int) error {
	query := `UPDATE vale 
              SET pago = TRUE 
              WHERE MONTH(data) = ? AND YEAR(data) = ? AND aprovado = TRUE AND ativo = TRUE`
	_, err := DB.Exec(query, mes, ano)
	if err != nil {
		return fmt.Errorf("erro ao marcar vales como pagos: %w", err)
	}
	return nil
}

// MarcarTodosValesComoPagos marca todos os vales aprovados e não pagos como pagos.
func MarcarTodosValesComoPagos() error {
	query := `UPDATE vale 
              SET pago = 1 
              WHERE aprovado = 1 AND ativo = 1 AND pago = 0`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("erro ao marcar todos os vales como pagos: %w", err)
	}
	return nil
}
