package Repository

import (
	"AutoGRH/pkg/Entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Cria um novo vale
func CreateVale(v *Entity.Vale) error {
	query := `INSERT INTO vale (funcionarioID, valor, data, aprovado, pago)
              VALUES (?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, v.FuncionarioId, v.Valor, v.Data.Format("2006-01-02"), v.Aprovado, v.Pago)
	if err != nil {
		return fmt.Errorf("erro ao criar vale: %w", err)
	}

	v.Id, err = result.LastInsertId()
	return err
}

// Busca um vale por ID
func GetValeByID(id int64) (*Entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago
              FROM vale WHERE valeID = ?`

	row := DB.QueryRow(query, id)

	var v Entity.Vale
	var dataStr string
	err := row.Scan(&v.Id, &v.FuncionarioId, &v.Valor, &dataStr, &v.Aprovado, &v.Pago)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar vale: %w", err)
	}

	v.Data, err = time.Parse("2006-01-02", dataStr)
	if err != nil {
		log.Printf("erro ao converter data do vale: %v", err)
	}

	return &v, nil
}

// Lista todos os vales de um funcionário
func GetValesByFuncionarioID(funcionarioId int64) ([]Entity.Vale, error) {
	query := `SELECT valeID, funcionarioID, valor, data, aprovado, pago
              FROM vale WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioId)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vales: %w", err)
	}
	defer rows.Close()

	var vales []Entity.Vale
	for rows.Next() {
		var v Entity.Vale
		var dataStr string
		err := rows.Scan(&v.Id, &v.FuncionarioId, &v.Valor, &dataStr, &v.Aprovado, &v.Pago)
		if err != nil {
			log.Println("erro ao ler vale:", err)
			continue
		}

		v.Data, err = time.Parse("2006-01-02", dataStr)
		if err != nil {
			log.Printf("erro ao converter data do vale: %v", err)
			continue
		}

		vales = append(vales, v)
	}
	return vales, nil
}

// Atualiza um vale
func UpdateVale(v *Entity.Vale) error {
	query := `UPDATE vale SET valor = ?, data = ?, aprovado = ?, pago = ? WHERE valeID = ?`

	_, err := DB.Exec(query, v.Valor, v.Data.Format("2006-01-02"), v.Aprovado, v.Pago, v.Id)
	return err
}

// Deleta um vale
func DeleteVale(id int64) error {
	query := `DELETE FROM vale WHERE valeID = ?`
	_, err := DB.Exec(query, id)
	return err
}
