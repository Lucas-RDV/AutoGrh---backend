package test

import (
	"AutoGRH/pkg/entity"
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
		FuncionarioId: funcionario.Id,
		Quantidade:    2,
		Mes:           time.Now(),
	}
	err = repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}
	if falta.Id == 0 {
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
		FuncionarioId: funcionario.Id,
		Quantidade:    3,
		Mes:           time.Now(),
	}
	err = repository.CreateFalta(falta)
	if err != nil {
		t.Fatalf("erro ao criar falta: %v", err)
	}

	faltas, err := repository.GetFaltasByFuncionarioID(funcionario.Id)
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
		FuncionarioId: funcionario.Id,
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
		FuncionarioId: funcionario.Id,
		Quantidade:    4,
		Mes:           time.Now(),
	}
	repository.CreateFalta(falta)

	err := repository.DeleteFalta(falta.Id)
	if err != nil {
		t.Fatalf("erro ao deletar falta: %v", err)
	}

	faltas, _ := repository.GetFaltasByFuncionarioID(funcionario.Id)
	for _, f := range faltas {
		if f.Id == falta.Id {
			t.Error("falta ainda existe após exclusão")
		}
	}
}
