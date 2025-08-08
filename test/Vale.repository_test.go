package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

var valeFuncionarioId int64
var valeEntity *Entity.Vale

func createValeFuncionario(t *testing.T) int64 {
	if valeFuncionarioId != 0 {
		return valeFuncionarioId
	}

	funcionario := Entity.NewFuncionario(
		"Vale Tester", "1234567", "98765432100", "12345678900", "123456789", "Rua A",
		"999999999", "888888888", "Auxiliar", time.Date(1990, 5, 5, 0, 0, 0, 0, time.UTC),
		time.Now(), 1200.00,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario de teste para vale: %v", err)
	}

	valeFuncionarioId = funcionario.Id
	return valeFuncionarioId
}

func TestCreateVale(t *testing.T) {
	funcId := createValeFuncionario(t)
	vale := Entity.NewVale(funcId, 300.0, time.Now())

	err := Repository.CreateVale(vale)
	if err != nil {
		t.Fatalf("erro ao criar vale: %v", err)
	}
	if vale.Id == 0 {
		t.Error("ID do vale nao foi atribuido")
	}
	valeEntity = vale
}

func TestGetValeByID(t *testing.T) {
	if valeEntity == nil {
		t.Skip("vale de teste nao criado")
	}
	v, err := Repository.GetValeByID(valeEntity.Id)
	if err != nil {
		t.Fatalf("erro ao buscar vale: %v", err)
	}
	if v == nil || v.Valor != valeEntity.Valor {
		t.Error("vale retornado incorreto")
	}
}

func TestGetValesByFuncionarioID(t *testing.T) {
	vales, err := Repository.GetValesByFuncionarioID(valeFuncionarioId)
	if err != nil {
		t.Fatalf("erro ao buscar vales do funcionario: %v", err)
	}
	if len(vales) == 0 {
		t.Error("nenhum vale retornado para o funcionario")
	}
}

func TestUpdateVale(t *testing.T) {
	if valeEntity == nil {
		t.Skip("vale de teste nao criado")
	}
	valeEntity.Aprovado = true
	valeEntity.Pago = true
	err := Repository.UpdateVale(valeEntity)
	if err != nil {
		t.Fatalf("erro ao atualizar vale: %v", err)
	}

	updated, _ := Repository.GetValeByID(valeEntity.Id)
	if !updated.Aprovado || !updated.Pago {
		t.Error("vale nao foi atualizado corretamente")
	}
}

func TestDeleteVale(t *testing.T) {
	if valeEntity == nil {
		t.Skip("vale de teste nao criado")
	}
	err := Repository.DeleteVale(valeEntity.Id)
	if err != nil {
		t.Fatalf("erro ao deletar vale: %v", err)
	}

	v, err := Repository.GetValeByID(valeEntity.Id)
	if err != nil {
		t.Fatalf("erro ao buscar vale apos exclusao: %v", err)
	}
	if v != nil {
		t.Error("vale nao foi excluido")
	}
}
