package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

func TestCreateFalta(t *testing.T) {
	funcionario := Entity.NewFuncionario("Falta Nome", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	falta := &Entity.Falta{
		FuncionarioId: funcionario.Id,
		Quantidade:    2,
		Mes:           time.Now(),
	}
	err = Repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}
	if falta.Id == 0 {
		t.Error("ID da falta não foi definido")
	}
}

func TestGetFaltasByFuncionarioID(t *testing.T) {
	funcionario := Entity.NewFuncionario("Falta Get", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	falta := &Entity.Falta{
		FuncionarioId: funcionario.Id,
		Quantidade:    3,
		Mes:           time.Now(),
	}
	err = Repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}

	faltas, err := Repository.GetFaltasByFuncionarioID(funcionario.Id)
	if err != nil {
		t.Fatalf("erro ao buscar faltas: %v", err)
	}
	if len(faltas) == 0 {
		t.Error("nenhuma falta retornada")
	}
}

func TestUpdateFalta(t *testing.T) {
	funcionario := Entity.NewFuncionario("Falta Update", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	Repository.CreateFuncionario(funcionario)

	falta := &Entity.Falta{
		FuncionarioId: funcionario.Id,
		Quantidade:    1,
		Mes:           time.Now(),
	}
	Repository.CreateFalta(falta)

	falta.Quantidade = 5
	falta.Mes = time.Now().AddDate(0, -1, 0)
	err := Repository.UpdateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao atualizar falta: %v", err)
	}
}

func TestDeleteFalta(t *testing.T) {
	funcionario := Entity.NewFuncionario("Falta Delete", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	Repository.CreateFuncionario(funcionario)

	falta := &Entity.Falta{
		FuncionarioId: funcionario.Id,
		Quantidade:    4,
		Mes:           time.Now(),
	}
	Repository.CreateFalta(falta)

	err := Repository.DeleteFalta(falta.Id)
	if err != nil {
		t.Fatalf("erro ao deletar falta: %v", err)
	}

	faltas, _ := Repository.GetFaltasByFuncionarioID(funcionario.Id)
	for _, f := range faltas {
		if f.Id == falta.Id {
			t.Error("falta ainda existe após exclusão")
		}
	}
}

func TestListFaltas(t *testing.T) {
	// Cria funcionário
	funcionario := Entity.NewFuncionario(
		"Falta Teste", "123", "456", "789", "000", "Rua Teste", "1111", "2222", "Operador",
		time.Now().AddDate(-20, 0, 0), time.Now(), 1500.0,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionário: %v", err)
	}

	// Cria falta
	f := &Entity.Falta{
		FuncionarioId: funcionario.Id,
		Quantidade:    1,
		Mes:           time.Now(),
	}
	err = Repository.CreateFalta(f)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}

	// Testa a listagem geral
	faltas, err := Repository.ListFaltas()
	if err != nil {
		t.Fatalf("erro ao listar faltas: %v", err)
	}

	found := false
	for _, falta := range faltas {
		if falta.Id == f.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("falta de teste não encontrada na listagem geral")
	}
}
