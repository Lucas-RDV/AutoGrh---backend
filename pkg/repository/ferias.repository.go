package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateFerias cria um novo período de férias
func CreateFerias(f *entity.Ferias) error {
	query := `INSERT INTO ferias 
	(funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		f.FuncionarioID, f.Dias, f.Inicio, f.Vencimento,
		f.Vencido, f.Valor, f.Pago, f.Terco, f.TercoPago,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir férias: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID de férias inseridas: %w", err)
	}
	f.ID = id
	return nil
}

// GetFeriasByID busca férias por ID, incluindo descansos
func GetFeriasByID(id int64) (*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago
	          FROM ferias WHERE feriasID = ?`

	row := DB.QueryRow(query, id)

	var f entity.Ferias
	var inicioStr, vencimentoStr string
	if err := row.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr,
		&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar férias: %w", err)
	}

	var err error
	f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de início: %w", err)
	}
	f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data de vencimento: %w", err)
	}

	// Carrega descansos
	descansos, err := GetDescansosByFeriasID(f.ID)
	if err != nil {
		log.Printf("erro ao carregar descansos: %v", err)
	} else {
		f.Descansos = make([]entity.Descanso, 0, len(descansos))
		for _, d := range descansos {
			f.Descansos = append(f.Descansos, *d)
		}
	}

	return &f, nil
}

// UpdateFerias atualiza um período de férias
func UpdateFerias(f *entity.Ferias) error {
	query := `UPDATE ferias SET dias = ?, inicio = ?, vencimento = ?, vencido = ?, 
	          valor = ?, pago = ?, terco = ?, tercoPago = ? WHERE feriasID = ?`

	_, err := DB.Exec(query,
		f.Dias, f.Inicio, f.Vencimento, f.Vencido,
		f.Valor, f.Pago, f.Terco, f.TercoPago, f.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar férias: %w", err)
	}
	return nil
}

// DeleteFerias remove um período de férias (por enquanto DELETE físico, futuramente -> soft delete)
func DeleteFerias(id int64) error {
	query := `DELETE FROM ferias WHERE feriasID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar férias: %w", err)
	}
	return nil
}

// GetFeriasByFuncionarioID lista todas as férias de um funcionário (com descansos)
func GetFeriasByFuncionarioID(funcionarioID int64) ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago
	          FROM ferias WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar férias por funcionário: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetFeriasByFuncionarioID: %v", cerr)
		}
	}()

	var lista []*entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		var inicioStr, vencimentoStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr,
			&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago); err != nil {
			log.Printf("erro ao ler férias: %v", err)
			continue
		}

		if f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr); err != nil {
			log.Printf("erro ao converter data de vencimento: %v", err)
			continue
		}

		// Carrega descansos
		descansos, derr := GetDescansosByFeriasID(f.ID)
		if derr != nil {
			log.Printf("erro ao carregar descansos: %v", derr)
		} else {
			for _, d := range descansos {
				f.Descansos = append(f.Descansos, *d)
			}
		}

		lista = append(lista, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar férias: %w", err)
	}
	return lista, nil
}

// ListFerias lista todos os registros de férias
func ListFerias() ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago FROM ferias`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar férias: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListFerias: %v", cerr)
		}
	}()

	var lista []*entity.Ferias
	for rows.Next() {
		var f *entity.Ferias
		var inicioStr, vencimentoStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &inicioStr, &vencimentoStr,
			&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago); err != nil {
			log.Printf("erro ao ler férias: %v", err)
			continue
		}

		if f.Inicio, err = dateStringToTime.DateStringToTime(inicioStr); err != nil {
			log.Printf("erro ao converter data de início: %v", err)
			continue
		}
		if f.Vencimento, err = dateStringToTime.DateStringToTime(vencimentoStr); err != nil {
			log.Printf("erro ao converter data de vencimento: %v", err)
			continue
		}

		lista = append(lista, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar férias: %w", err)
	}
	return lista, nil
}

// GetFeriasAtivas retorna férias não vencidas
func GetFeriasAtivas(funcionarioID int64) ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago
			  FROM ferias WHERE funcionarioID = ? AND vencido = FALSE`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar férias ativas: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &f.Inicio, &f.Vencimento,
			&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler férias ativas: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}

// GetFeriasVencidas retorna férias vencidas
func GetFeriasVencidas(funcionarioID int64) ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago
			  FROM ferias WHERE funcionarioID = ? AND vencido = TRUE`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar férias vencidas: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &f.Inicio, &f.Vencimento,
			&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler férias vencidas: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}

// GetFeriasNaoPagas retorna férias não pagas
func GetFeriasNaoPagas(funcionarioID int64) ([]*entity.Ferias, error) {
	query := `SELECT feriasID, funcionarioID, dias, inicio, vencimento, vencido, valor, pago, terco, tercoPago
			  FROM ferias WHERE funcionarioID = ? AND pago = FALSE`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar férias não pagas: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Ferias
	for rows.Next() {
		var f entity.Ferias
		err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Dias, &f.Inicio, &f.Vencimento,
			&f.Vencido, &f.Valor, &f.Pago, &f.Terco, &f.TercoPago)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler férias não pagas: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}

// MarcarFeriasComoPagas marca férias como quitadas
func MarcarFeriasComoPagas(feriasID int64) error {
	query := `UPDATE ferias SET pago = TRUE, tercoPago = TRUE WHERE feriasID = ?`
	_, err := DB.Exec(query, feriasID)
	if err != nil {
		return fmt.Errorf("erro ao marcar férias como pagas: %w", err)
	}
	return nil
}

// ConsumirDiasFerias subtrai uma quantidade de dias de férias
func ConsumirDiasFerias(feriasID int64, dias int) error {
	query := `UPDATE ferias SET dias = dias - ? WHERE feriasID = ? AND dias >= ?`
	_, err := DB.Exec(query, dias, feriasID, dias)
	if err != nil {
		return fmt.Errorf("erro ao consumir dias de férias: %w", err)
	}
	return nil
}
