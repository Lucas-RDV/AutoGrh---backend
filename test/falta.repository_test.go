package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/repository"
	"testing"
	"time"
)

func TestCreateFalta(t *testing.T) {
	funcionario := entity.NewFuncionario("Falta Nome", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	err := repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	falta := &entity.Falta{
		FuncionarioID: funcionario.ID,
		Quantidade:    2,
		Mes:           time.Now(),
	}
	err = repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}
	if falta.ID == 0 {
		t.Error("ID da falta não foi definido")
	}
}

func TestGetFaltasByFuncionarioID(t *testing.T) {
	funcionario := entity.NewFuncionario("Falta Get", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	err := repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}

	falta := &entity.Falta{
		FuncionarioID: funcionario.ID,
		Quantidade:    3,
		Mes:           time.Now(),
	}
	err = repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}

	faltas, err := repository.GetFaltasByFuncionarioID(funcionario.ID)
	if err != nil {
		t.Fatalf("erro ao buscar faltas: %v", err)
	}
	if len(faltas) == 0 {
		t.Error("nenhuma falta retornada")
	}
}

func TestUpdateFalta(t *testing.T) {
	funcionario := entity.NewFuncionario("Falta Update", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	repository.CreateFuncionario(funcionario)

	falta := &entity.Falta{
		FuncionarioID: funcionario.ID,
		Quantidade:    1,
		Mes:           time.Now(),
	}
	repository.CreateFalta(falta)

	falta.Quantidade = 5
	falta.Mes = time.Now().AddDate(0, -1, 0)
	err := repository.UpdateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao atualizar falta: %v", err)
	}
}

func TestDeleteFalta(t *testing.T) {
	funcionario := entity.NewFuncionario("Falta Delete", "RG", "CPF", "PIS", "CTPF", "End", "Contato", "Emergência", "Cargo", time.Now().AddDate(-30, 0, 0), time.Now(), 1000)
	repository.CreateFuncionario(funcionario)

	falta := &entity.Falta{
		FuncionarioID: funcionario.ID,
		Quantidade:    4,
		Mes:           time.Now(),
	}
	repository.CreateFalta(falta)

	err := repository.DeleteFalta(falta.ID)
	if err != nil {
		t.Fatalf("erro ao deletar falta: %v", err)
	}

	faltas, _ := repository.GetFaltasByFuncionarioID(funcionario.ID)
	for _, f := range faltas {
		if f.ID == falta.ID {
			t.Error("falta ainda existe após exclusão")
		}
	}
}

func TestListFaltas(t *testing.T) {
	// Cria funcionário
	funcionario := entity.NewFuncionario(
		"Falta Teste", "123", "456", "789", "000", "Rua Teste", "1111", "2222", "Operador",
		time.Now().AddDate(-20, 0, 0), time.Now(), 1500.0,
	)
	err := repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionário: %v", err)
	}

	// Cria falta
	f := &entity.Falta{
		FuncionarioID: funcionario.ID,
		Quantidade:    1,
		Mes:           time.Now(),
	}
	err = repository.CreateFalta(f)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}

	// Testa a listagem geral
	faltas, err := repository.ListFaltas()
	if err != nil {
		t.Fatalf("erro ao listar faltas: %v", err)
	}

	found := false
	for _, falta := range faltas {
		if falta.ID == f.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("falta de teste não encontrada na listagem geral")
	}
}
