package test

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"testing"
	"time"
)

var salarioFuncionarioId int64
var salarioEntity *entity.Salario

// Cria um funcionário de teste (evita duplicatas)
func ensureFuncionarioDeTeste(t *testing.T) int64 {
	nome := "Funcionario Salario Teste"
	rg := "1234567"
	cpf := "12345678901"
	pis := "12345678900"
	ctpf := "1234567"
	endereco := "Rua X"
	contato := "(11) 99999-9999"
	contatoEmerg := "(11) 98888-8888"
	cargo := "Analista"
	nascimento := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	admissao := time.Now()
	salarioInicial := 3000.00

	// Tenta buscar por nome se já existe
	existentes, _ := repository.ListFuncionarios()
	for _, f := range existentes {
		if f.Nome == nome {
			return f.Id
		}
	}

	f := entity.NewFuncionario(nome, rg, cpf, pis, ctpf, endereco, contato, contatoEmerg, cargo, nascimento, admissao, salarioInicial)
	err := repository.CreateFuncionario(f)
	if err != nil {
		t.Fatalf("erro ao criar funcionário de teste: %v", err)
	}
	return f.Id
}

func TestCreateSalario(t *testing.T) {
	salarioFuncionarioId = ensureFuncionarioDeTeste(t)
	inicio := time.Now()
	valor := 3200.50

	s := entity.NewSalario(salarioFuncionarioId, inicio, valor)
	err := repository.CreateSalario(s)
	if err != nil {
		t.Fatalf("erro ao criar salário: %v", err)
	}
	if s.Id == 0 {
		t.Error("ID do salário não foi atribuído")
	}
	salarioEntity = s
}

func TestGetSalariosByFuncionarioID(t *testing.T) {
	salarios, err := repository.GetSalariosByFuncionarioID(salarioFuncionarioId)
	if err != nil {
		t.Fatalf("erro ao buscar salários: %v", err)
	}
	if len(salarios) == 0 {
		t.Error("nenhum salário encontrado para o funcionário")
	}
}

func TestUpdateSalario(t *testing.T) {
	if salarioEntity == nil {
		t.Skip("salário não criado")
	}
	salarioEntity.Valor = 3500.75
	err := repository.UpdateSalario(salarioEntity)
	if err != nil {
		t.Fatalf("erro ao atualizar salário: %v", err)
	}
}

func TestDeleteSalario(t *testing.T) {
	if salarioEntity == nil {
		t.Skip("salário não criado")
	}
	err := repository.DeleteSalario(salarioEntity.Id)
	if err != nil {
		t.Fatalf("erro ao deletar salário: %v", err)
	}
}
