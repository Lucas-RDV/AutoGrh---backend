package test

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"testing"
	"time"
)

var folhaCriada *entity.FolhaPagamentos

func TestCreateFolha(t *testing.T) {
	folha := entity.NewFolhaPagamentos(time.Now())
	err := repository.CreateFolha(folha)
	if err != nil {
		t.Fatalf("erro ao criar folha: %v", err)
	}
	if folha.Id == 0 {
		t.Error("ID da folha não foi definido")
	}
	folhaCriada = folha
}

func TestGetFolhaByID(t *testing.T) {
	if folhaCriada == nil {
		t.Fatal("Folha não criada previamente")
	}

	folha, err := repository.GetFolhaByID(folhaCriada.Id)
	if err != nil {
		t.Fatalf("erro ao buscar folha por ID: %v", err)
	}
	if folha == nil || folha.Id != folhaCriada.Id {
		t.Error("Folha retornada difere da criada")
	}
}

func TestListFolhas(t *testing.T) {
	folhas, err := repository.ListFolhas()
	if err != nil {
		t.Fatalf("erro ao listar folhas: %v", err)
	}

	found := false
	for _, f := range folhas {
		if f.Id == folhaCriada.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("folha criada não encontrada na listagem")
	}
}
