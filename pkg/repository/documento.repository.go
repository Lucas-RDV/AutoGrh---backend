package repository

import (
	"AutoGRH/pkg/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// CreateDocumento insere um documento no banco
func CreateDocumento(ctx context.Context, d *entity.Documento) error {
	query := `INSERT INTO documento (funcionarioID, caminho)
			  VALUES (?, ?)`
	result, err := DB.ExecContext(ctx, query, d.FuncionarioID, d.Caminho)
	if err != nil {
		return fmt.Errorf("erro ao inserir documento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do documento: %w", err)
	}
	d.ID = id
	return nil
}

// GetDocumentosByFuncionarioID retorna todos documentos de um funcionário
func GetDocumentosByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, caminho
			  FROM documento WHERE funcionarioID = ?`

	rows, err := DB.QueryContext(ctx, query, funcionarioID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar documentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetDocumentosByFuncionarioID: %v", cerr)
		}
	}()

	var docs []*entity.Documento
	for rows.Next() {
		var d entity.Documento
		if err := rows.Scan(&d.ID, &d.FuncionarioID, &d.Caminho); err != nil {
			return nil, fmt.Errorf("erro ao ler documento: %w", err)
		}
		docs = append(docs, &d)
	}
	return docs, nil
}

// GetByID retorna um documento pelo ID
func GetByID(ctx context.Context, id int64) (*entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, caminho
			  FROM documento WHERE documentoID = ?`

	row := DB.QueryRowContext(ctx, query, id)

	var d entity.Documento
	if err := row.Scan(&d.ID, &d.FuncionarioID, &d.Caminho); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar documento por ID: %w", err)
	}

	return &d, nil
}

// ListDocumentos retorna todos os documentos cadastrados
func ListDocumentos(ctx context.Context) ([]*entity.Documento, error) {
	query := `SELECT documentoID, funcionarioID, caminho FROM documento`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar documentos: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListDocumentos: %v", cerr)
		}
	}()

	var docs []*entity.Documento
	for rows.Next() {
		var d entity.Documento
		if err := rows.Scan(&d.ID, &d.FuncionarioID, &d.Caminho); err != nil {
			return nil, fmt.Errorf("erro ao ler documento: %w", err)
		}
		docs = append(docs, &d)
	}
	return docs, nil
}

// DeleteDocumento remove um documento pelo ID
func DeleteDocumento(ctx context.Context, id int64) error {
	query := `DELETE FROM documento WHERE documentoID = ?`
	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar documento: %w", err)
	}
	return nil
}
