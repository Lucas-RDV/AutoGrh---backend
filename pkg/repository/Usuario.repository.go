package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
)

// CreateUsuario Cria um novo usuário no banco
func CreateUsuario(u *entity.Usuario) error {
	query := "INSERT INTO usuario (username, password, isAdmin) VALUES (?, ?, ?)"
	result, err := DB.Exec(query, u.Username, u.Password, u.IsAdmin)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	u.Id, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID inserido: %w", err)
	}

	return nil
}

// GetUsuarioByID Busca um usuário pelo ID
func GetUsuarioByID(id int64) (*entity.Usuario, error) {
	query := "SELECT usuarioID, username, password, isAdmin FROM usuario WHERE usuarioID = ?"
	row := DB.QueryRow(query, id)

	var u entity.Usuario
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuário não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return &u, nil
}

// UpdateUsuario Atualiza um usuário existente
func UpdateUsuario(u *entity.Usuario) error {
	query := "UPDATE usuario SET username = ?, password = ?, isAdmin = ? WHERE usuarioID = ?"
	_, err := DB.Exec(query, u.Username, u.Password, u.IsAdmin, u.Id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}
	return nil
}

// DeleteUsuario Deleta um usuário pelo ID
func DeleteUsuario(id int64) error {
	query := "DELETE FROM usuario WHERE usuarioID = ?"
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}

// GetAllUsuarios Lista todos os usuários
func GetAllUsuarios() ([]*entity.Usuario, error) {
	query := "SELECT usuarioID, username, password, isAdmin FROM usuario"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}
	defer rows.Close()

	var usuarios []*entity.Usuario
	for rows.Next() {
		var u entity.Usuario
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.IsAdmin)
		if err != nil {
			log.Printf("erro ao ler linha: %v", err)
			continue
		}
		usuarios = append(usuarios, &u)
	}

	return usuarios, nil
}
