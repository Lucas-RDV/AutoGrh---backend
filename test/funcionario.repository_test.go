package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
	"time"
)

func TestCreateFuncionario(t *testing.T) {
	f := Entity.NewFuncionario(
		"João Teste", "1234567", "12345678900", "123456789", "123456", "Rua A", "1111-1111",
		"9999-9999", "Analista", time.Now().AddDate(-30, 0, 0), time.Now(), 2500.00,
	)

	err := Repository.CreateFuncionario(f)
	if err != nil {
		t.Fatalf("erro ao criar funcionario: %v", err)
	}
	if f.Id == 0 {
		t.Error("ID do funcionario não foi definido")
	}
}

func TestGetFuncionarioByID(t *testing.T) {
	f := Entity.NewFuncionario(
		"Maria Busca", "7654321", "98765432100", "987654321", "654321", "Rua B", "2222-2222",
		"8888-8888", "Gerente", time.Now().AddDate(-35, 0, 0), time.Now(), 4000.00,
	)

	err := Repository.CreateFuncionario(f)
	if err != nil {
		t.Fatalf("erro ao criar funcionario para busca: %v", err)
	}

	fetched, err := Repository.GetFuncionarioByID(f.Id)
	if err != nil {
		t.Fatalf("erro ao buscar funcionario: %v", err)
	}
	if fetched == nil || fetched.Nome != f.Nome {
		t.Error("funcionario buscado difere do original")
	}
}

func TestUpdateFuncionario(t *testing.T) {
	f := Entity.NewFuncionario(
		"Lucas Update", "1111111", "11111111111", "111111111", "111111", "Rua C", "3333-3333",
		"7777-7777", "Auxiliar", time.Now().AddDate(-25, 0, 0), time.Now(), 1800.00,
	)

	Repository.CreateFuncionario(f)

	f.Cargo = "Coordenador"
	f.SalarioInicial = 3200.00
	err := Repository.UpdateFuncionario(f)
	if err != nil {
		t.Fatalf("erro ao atualizar funcionario: %v", err)
	}

	updated, _ := Repository.GetFuncionarioByID(f.Id)
	if updated.Cargo != "Coordenador" || updated.SalarioInicial != 3200.00 {
		t.Error("atualização de funcionario falhou")
	}
}

func TestDeleteFuncionario(t *testing.T) {
	f := Entity.NewFuncionario(
		"Ana Deletar", "2222222", "22222222222", "222222222", "222222", "Rua D", "4444-4444",
		"6666-6666", "Técnico", time.Now().AddDate(-20, 0, 0), time.Now(), 2300.00,
	)

	Repository.CreateFuncionario(f)

	err := Repository.DeleteFuncionario(f.Id)
	if err != nil {
		t.Fatalf("erro ao deletar funcionario: %v", err)
	}

	deleted, _ := Repository.GetFuncionarioByID(f.Id)
	if deleted != nil {
		t.Error("funcionario ainda existe após exclusão")
	}
}

func TestListFuncionarios(t *testing.T) {
	// Cria um funcionário de teste
	f := Entity.NewFuncionario(
		"Funcionario Listagem", "999", "888", "777", "666", "Rua Listagem", "1111-1111", "2222-2222",
		"Auxiliar", time.Now().AddDate(-25, 0, 0), time.Now(), 2000.00,
	)

	err := Repository.CreateFuncionario(f)
	if err != nil {
		t.Fatalf("erro ao criar funcionário para listagem: %v", err)
	}

	// Executa a listagem
	funcionarios, err := Repository.ListFuncionarios()
	if err != nil {
		t.Fatalf("erro ao listar funcionários: %v", err)
	}

	// Verifica se o funcionário recém-criado está na lista
	found := false
	for _, funci := range funcionarios {
		if funci.Id == f.Id {
			found = true
			break
		}
	}

	if !found {
		t.Error("funcionário de teste não encontrado na listagem de funcionários")
	}
}
