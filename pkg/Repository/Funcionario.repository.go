package Repository

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/utils/DateStringToTime"
	"AutoGRH/pkg/utils/NullTimeToPtr"
	"AutoGRH/pkg/utils/PtrToNullTime"
	"database/sql"
	"fmt"
	"log"
)

// CreateFuncionario cria um novo funcionário no banco
func CreateFuncionario(f *Entity.Funcionario) error {
	query := `INSERT INTO funcionario (
		nome, rg, cpf, pis, ctpf, endereco, contato, contatoEmergencia,
		nascimento, admissao, demissao, cargo, salarioInicial, feriasDisponiveis)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		f.Nome, f.RG, f.CPF, f.PIS, f.CTPF, f.Endereco, f.Contato, f.ContatoEmergencia,
		f.Nascimento, f.Admissao, PtrToNullTime.PtrToNullTime(f.Demissao), f.Cargo, f.SalarioInicial, f.FeriasDisponiveis,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir funcionário: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do novo funcionário: %w", err)
	}
	f.ID = id
	return nil
}

// GetFuncionarioByID busca um funcionário pelo ID com todos os relacionamentos
func GetFuncionarioByID(id int64) (*Entity.Funcionario, error) {
	query := `SELECT funcionarioID, nome, rg, cpf, pis, ctpf, endereco, contato,
		contatoEmergencia, nascimento, admissao, demissao, cargo, salarioInicial, feriasDisponiveis
		FROM funcionario WHERE funcionarioID = ?`

	row := DB.QueryRow(query, id)

	var f Entity.Funcionario
	var nascimentoStr, admissaoStr string
	var demissao sql.NullTime

	err := row.Scan(
		&f.ID, &f.Nome, &f.RG, &f.CPF, &f.PIS, &f.CTPF, &f.Endereco, &f.Contato,
		&f.ContatoEmergencia, &nascimentoStr, &admissaoStr, &demissao,
		&f.Cargo, &f.SalarioInicial, &f.FeriasDisponiveis,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar funcionário: %w", err)
	}

	f.Nascimento, err = DateStringToTime.DateStringToTime(nascimentoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter nascimento: %w", err)
	}
	f.Admissao, err = DateStringToTime.DateStringToTime(admissaoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter admissão: %w", err)
	}
	f.Demissao = NullTimeToPtr.NullTimeToPtr(demissao)

	err = carregarRelacionamentos(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func carregarRelacionamentos(f *Entity.Funcionario) error {
	if ferias, err := GetFeriasByFuncionarioID(f.ID); err == nil {
		for _, fr := range ferias {
			f.Ferias = append(f.Ferias, *fr)
		}
	} else {
		return fmt.Errorf("erro ao carregar férias: %w", err)
	}
	if salarios, err := GetSalariosByFuncionarioID(f.ID); err == nil {
		f.Salarios = salarios
	} else {
		return fmt.Errorf("erro ao carregar salários: %w", err)
	}
	if documentos, err := GetDocumentosByFuncionarioID(f.ID); err == nil {
		f.Documentos = documentos
	} else {
		return fmt.Errorf("erro ao carregar documentos: %w", err)
	}
	if faltas, err := GetFaltasByFuncionarioID(f.ID); err == nil {
		f.Faltas = faltas
	} else {
		return fmt.Errorf("erro ao carregar faltas: %w", err)
	}
	if pagamentos, err := GetPagamentosByFuncionarioID(f.ID); err == nil {
		f.Pagamentos = pagamentos
	} else {
		return fmt.Errorf("erro ao carregar pagamentos: %w", err)
	}
	if vales, err := GetValesByFuncionarioID(f.ID); err == nil {
		f.Vales = vales
	} else {
		return fmt.Errorf("erro ao carregar vales: %w", err)
	}
	return nil
}

// UpdateFuncionario atualiza os dados de um funcionário
func UpdateFuncionario(f *Entity.Funcionario) error {
	query := `UPDATE funcionario SET
		nome = ?, rg = ?, cpf = ?, pis = ?, ctpf = ?, endereco = ?, contato = ?,
		contatoEmergencia = ?, nascimento = ?, admissao = ?, demissao = ?,
		cargo = ?, salarioInicial = ?, feriasDisponiveis = ?
		WHERE funcionarioID = ?`

	_, err := DB.Exec(query,
		f.Nome, f.RG, f.CPF, f.PIS, f.CTPF, f.Endereco, f.Contato, f.ContatoEmergencia,
		f.Nascimento, f.Admissao, PtrToNullTime.PtrToNullTime(f.Demissao), f.Cargo, f.SalarioInicial, f.FeriasDisponiveis,
		f.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar funcionário: %w", err)
	}
	return nil
}

// DeleteFuncionario remove um funcionário do banco
func DeleteFuncionario(id int64) error {
	query := `DELETE FROM funcionario WHERE funcionarioID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar funcionário: %w", err)
	}
	return nil
}

// ListFuncionarios retorna uma lista leve de funcionários
func ListFuncionarios() ([]*Entity.Funcionario, error) {
	query := `SELECT funcionarioID, nome FROM funcionario`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar funcionários: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows em ListFuncionarios: %v", cerr)
		}
	}()
	var lista []*Entity.Funcionario
	for rows.Next() {
		var f Entity.Funcionario
		err := rows.Scan(&f.ID, &f.Nome)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler funcionário: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}
