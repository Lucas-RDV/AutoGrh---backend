package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/timeToDateString"
	"database/sql"
	"fmt"
	"log"
)

// CreateDescanso cria um descanso vinculado a um período de férias
func CreateDescanso(d *entity.Descanso) error {
	query := `INSERT INTO descanso (feriasID, inicio, fim, valor, pago, aprovado)
              VALUES (?, ?, ?, ?, ?, ?)`

	res, err := DB.Exec(
		query,
		d.FeriasID,
		timeToDateString.TimeToDateString(d.Inicio),
		timeToDateString.TimeToDateString(d.Fim),
		d.Valor,
		d.Pago,
		d.Aprovado,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir descanso: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do descanso inserido: %w", err)
	}
	d.ID = id
	return nil
}

// GetDescansoByID busca um descanso por ID
func GetDescansoByID(id int64) (*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
	          FROM descanso WHERE descansoID = ?`

	row := DB.QueryRow(query, id)

	var d entity.Descanso
	var inicioStr, fimStr string

	if err := row.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar descanso: %w", err)
	}

	var err error
	if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
		return nil, fmt.Errorf("erro ao converter data de início: %w", err)
	}
	if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
		return nil, fmt.Errorf("erro ao converter data de fim: %w", err)
	}

	return &d, nil
}

// GetDescansosByFeriasID busca todos os descansos de um período de férias
func GetDescansosByFeriasID(feriasID int64) ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado
	          FROM descanso WHERE feriasID = ?`

	rows, err := DB.Query(query, feriasID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar descansos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetDescansosByFeriasID: %v", cerr)
		}
	}()

	var descansos []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		if err := rows.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado); err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
			log.Printf("erro ao converter data de fim: %v", err)
			continue
		}

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// ListDescansos lista todos os descansos
func ListDescansos() ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, pago, aprovado FROM descanso`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar descansos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListDescansos: %v", cerr)
		}
	}()

	var descansos []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		if err := rows.Scan(&d.ID, &d.FeriasID, &inicioStr, &fimStr, &d.Valor, &d.Pago, &d.Aprovado); err != nil {
			log.Printf("erro ao ler descanso: %v", err)
			continue
		}

		if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
			log.Printf("erro ao converter data de fim: %v", err)
			continue
		}

		descansos = append(descansos, &d)
	}
	return descansos, nil
}

// UpdateDescanso atualiza um descanso
func UpdateDescanso(d *entity.Descanso) error {
	query := `UPDATE descanso SET inicio = ?, fim = ?, valor = ?, pago = ?, aprovado = ?
	          WHERE descansoID = ?`

	_, err := DB.Exec(
		query,
		timeToDateString.TimeToDateString(d.Inicio),
		timeToDateString.TimeToDateString(d.Fim),
		d.Valor,
		d.Pago,
		d.Aprovado,
		d.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar descanso: %w", err)
	}
	return nil
}

// DeleteDescanso deleta um descanso (hard delete)
func DeleteDescanso(id int64) error {
	query := `DELETE FROM descanso WHERE descansoID = ?`
	if _, err := DB.Exec(query, id); err != nil {
		return fmt.Errorf("erro ao deletar descanso: %w", err)
	}
	return nil
}

// GetDescansosAprovados retorna todos os descansos aprovados
func GetDescansosAprovados() ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, aprovado, pago
			  FROM descanso WHERE aprovado = 1`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar descansos aprovados: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		if err := rows.Scan(
			&d.ID,
			&d.FeriasID,
			&inicioStr,
			&fimStr,
			&d.Valor,
			&d.Aprovado,
			&d.Pago,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler descanso aprovado: %w", err)
		}

		if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			return nil, fmt.Errorf("erro ao converter inicio: %w", err)
		}
		if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
			return nil, fmt.Errorf("erro ao converter fim: %w", err)
		}

		lista = append(lista, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no iterador de descansos aprovados: %w", err)
	}
	return lista, nil
}

// GetDescansosPendentes retorna todos os descansos pendentes (a aprovar)
func GetDescansosPendentes() ([]*entity.Descanso, error) {
	query := `SELECT descansoID, feriasID, inicio, fim, valor, aprovado, pago
			  FROM descanso WHERE aprovado = 0`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar descansos pendentes: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		if err := rows.Scan(
			&d.ID,
			&d.FeriasID,
			&inicioStr,
			&fimStr,
			&d.Valor,
			&d.Aprovado,
			&d.Pago,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler descanso pendente: %w", err)
		}

		if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			return nil, fmt.Errorf("erro ao converter inicio: %w", err)
		}
		if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
			return nil, fmt.Errorf("erro ao converter fim: %w", err)
		}

		lista = append(lista, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no iterador de descansos pendentes: %w", err)
	}
	return lista, nil
}

// GetDescansosByFuncionarioID retorna todos os descansos de um funcionário (via períodos de férias)
func GetDescansosByFuncionarioID(funcionarioID int64) ([]*entity.Descanso, error) {
	query := `SELECT d.descansoID, d.feriasID, d.inicio, d.fim, d.valor, d.pago, d.aprovado
			  FROM descanso d
			  INNER JOIN ferias f ON d.feriasID = f.feriasID
			  WHERE f.funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar descansos por funcionário: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetDescansosByFuncionarioID: %v", cerr)
		}
	}()

	var lista []*entity.Descanso
	for rows.Next() {
		var d entity.Descanso
		var inicioStr, fimStr string

		if err := rows.Scan(
			&d.ID,
			&d.FeriasID,
			&inicioStr,
			&fimStr,
			&d.Valor,
			&d.Pago,
			&d.Aprovado,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler descanso por funcionário: %w", err)
		}

		if d.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			return nil, fmt.Errorf("erro ao converter inicio: %w", err)
		}
		if d.Fim, err = dateStringToTime.DateStringToTime(fimStr); err != nil {
			return nil, fmt.Errorf("erro ao converter fim: %w", err)
		}

		lista = append(lista, &d)
	}
	return lista, nil
}

// MarcarDescansoComoAprovadoSePendente faz a transição pendente->aprovado de forma atômica.
// Retorna true quando houve transição; false quando já estava aprovado (ou inexistente).
func MarcarDescansoComoAprovadoSePendente(id int64) (bool, error) {
	query := `UPDATE descanso SET aprovado = 1 WHERE descansoID = ? AND aprovado = 0`
	res, err := DB.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("erro ao marcar descanso como aprovado: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar atualização de aprovação do descanso: %w", err)
	}
	return affected > 0, nil
}

// ReverterAprovacaoDescanso reverte a aprovação para pendente (uso compensatório em falhas de consumo).
func ReverterAprovacaoDescanso(id int64) error {
	query := `UPDATE descanso SET aprovado = 0 WHERE descansoID = ?`
	if _, err := DB.Exec(query, id); err != nil {
		return fmt.Errorf("erro ao reverter aprovação do descanso: %w", err)
	}
	return nil
}
