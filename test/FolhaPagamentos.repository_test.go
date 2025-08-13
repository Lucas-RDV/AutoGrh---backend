package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

var folhaCriada *Entity.FolhaPagamentos

func TestCreateFolha(t *testing.T) {
	folha := Entity.NewFolhaPagamentos(time.Now())
	err := Repository.CreateFolha(folha)
	if err != nil {
		t.Fatalf("erro ao criar folha: %v", err)
	}
	if folha.ID == 0 {
		t.Error("ID da folha não foi definido")
	}
	folhaCriada = folha
}

func TestGetFolhaByID(t *testing.T) {
	if folhaCriada == nil {
		t.Fatal("Folha não criada previamente")
	}

	folha, err := Repository.GetFolhaByID(folhaCriada.ID)
	if err != nil {
		t.Fatalf("erro ao buscar folha por ID: %v", err)
	}
	if folha == nil || folha.ID != folhaCriada.ID {
		t.Error("Folha retornada difere da criada")
	}
}

func TestListFolhas(t *testing.T) {
	folhas, err := Repository.ListFolhas()
	if err != nil {
		t.Fatalf("erro ao listar folhas: %v", err)
	}

	found := false
	for _, f := range folhas {
		if f.ID == folhaCriada.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("folha criada não encontrada na listagem")
	}
}
