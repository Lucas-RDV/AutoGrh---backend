package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/dateStringToTime"
	"AutoGRH/pkg/utils/nullTimeToPtr"
	"AutoGRH/pkg/utils/ptrToNullTime"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// CreateFuncionario cria um novo funcionário no banco
func CreateFuncionario(f *entity.Funcionario) error {
	if f.PessoaID == 0 {
		return fmt.Errorf("pessoa associada ao funcionário é inválida ou inexistente")
	}

	query := `INSERT INTO funcionario (
		pessoaID, pis, ctpf, nascimento, admissao, demissao,
		cargo, salarioInicial, feriasDisponiveis, ativo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		f.PessoaID, f.PIS, f.CTPF, f.Nascimento, f.Admissao,
		ptrToNullTime.PtrToNullTime(f.Demissao),
		f.Cargo, f.SalarioInicial, f.FeriasDisponiveis, true,
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
func GetFuncionarioByID(id int64) (*entity.Funcionario, error) {
	query := `SELECT funcionarioID, pessoaID, pis, ctpf, nascimento, admissao, demissao,
		cargo, salarioInicial, feriasDisponiveis FROM funcionario WHERE funcionarioID = ?`

	row := DB.QueryRow(query, id)

	var f entity.Funcionario
	var nascimentoStr, admissaoStr string
	var demissao sql.NullTime

	err := row.Scan(
		&f.ID, &f.PessoaID, &f.PIS, &f.CTPF,
		&nascimentoStr, &admissaoStr, &demissao,
		&f.Cargo, &f.SalarioInicial, &f.FeriasDisponiveis,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar funcionário: %w", err)
	}

	f.Nascimento, err = dateStringToTime.DateStringToTime(nascimentoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter nascimento: %w", err)
	}
	f.Admissao, err = dateStringToTime.DateStringToTime(admissaoStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter admissão: %w", err)
	}
	f.Demissao = nullTimeToPtr.NullTimeToPtr(demissao)

	err = carregarRelacionamentos(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// carregarRelacionamentos popula dados relacionados ao funcionário
func carregarRelacionamentos(f *entity.Funcionario) error {
	if ferias, err := GetFeriasByFuncionarioID(f.ID); err == nil {
		for _, fr := range ferias {
			f.Ferias = append(f.Ferias, *fr)
		}
	} else {
		return fmt.Errorf("erro ao carregar férias: %w", err)
	}

	if salarios, err := GetSalariosByFuncionarioID(f.ID); err == nil {
		f.SalariosRegistrados = salarios
	} else {
		return fmt.Errorf("erro ao carregar salários registrados: %w", err)
	}

	if salarioAtual, err := GetSalarioAtual(f.ID); err == nil {
		f.SalarioRegistradoAtual = salarioAtual
	} else {
		return fmt.Errorf("erro ao carregar salário registrado atual: %w", err)
	}

	if salariosReais, err := GetSalariosReaisByFuncionarioID(f.ID); err == nil {
		f.SalariosReais = salariosReais
	} else {
		return fmt.Errorf("erro ao carregar salários reais: %w", err)
	}

	if salarioRealAtual, err := GetSalarioRealAtual(f.ID); err == nil {
		f.SalarioRealAtual = salarioRealAtual
	} else {
		return fmt.Errorf("erro ao carregar salário real atual: %w", err)
	}

	if documentos, err := GetDocumentosByFuncionarioID(context.Background(), f.ID); err == nil {
		for _, d := range documentos {
			f.Documentos = append(f.Documentos, *d)
		}
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
func UpdateFuncionario(f *entity.Funcionario) error {
	query := `UPDATE funcionario SET
		pis = ?, ctpf = ?, nascimento = ?, admissao = ?, demissao = ?,
		cargo = ?, salarioInicial = ?, feriasDisponiveis = ?
		WHERE funcionarioID = ?`

	_, err := DB.Exec(query,
		f.PIS, f.CTPF, f.Nascimento, f.Admissao,
		ptrToNullTime.PtrToNullTime(f.Demissao),
		f.Cargo, f.SalarioInicial, f.FeriasDisponiveis, f.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar funcionário: %w", err)
	}
	return nil
}

// DeleteFuncionario faz soft delete de um funcionário
func DeleteFuncionario(id int64) error {
	query := `UPDATE funcionario SET ativo = FALSE WHERE funcionarioID = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar funcionário: %w", err)
	}
	return nil
}

// listFuncionariosByAtivo é uma função auxiliar para consultas com base no status ativo
func listFuncionariosByAtivo(ativo bool) ([]*entity.Funcionario, error) {
	query := `SELECT funcionarioID, pessoaID FROM funcionario WHERE ativo = ?`

	rows, err := DB.Query(query, ativo)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar funcionários: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows: %v", cerr)
		}
	}()

	var lista []*entity.Funcionario
	for rows.Next() {
		var f entity.Funcionario
		if err := rows.Scan(&f.ID, &f.PessoaID); err != nil {
			return nil, fmt.Errorf("erro ao ler funcionário: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}

// ListFuncionariosAtivos retorna lista de funcionários ativos
func ListFuncionariosAtivos() ([]*entity.Funcionario, error) {
	return listFuncionariosByAtivo(true)
}

// ListFuncionariosInativos retorna lista de funcionários desligados
func ListFuncionariosInativos() ([]*entity.Funcionario, error) {
	return listFuncionariosByAtivo(false)
}

// ListTodosFuncionarios retorna todos os funcionários sem filtro
func ListTodosFuncionarios() ([]*entity.Funcionario, error) {
	query := `SELECT funcionarioID, pessoaID FROM funcionario`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar todos os funcionários: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("erro ao fechar rows: %v", cerr)
		}
	}()

	var lista []*entity.Funcionario
	for rows.Next() {
		var f entity.Funcionario
		if err := rows.Scan(&f.ID, &f.PessoaID); err != nil {
			return nil, fmt.Errorf("erro ao ler funcionário: %w", err)
		}
		lista = append(lista, &f)
	}
	return lista, nil
}

//
// === NOVAS FUNÇÕES AUXILIARES PARA SALÁRIOS ===
//
