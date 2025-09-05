package repository

import (
	"AutoGRH/pkg/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// CreateUsuario cria um novo usuário no banco
func CreateUsuario(u *entity.Usuario) error {
	query := `INSERT INTO usuario (username, password, isAdmin) VALUES (?, ?, ?)`
	result, err := DB.Exec(query, u.Username, u.Password, u.IsAdmin)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID inserido: %w", err)
	}
	u.ID = id
	return nil
}

// GetUsuarioByID busca um usuário pelo ID
func GetUsuarioByID(id int64) (*entity.Usuario, error) {
	query := `SELECT usuarioID, username, password, isAdmin, ativo FROM usuario WHERE usuarioID = ? AND ativo = TRUE`
	row := DB.QueryRow(query, id)

	var u entity.Usuario
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin, &u.Ativo); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return &u, nil
}

// UpdateUsuario atualiza um usuário existente
func UpdateUsuario(u *entity.Usuario) error {
	query := `UPDATE usuario SET username = ?, password = ?, isAdmin = ?, ativo = ? WHERE usuarioID = ?`
	_, err := DB.Exec(query, u.Username, u.Password, u.IsAdmin, u.Ativo, u.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}
	return nil
}

// DeleteUsuario desativa um usuário pelo ID
func DeleteUsuario(id int64) error {
	query := `UPDATE usuario SET ativo = FALSE WHERE usuarioID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao desativar usuário: %w", err)
	}
	return nil
}

// GetAllUsuarios lista todos os usuários
func GetAllUsuarios() ([]*entity.Usuario, error) {
	query := `SELECT usuarioID, username, password, isAdmin, ativo FROM usuario WHERE ativo = TRUE`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em GetAllUsuarios: %v", cerr)
		}
	}()

	var usuarios []*entity.Usuario
	for rows.Next() {
		var u entity.Usuario
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.IsAdmin, &u.Ativo); err != nil {
			log.Printf("erro ao ler linha de usuário: %v", err)
			continue
		}
		usuarios = append(usuarios, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar usuários: %w", err)
	}
	return usuarios, nil
}

func GetUsuarioByUsername(ctx context.Context, username string) (*entity.Usuario, error) {
	if DB == nil {
		return nil, errors.New("repository DB não inicializado")
	}
	u := strings.TrimSpace(username)
	if u == "" {
		return nil, nil // nada a buscar
	}

	const q = `SELECT usuarioID, username, password, isAdmin, ativo FROM usuario WHERE username = ? LIMIT 1`
	row := DB.QueryRowContext(ctx, q, u)

	var out entity.Usuario
	if err := row.Scan(&out.ID, &out.Username, &out.Password, &out.IsAdmin, &out.Ativo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}
