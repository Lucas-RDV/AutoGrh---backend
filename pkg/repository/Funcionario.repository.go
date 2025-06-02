package repository

import (
	"AutoGRH/pkg/entity"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateFuncionario cria um novo funcionário no banco
func CreateFuncionario(f *entity.Funcionario) error {
	query := `INSERT INTO funcionario (
		nome, rg, cpf, pis, ctpf, endereco, contato, contatoEmergencia,
		nascimento, admissao, demissao, cargo, salarioInicial, feriasDisponiveis)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		f.Nome, f.RG, f.CPF, f.PIS, f.CTPF, f.Endereco, f.Contato, f.ContatoEmergencia,
		f.Nascimento, f.Admissao, f.Demissao, f.Cargo, f.SalarioInicial, f.FeriasDisponiveis,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir funcionário: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		f.Id = id
	}
	return err
}

// GetFuncionarioByID busca um funcionário pelo ID com todos os relacionamentos
func GetFuncionarioByID(id int64) (*entity.Funcionario, error) {
	query := `SELECT funcionarioID, nome, rg, cpf, pis, ctpf, endereco, contato,
		contatoEmergencia, nascimento, admissao, demissao, cargo, salarioInicial, feriasDisponiveis
		FROM funcionario WHERE funcionarioID = ?`

	row := DB.QueryRow(query, id)

	var f entity.Funcionario
	var nascimentoStr, admissaoStr string
	var demissao sql.NullTime

	err := row.Scan(
		&f.Id, &f.Nome, &f.RG, &f.CPF, &f.PIS, &f.CTPF, &f.Endereco, &f.Contato,
		&f.ContatoEmergencia, &nascimentoStr, &admissaoStr, &demissao,
		&f.Cargo, &f.SalarioInicial, &f.FeriasDisponiveis,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar funcionário: %w", err)
	}

	// Converte strings para time.Time
	f.Nascimento, err = time.Parse("2006-01-02", nascimentoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter nascimento: %w", err)
	}
	f.Admissao, err = time.Parse("2006-01-02", admissaoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter admissão: %w", err)
	}

	// Trata demissão (pode ser nula)
	if demissao.Valid {
		f.Demissao = &demissao.Time
	} else {
		f.Demissao = nil
	}

	// Carrega os relacionamentos compostos
	if ferias, err := GetFeriasByFuncionarioID(f.Id); err == nil {
		for _, fr := range ferias {
			f.Ferias = append(f.Ferias, *fr)
		}
	}
	if salarios, err := GetSalariosByFuncionarioID(f.Id); err == nil {
		f.Salarios = salarios
	}
	if documentos, err := GetDocumentosByFuncionarioID(f.Id); err == nil {
		f.Documentos = documentos
	}
	if faltas, err := GetFaltasByFuncionarioID(f.Id); err == nil {
		f.Faltas = faltas
	}
	if pagamentos, err := GetPagamentosByFuncionarioID(f.Id); err == nil {
		f.Pagamentos = pagamentos
	}
	if vales, err := GetValesByFuncionarioID(f.Id); err == nil {
		f.Vales = vales
	}

	return &f, nil
}

// UpdateFuncionario atualiza os dados de um funcionário
func UpdateFuncionario(f *entity.Funcionario) error {
	query := `UPDATE funcionario SET
		nome = ?, rg = ?, cpf = ?, pis = ?, ctpf = ?, endereco = ?, contato = ?,
		contatoEmergencia = ?, nascimento = ?, admissao = ?, demissao = ?,
		cargo = ?, salarioInicial = ?, feriasDisponiveis = ?
		WHERE funcionarioID = ?`

	var demissao sql.NullTime
	if f.Demissao != nil {
		demissao = sql.NullTime{Time: *f.Demissao, Valid: true}
	} else {
		demissao = sql.NullTime{Valid: false}
	}

	_, err := DB.Exec(query,
		f.Nome, f.RG, f.CPF, f.PIS, f.CTPF, f.Endereco, f.Contato, f.ContatoEmergencia,
		f.Nascimento, f.Admissao, demissao, f.Cargo, f.SalarioInicial, f.FeriasDisponiveis,
		f.Id,
	)
	return err
}

// DeleteFuncionario remove um funcionário do banco
func DeleteFuncionario(id int64) error {
	query := `DELETE FROM funcionario WHERE funcionarioID = ?`
	_, err := DB.Exec(query, id)
	return err
}

// ListFuncionarios retorna uma lista leve de funcionários
func ListFuncionarios() ([]*entity.Funcionario, error) {
	query := `SELECT funcionarioID, nome FROM funcionario`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar funcionários: %w", err)
	}
	defer rows.Close()

	var lista []*entity.Funcionario
	for rows.Next() {
		var f entity.Funcionario
		err := rows.Scan(&f.Id, &f.Nome)
		if err != nil {
			log.Println("erro ao ler funcionário:", err)
			continue
		}
		lista = append(lista, &f)
	}
	return lista, nil
}
