package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/repository"
	"testing"
	"time"
)

var feriasFuncionarioID int64
var feriasEntity entity.Ferias

func createFeriasFuncionario(t *testing.T) {
	funcionario := entity.NewFuncionario(
		"Funcionario Férias", "12345678", "99999999900", "123456789", "111111", "Rua A", "1234-5678",
		"9999-9999", "Analista", time.Now().AddDate(-25, 0, 0), time.Now(), 3000.00,
	)
	err := repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionário: %v", err)
	}
	feriasFuncionarioID = funcionario.ID
}

func TestCreateFerias(t *testing.T) {
	createFeriasFuncionario(t)

	inicio := time.Now()
	vencimento := inicio.AddDate(0, 1, 0)

	f := &entity.Ferias{
		FuncionarioID: feriasFuncionarioID,
		Dias:          30,
		Inicio:        inicio,
		Vencimento:    vencimento,
		Vencido:       false,
		Valor:         2500.0,
	}

	err := repository.CreateFerias(f)
	if err != nil {
		t.Fatalf("erro ao criar férias: %v", err)
	}
	if f.ID == 0 {
		t.Error("ID das férias não foi definido")
	}
	feriasEntity = *f
}

func TestGetFeriasByFuncionarioID(t *testing.T) {
	fList, err := repository.GetFeriasByFuncionarioID(feriasFuncionarioID)
	if err != nil {
		t.Fatalf("erro ao buscar férias por funcionário: %v", err)
	}

	found := false
	for _, f := range fList {
		if f.ID == feriasEntity.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("férias de teste não encontradas para o funcionário")
	}
}

func TestUpdateFerias(t *testing.T) {
	feriasEntity.Dias = 20
	feriasEntity.Valor = 2000.0
	feriasEntity.Vencido = true

	err := repository.UpdateFerias(&feriasEntity)
	if err != nil {
		t.Fatalf("erro ao atualizar férias: %v", err)
	}

	list, _ := repository.GetFeriasByFuncionarioID(feriasFuncionarioID)
	updated := false
	for _, f := range list {
		if f.ID == feriasEntity.ID && f.Dias == 20 && f.Vencido {
			updated = true
			break
		}
	}
	if !updated {
		t.Error("férias não foram atualizadas corretamente")
	}
}

func TestDeleteFerias(t *testing.T) {
	err := repository.DeleteFerias(feriasEntity.ID)
	if err != nil {
		t.Fatalf("erro ao deletar férias: %v", err)
	}

	list, _ := repository.GetFeriasByFuncionarioID(feriasFuncionarioID)
	for _, f := range list {
		if f.ID == feriasEntity.ID {
			t.Error("férias ainda existem após exclusão")
		}
	}
}

func TestListFerias(t *testing.T) {
	// Cria férias temporárias para teste
	inicio := time.Now()
	vencimento := inicio.AddDate(0, 1, 0)
	f := &entity.Ferias{
		FuncionarioID: feriasFuncionarioID,
		Dias:          10,
		Inicio:        inicio,
		Vencimento:    vencimento,
		Vencido:       false,
		Valor:         1000.0,
	}
	err := repository.CreateFerias(f)
	if err != nil {
		t.Fatalf("erro ao criar férias para listagem: %v", err)
	}

	lista, err := repository.ListFerias()
	if err != nil {
		t.Fatalf("erro ao listar férias: %v", err)
	}

	found := false
	for _, fer := range lista {
		if fer.ID == f.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("férias criadas para listagem não foram encontradas")
	}
}
