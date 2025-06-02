package repository

import (
	"AutoGRH/pkg/entity"
	"fmt"
	"log"
)

// Cria um documento
func CreateDocumento(d *entity.Documento) error {
	query := `INSERT INTO documento (funcionarioID, documento) VALUES (?, ?)`

	result, err := DB.Exec(query, d.FuncionarioId, d.Doc)
	if err != nil {
		return fmt.Errorf("erro ao inserir documento: %w", err)
	}

	d.Id, err = result.LastInsertId()
	return err
}

// Busca documentos por funcionário
func GetDocumentosByFuncionarioID(funcionarioId int64) ([]entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, documento FROM documento WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioId)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documentos: %w", err)
	}
	defer rows.Close()

	var documentos []entity.Documento
	for rows.Next() {
		var d entity.Documento
		err := rows.Scan(&d.Id, &d.FuncionarioId, &d.Doc)
		if err != nil {
			log.Printf("erro ao ler documento: %v", err)
			continue
		}
		documentos = append(documentos, d)
	}
	return documentos, nil
}

// Deleta um documento
func DeleteDocumento(id int64) error {
	query := `DELETE FROM documento WHERE documentoID = ?`
	_, err := DB.Exec(query, id)
	return err
}
