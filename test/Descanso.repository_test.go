package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

var descansoID int64
var descansoEntity *Entity.Descanso

func createTestDescanso(t *testing.T) *Entity.Descanso {
	// Criação de férias para o descanso
	f := &Entity.Ferias{
		FuncionarioID: 1,
		Dias:          30,
		Inicio:        time.Now(),
		Vencimento:    time.Now().AddDate(1, 0, 0),
		Vencido:       false,
		Valor:         2000.0,
	}
	err := Repository.CreateFerias(f)
	if err != nil {
		t.Fatalf("erro ao criar ferias: %v", err)
	}

	d := &Entity.Descanso{
		FeriasID: f.Id,
		Inicio:   time.Now(),
		Fim:      time.Now().AddDate(0, 0, 5),
		Valor:    500.0,
		Pago:     false,
		Aprovado: false,
	}

	err = Repository.CreateDescanso(d)
	if err != nil {
		t.Fatalf("erro ao criar descanso: %v", err)
	}
	descansoID = d.Id
	descansoEntity = d
	return d
}

func TestCreateDescanso(t *testing.T) {
	createTestDescanso(t)
	if descansoID == 0 {
		t.Error("descanso ID não foi definido")
	}
}

func TestGetDescansoByID(t *testing.T) {
	createTestDescanso(t)
	d, err := Repository.GetDescansoByID(descansoID)
	if err != nil {
		t.Fatalf("erro ao buscar descanso: %v", err)
	}
	if d == nil || d.Id != descansoID {
		t.Error("descanso buscado não corresponde ao esperado")
	}
}

func TestGetDescansosByFeriasID(t *testing.T) {
	d := createTestDescanso(t)
	descansos, err := Repository.GetDescansosByFeriasID(d.FeriasID)
	if err != nil {
		t.Fatalf("erro ao buscar descansos por ferias: %v", err)
	}
	found := false
	for _, item := range descansos {
		if item.Id == descansoID {
			found = true
			break
		}
	}
	if !found {
		t.Error("descanso de teste não encontrado na listagem por ferias")
	}
}

func TestUpdateDescanso(t *testing.T) {
	d := createTestDescanso(t)
	d.Valor = 999.0
	err := Repository.UpdateDescanso(d)
	if err != nil {
		t.Fatalf("erro ao atualizar descanso: %v", err)
	}

	dAtualizado, _ := Repository.GetDescansoByID(d.Id)
	if dAtualizado.Valor != 999.0 {
		t.Error("valor do descanso não foi atualizado corretamente")
	}
}

func TestDeleteDescanso(t *testing.T) {
	d := createTestDescanso(t)
	err := Repository.DeleteDescanso(d.Id)
	if err != nil {
		t.Fatalf("erro ao deletar descanso: %v", err)
	}

	dExcluido, _ := Repository.GetDescansoByID(d.Id)
	if dExcluido != nil {
		t.Error("descanso ainda existe após exclusão")
	}
}

func TestListDescansos(t *testing.T) {
	createTestDescanso(t)
	descansos, err := Repository.ListDescansos()
	if err != nil {
		t.Fatalf("erro ao listar descansos: %v", err)
	}
	if len(descansos) == 0 {
		t.Error("nenhum descanso retornado na listagem geral")
	}
}
