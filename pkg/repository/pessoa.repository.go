package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
)

// CreatePessoa insere uma nova pessoa no banco de dados
func CreatePessoa(p *entity.Pessoa) error {
	query := `INSERT INTO pessoa (nome, cpf, rg, endereco, contato, contatoEmergencia)
			  VALUES (?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, p.Nome, p.CPF, p.RG, p.Endereco, p.Contato, p.ContatoEmergencia)
	if err != nil {
		return fmt.Errorf("erro ao inserir pessoa: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da nova pessoa: %w", err)
	}
	p.ID = id
	return nil
}

// GetPessoaByID retorna uma pessoa pelo ID
func GetPessoaByID(id int64) (*entity.Pessoa, error) {
	query := `SELECT pessoaID, nome, cpf, rg, endereco, contato, contatoEmergencia
			  FROM pessoa WHERE pessoaID = ?`
	row := DB.QueryRow(query, id)

	var p entity.Pessoa
	if err := row.Scan(&p.ID, &p.Nome, &p.CPF, &p.RG, &p.Endereco, &p.Contato, &p.ContatoEmergencia); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pessoa: %w", err)
	}
	return &p, nil
}

// GetPessoaByCPF retorna uma pessoa pelo CPF
func GetPessoaByCPF(cpf string) (*entity.Pessoa, error) {
	query := `SELECT pessoaID, nome, cpf, rg, endereco, contato, contatoEmergencia
			  FROM pessoa WHERE cpf = ?`
	row := DB.QueryRow(query, cpf)

	var p entity.Pessoa
	if err := row.Scan(&p.ID, &p.Nome, &p.CPF, &p.RG, &p.Endereco, &p.Contato, &p.ContatoEmergencia); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pessoa por CPF: %w", err)
	}
	return &p, nil
}

// UpdatePessoa atualiza os dados de uma pessoa
func UpdatePessoa(p *entity.Pessoa) error {
	query := `UPDATE pessoa SET nome = ?, cpf = ?, rg = ?, endereco = ?, contato = ?, contatoEmergencia = ?
			  WHERE pessoaID = ?`
	_, err := DB.Exec(query, p.Nome, p.CPF, p.RG, p.Endereco, p.Contato, p.ContatoEmergencia, p.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pessoa: %w", err)
	}
	return nil
}

// DeletePessoa remove uma pessoa do banco de dados
func DeletePessoa(id int64) error {
	query := `DELETE FROM pessoa WHERE pessoaID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar pessoa: %w", err)
	}
	return nil
}

// ExistsPessoaByCPF verifica se já existe uma pessoa com o CPF informado
func ExistsPessoaByCPF(cpf string) (bool, error) {
	query := `SELECT COUNT(*) FROM pessoa WHERE cpf = ?`
	var count int
	err := DB.QueryRow(query, cpf).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência por CPF: %w", err)
	}
	return count > 0, nil
}

// ExistsPessoaByRG verifica se já existe uma pessoa com o RG informado
func ExistsPessoaByRG(rg string) (bool, error) {
	query := `SELECT COUNT(*) FROM pessoa WHERE rg = ?`
	var count int
	err := DB.QueryRow(query, rg).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência por RG: %w", err)
	}
	return count > 0, nil
}

// SearchPessoaByNome busca pessoas com nome semelhante
func SearchPessoaByNome(nome string) ([]*entity.Pessoa, error) {
	query := `SELECT pessoaID, nome, cpf, rg, endereco, contato, contatoEmergencia
			  FROM pessoa WHERE nome LIKE ?`

	nomeLike := fmt.Sprintf("%%%s%%", nome)
	rows, err := DB.Query(query, nomeLike)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pessoas por nome: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em SearchPessoaByNome: %v", cerr)
		}
	}()

	var lista []*entity.Pessoa
	for rows.Next() {
		var p entity.Pessoa
		err := rows.Scan(&p.ID, &p.Nome, &p.CPF, &p.RG, &p.Endereco, &p.Contato, &p.ContatoEmergencia)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler pessoa: %w", err)
		}
		lista = append(lista, &p)
	}
	return lista, nil
}

// ListPessoas retorna todas as pessoas cadastradas
func ListPessoas() ([]*entity.Pessoa, error) {
	query := `SELECT pessoaID, nome, cpf, rg, endereco, contato, contatoEmergencia FROM pessoa`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pessoas: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListPessoas: %v", cerr)
		}
	}()

	var lista []*entity.Pessoa
	for rows.Next() {
		var p entity.Pessoa
		err := rows.Scan(&p.ID, &p.Nome, &p.CPF, &p.RG, &p.Endereco, &p.Contato, &p.ContatoEmergencia)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler pessoa: %w", err)
		}
		lista = append(lista, &p)
	}
	return lista, nil
}
