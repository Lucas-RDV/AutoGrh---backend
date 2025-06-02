package test

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"testing"
	"time"
)

var documentoFuncionarioID int64
var documentoCriado entity.Documento

func createDocumentoFuncionario(t *testing.T) {
	funcionario := entity.NewFuncionario(
		"Funcionario Documento", "12345678", "99999999900", "123456789", "111111", "Rua A", "1234-5678",
		"9999-9999", "Analista", fakeDate(25), now(), 3000.00,
	)
	err := repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}
	documentoFuncionarioID = funcionario.Id
}

func TestCreateDocumento(t *testing.T) {
	createDocumentoFuncionario(t)
	doc := &entity.Documento{
		FuncionarioId: documentoFuncionarioID,
		Doc:           []byte("RG escaneado"),
	}
	err := repository.CreateDocumento(doc)
	if err != nil {
		t.Fatalf("erro ao criar documento: %v", err)
	}
	if doc.Id == 0 {
		t.Error("documento criado sem ID")
	}
	documentoCriado = *doc
}

func TestGetDocumentosByFuncionarioID(t *testing.T) {
	docs, err := repository.GetDocumentosByFuncionarioID(documentoFuncionarioID)
	if err != nil {
		t.Fatalf("erro ao buscar documentos: %v", err)
	}

	found := false
	for _, d := range docs {
		if d.Id == documentoCriado.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("documento de teste não encontrado na busca por funcionário")
	}
}

func TestListDocumentos(t *testing.T) {
	docs, err := repository.ListDocumentos()
	if err != nil {
		t.Fatalf("erro ao listar documentos: %v", err)
	}

	found := false
	for _, d := range docs {
		if d.Id == documentoCriado.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("documento de teste não encontrado na listagem geral")
	}
}

func TestDeleteDocumento(t *testing.T) {
	err := repository.DeleteDocumento(documentoCriado.Id)
	if err != nil {
		t.Fatalf("erro ao deletar documento: %v", err)
	}
	docs, err := repository.GetDocumentosByFuncionarioID(documentoFuncionarioID)
	if err != nil {
		t.Fatalf("erro ao buscar documentos após deleção: %v", err)
	}
	for _, d := range docs {
		if d.Id == documentoCriado.Id {
			t.Error("documento ainda existe após exclusão")
		}
	}
}

func fakeDate(yearsAgo int) time.Time {
	return time.Now().AddDate(-yearsAgo, 0, 0)
}

func now() time.Time {
	return time.Now()
}
