package Repository

import (
	"AutoGRH/pkg/Entity"
	"fmt"
	"log"
)

// CreateDocumento cria um documento vinculado a um funcionário
func CreateDocumento(d *Entity.Documento) error {
	query := `INSERT INTO documento (funcionarioID, documento) VALUES (?, ?)`

	result, err := DB.Exec(query, d.FuncionarioID, d.Doc)
	if err != nil {
		return fmt.Errorf("erro ao inserir documento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do novo documento: %w", err)
	}
	d.ID = id
	return nil
}

// GetDocumentosByFuncionarioID busca documentos por ID de funcionário
func GetDocumentosByFuncionarioID(funcionarioID int64) ([]Entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, documento FROM documento WHERE funcionarioID = ?`

	rows, err := DB.Query(query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetDocumentosByFuncionarioID: %v", cerr)
		}
	}()

	var documentos []Entity.Documento
	for rows.Next() {
		var d Entity.Documento
		if err := rows.Scan(&d.ID, &d.FuncionarioID, &d.Doc); err != nil {
			log.Printf("erro ao ler documento: %v", err)
			continue
		}
		documentos = append(documentos, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar documentos: %w", err)
	}
	return documentos, nil
}

// ListDocumentos lista todos os documentos
func ListDocumentos() ([]Entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, documento FROM documento`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar documentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListDocumentos: %v", cerr)
		}
	}()

	var documentos []Entity.Documento
	for rows.Next() {
		var d Entity.Documento
		if err := rows.Scan(&d.ID, &d.FuncionarioID, &d.Doc); err != nil {
			log.Printf("erro ao ler documento: %v", err)
			continue
		}
		documentos = append(documentos, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar documentos: %w", err)
	}
	return documentos, nil
}

// DeleteDocumento deleta um documento por ID
func DeleteDocumento(id int64) error {
	query := `DELETE FROM documento WHERE documentoID = ?`
	_, err := DB.Exec(query, id)
	return err
}
