package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

var pagamentoFuncionarioId int64
var pagamentoEntity Entity.Pagamento

func createPagamentoFuncionario(t *testing.T) {
	funcionario := Entity.NewFuncionario(
		"Teste", "123", "456", "789", "000", "Rua Teste", "9999", "8888", "Cargo",
		time.Now().AddDate(-30, 0, 0), time.Now(), 3000.0,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}
	pagamentoFuncionarioId = funcionario.Id
}

func GetTipoPagamentoID(tipo string) (int64, error) {
	var id int64
	row := Repository.DB.QueryRow("SELECT tipoID FROM tipo_pagamento WHERE tipo = ?", tipo)
	err := row.Scan(&id)
	return id, err
}

func TestCreatePagamento(t *testing.T) {
	createPagamentoFuncionario(t)
	folha := Entity.NewFolhaPagamentos(time.Now())
	err := Repository.CreateFolha(folha)
	if err != nil {
		t.Fatalf("erro ao criar folha: %v", err)
	}
	tipoId, err := GetTipoPagamentoID("salario")
	if err != nil {
		t.Fatalf("erro ao obter tipo_pagamento: %v", err)
	}
	p := Entity.NewPagamento(tipoId, time.Now(), 3000.00)
	p.FuncionarioId = pagamentoFuncionarioId
	p.FolhaId = folha.Id
	err = Repository.CreatePagamento(p)
	if err != nil {
		t.Fatalf("erro ao criar pagamento: %v", err)
	}
	if p.Id == 0 {
		t.Error("ID do pagamento não foi definido")
	}
	pagamentoEntity = *p
}

func TestGetPagamentosByFuncionarioID(t *testing.T) {
	// Cria um funcionário de teste
	funcionario := Entity.NewFuncionario(
		"Funcionario Pagamento", "12345678", "99999999900", "123456789", "111111", "Rua A", "1234-5678",
		"9999-9999", "Analista", time.Now().AddDate(-25, 0, 0), time.Now(), 3000.00,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	// Cria uma folha de pagamento
	folha := Entity.NewFolhaPagamentos(time.Now())
	err = Repository.CreateFolha(folha)
	if err != nil {
		t.Fatalf("erro ao criar folha: %v", err)
	}

	// Garante que o tipo "salario" existe e pega o ID
	tipoId, err := GetTipoPagamentoID("salario")
	if err != nil {
		t.Fatalf("erro ao obter tipo_pagamento: %v", err)
	}

	// Cria um pagamento
	p := Entity.NewPagamento(tipoId, time.Now(), 3000.00)
	p.FuncionarioId = funcionario.Id
	p.FolhaId = folha.Id

	err = Repository.CreatePagamento(p)
	if err != nil {
		t.Fatalf("erro ao criar pagamento: %v", err)
	}

	// Busca os pagamentos do funcionário
	pags, err := Repository.GetPagamentosByFuncionarioID(funcionario.Id)
	if err != nil {
		t.Fatalf("erro ao buscar pagamentos: %v", err)
	}

	// Verifica se o pagamento recém-criado está na lista
	found := false
	for _, pag := range pags {
		if pag.Id == p.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("pagamento não encontrado na listagem")
	}
}

func TestUpdatePagamento(t *testing.T) {
	pagamentoEntity.Valor = 3500.00
	err := Repository.UpdatePagamento(&pagamentoEntity)
	if err != nil {
		t.Fatalf("erro ao atualizar pagamento: %v", err)
	}

	pags, _ := Repository.GetPagamentosByFuncionarioID(pagamentoFuncionarioId)
	updated := false
	for _, p := range pags {
		if p.Id == pagamentoEntity.Id && p.Valor == 3500.00 {
			updated = true
			break
		}
	}
	if !updated {
		t.Error("pagamento não foi atualizado corretamente")
	}
}

func TestDeletePagamento(t *testing.T) {
	err := Repository.DeletePagamento(pagamentoEntity.Id)
	if err != nil {
		t.Fatalf("erro ao deletar pagamento: %v", err)
	}

	pags, _ := Repository.GetPagamentosByFuncionarioID(pagamentoFuncionarioId)
	for _, p := range pags {
		if p.Id == pagamentoEntity.Id {
			t.Error("pagamento ainda existe após exclusão")
		}
	}
}

func TestListPagamentos(t *testing.T) {
	// Cria um funcionário
	funcionario := Entity.NewFuncionario(
		"Funcionario List", "12345678", "88888888800", "222222222", "333333", "Rua B", "4444-5555",
		"5555-6666", "Desenvolvedor", time.Now().AddDate(-28, 0, 0), time.Now(), 3200.00,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	// Cria folha
	folha := Entity.NewFolhaPagamentos(time.Now())
	err = Repository.CreateFolha(folha)
	if err != nil {
		t.Fatalf("erro ao criar folha: %v", err)
	}

	// Garante tipo
	tipoId, err := GetTipoPagamentoID("salario")
	if err != nil {
		t.Fatalf("erro ao obter tipo_pagamento: %v", err)
	}

	// Cria pagamento
	p := Entity.NewPagamento(tipoId, time.Now(), 3200.00)
	p.FuncionarioId = funcionario.Id
	p.FolhaId = folha.Id

	err = Repository.CreatePagamento(p)
	if err != nil {
		t.Fatalf("erro ao criar pagamento: %v", err)
	}

	// Lista todos os pagamentos
	pagamentos, err := Repository.ListPagamentos()
	if err != nil {
		t.Fatalf("erro ao listar pagamentos: %v", err)
	}

	// Confirma se o criado está na lista
	found := false
	for _, pag := range pagamentos {
		if pag.Id == p.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("pagamento de teste não encontrado na listagem geral")
	}
}
