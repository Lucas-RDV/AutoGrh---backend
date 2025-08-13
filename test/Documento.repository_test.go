package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

var documentoFuncionarioID int64
var documentoCriado Entity.Documento

func createDocumentoFuncionario(t *testing.T) {
	funcionario := Entity.NewFuncionario(
		"Funcionario Documento", "12345678", "99999999900", "123456789", "111111", "Rua A", "1234-5678",
		"9999-9999", "Analista", fakeDate(25), now(), 3000.00,
	)
	err := Repository.CreateFuncionario(funcionario)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}
	documentoFuncionarioID = funcionario.ID
}

func TestCreateDocumento(t *testing.T) {
	createDocumentoFuncionario(t)
	doc := &Entity.Documento{
		FuncionarioID: documentoFuncionarioID,
		Doc:           []byte("RG escaneado"),
	}
	err := Repository.CreateDocumento(doc)
	if err != nil {
		t.Fatalf("erro ao criar documento: %v", err)
	}
	if doc.ID == 0 {
		t.Error("documento criado sem ID")
	}
	documentoCriado = *doc
}

func TestGetDocumentosByFuncionarioID(t *testing.T) {
	docs, err := Repository.GetDocumentosByFuncionarioID(documentoFuncionarioID)
	if err != nil {
		t.Fatalf("erro ao buscar documentos: %v", err)
	}

	found := false
	for _, d := range docs {
		if d.ID == documentoCriado.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("documento de teste não encontrado na busca por funcionário")
	}
}

func TestListDocumentos(t *testing.T) {
	docs, err := Repository.ListDocumentos()
	if err != nil {
		t.Fatalf("erro ao listar documentos: %v", err)
	}

	found := false
	for _, d := range docs {
		if d.ID == documentoCriado.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("documento de teste não encontrado na listagem geral")
	}
}

func TestDeleteDocumento(t *testing.T) {
	err := Repository.DeleteDocumento(documentoCriado.ID)
	if err != nil {
		t.Fatalf("erro ao deletar documento: %v", err)
	}
	docs, err := Repository.GetDocumentosByFuncionarioID(documentoFuncionarioID)
	if err != nil {
		t.Fatalf("erro ao buscar documentos após deleção: %v", err)
	}
	for _, d := range docs {
		if d.ID == documentoCriado.ID {
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
