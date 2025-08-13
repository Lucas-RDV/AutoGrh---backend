package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/Repository"
	"testing"
)

var testUsuario *Entity.Usuario

func TestCreateUsuario(t *testing.T) {
	// Tenta deletar usuário existente para evitar conflito
	Repository.DB.Exec("DELETE FROM usuario WHERE username = ?", "testuser")

	testUsuario = &Entity.Usuario{
		Username: "testuser",
		Password: "123456",
		IsAdmin:  false,
	}
	err := Repository.CreateUsuario(testUsuario)
	if err != nil {
		t.Fatalf("erro ao criar usuario: %v", err)
	}
	if testUsuario.ID == 0 {
		t.Error("ID do usuário não foi atribuído")
	}
}

func TestGetUsuarioByID(t *testing.T) {
	if testUsuario == nil {
		t.Fatal("usuário de teste não foi criado")
	}
	u, err := Repository.GetUsuarioByID(testUsuario.ID)
	if err != nil {
		t.Fatalf("erro ao buscar usuario: %v", err)
	}
	if u == nil || u.Username != testUsuario.Username {
		t.Error("usuario buscado difere do original")
	}
}

func TestUpdateUsuario(t *testing.T) {
	if testUsuario == nil {
		t.Fatal("usuário de teste não foi criado")
	}
	testUsuario.Password = "nova_senha"
	testUsuario.IsAdmin = true
	err := Repository.UpdateUsuario(testUsuario)
	if err != nil {
		t.Fatalf("erro ao atualizar usuario: %v", err)
	}

	u, _ := Repository.GetUsuarioByID(testUsuario.ID)
	if u.Password != "nova_senha" || !u.IsAdmin {
		t.Error("usuario não foi atualizado corretamente")
	}
}

func TestListUsuarios(t *testing.T) {
	usuarios, err := Repository.GetAllUsuarios()
	if err != nil {
		t.Fatalf("erro ao listar usuarios: %v", err)
	}
	found := false
	for _, u := range usuarios {
		if u.ID == testUsuario.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("usuario de teste não encontrado na listagem")
	}
}
