package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"database/sql"
	"fmt"
)

// CreateFalta cria um registro de falta
func CreateFalta(f *entity.Falta) error {
	query := `INSERT INTO falta (funcionarioID, quantidade, data) VALUES (?, ?, ?)`

	result, err := DB.Exec(query, f.FuncionarioID, f.Quantidade, f.Mes)
	if err != nil {
		return fmt.Errorf("erro ao inserir falta: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da falta inserida: %w", err)
	}
	f.ID = id
	return nil
}

// GetFaltasByFuncionarioID busca todas as faltas de um funcionário
func GetFaltasByFuncionarioID(funcionarioID int64) ([]*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas por funcionário: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Falta
	for rows.Next() {
		var f entity.Falta
		var dataStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
			return nil, fmt.Errorf("erro ao ler falta: %w", err)
		}

		f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data: %w", err)
		}

		lista = append(lista, &f)
	}
	return lista, nil
}

// GetFaltaByID retorna uma falta pelo ID
func GetFaltaByID(id int64) (*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta WHERE faltaID = ?`
	row := DB.QueryRow(query, id)

	var f entity.Falta
	var dataStr string
	if err := row.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar falta: %w", err)
	}

	var err error
	f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter data: %w", err)
	}

	return &f, nil
}

// UpdateFalta atualiza um registro de falta
func UpdateFalta(f *entity.Falta) error {
	query := `UPDATE falta SET quantidade = ?, data = ? WHERE faltaID = ?`
	_, err := DB.Exec(query, f.Quantidade, f.Mes, f.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar falta: %w", err)
	}
	return nil
}

// DeleteFalta remove uma falta por ID
func DeleteFalta(id int64) error {
	query := `DELETE FROM falta WHERE faltaID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar falta: %w", err)
	}
	return nil
}

// ListAllFaltas retorna todas as faltas cadastradas
func ListAllFaltas() ([]*entity.Falta, error) {
	query := `SELECT faltaID, funcionarioID, quantidade, data FROM falta`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar faltas: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Falta
	for rows.Next() {
		var f entity.Falta
		var dataStr string
		if err := rows.Scan(&f.ID, &f.FuncionarioID, &f.Quantidade, &dataStr); err != nil {
			return nil, fmt.Errorf("erro ao ler falta: %w", err)
		}

		f.Mes, err = dateStringToTime.DateStringToTime(dataStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data: %w", err)
		}

		lista = append(lista, &f)
	}
	return lista, nil
}

// GetTotalFaltasByFuncionarioMesAno retorna o total de faltas de um funcionário em um mês/ano específico
func GetTotalFaltasByFuncionarioMesAno(funcionarioID int64, mes int, ano int) (int, error) {
	query := `
SELECT COALESCE(SUM(quantidade),0)
FROM falta
WHERE funcionarioID = ?
  AND MONTH(data) = ?
  AND YEAR(data)  = ?`

	row := DB.QueryRow(query, funcionarioID, mes, ano)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("erro ao contar faltas do funcionário %d em %02d/%d: %w",
			funcionarioID, mes, ano, err)
	}

	return total, nil
}

func SetFaltasMensais(funcionarioID int64, mes int, ano int, quantidade int) error {
	if quantidade < 0 {
		quantidade = 0
	}

	// 1) Tenta atualizar a linha do mês/ano
	res, err := DB.Exec(`
		UPDATE falta
		   SET quantidade = ?
		 WHERE funcionarioID = ?
		   AND MONTH(data)   = ?
		   AND YEAR(data)    = ?`,
		quantidade, funcionarioID, mes, ano,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar faltas mensais: %w", err)
	}

	// 2) Se não havia linha para esse mês/ano e quantidade > 0, insere
	rows, _ := res.RowsAffected()
	if rows == 0 && quantidade > 0 {
		primeiroDia := fmt.Sprintf("%04d-%02d-01", ano, mes)
		if _, err := DB.Exec(`
			INSERT INTO falta (funcionarioID, quantidade, data)
			VALUES (?, ?, ?)`,
			funcionarioID, quantidade, primeiroDia,
		); err != nil {
			return fmt.Errorf("erro ao inserir faltas mensais: %w", err)
		}
	}

	// se rows==0 e quantidade==0, não cria linha "vazia" — está ok
	return nil
}
