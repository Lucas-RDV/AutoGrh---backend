package test

import (
	"AutoGRH/pkg/entity"
	repository "AutoGRH/pkg/repository"
	"testing"
)

func TestCreateUsuario(t *testing.T) {
	u := entity.NewUsuario("testuser", "testpass", true)
	err := repository.CreateUsuario(u)
	if err != nil {
		t.Fatalf("erro ao criar usuario: %v", err)
	}
	if u.Id == 0 {
		t.Error("ID do usuário não foi definido")
	}
}

func TestGetUsuarioByID(t *testing.T) {
	u := entity.NewUsuario("test2", "test2", false)
	repository.CreateUsuario(u)

	uFetched, err := repository.GetUsuarioByID(u.Id)
	if err != nil {
		t.Fatalf("erro ao buscar usuario: %v", err)
	}
	if uFetched == nil || uFetched.Username != u.Username {
		t.Error("usuario buscado difere do original")
	}
}

func TestUpdateUsuario(t *testing.T) {
	u := entity.NewUsuario("updateuser", "pass", false)
	repository.CreateUsuario(u)

	u.Password = "newpass"
	u.IsAdmin = true
	err := repository.UpdateUsuario(u)
	if err != nil {
		t.Fatalf("erro ao atualizar usuario: %v", err)
	}

	uCheck, _ := repository.GetUsuarioByID(u.Id)
	if uCheck.Password != "newpass" || !uCheck.IsAdmin {
		t.Error("atualização de usuario falhou")
	}
}

func TestDeleteUsuario(t *testing.T) {
	u := entity.NewUsuario("deleteuser", "deletepass", false)
	repository.CreateUsuario(u)

	err := repository.DeleteUsuario(u.Id)
	if err != nil {
		t.Fatalf("erro ao deletar usuario: %v", err)
	}

	uDel, _ := repository.GetUsuarioByID(u.Id)
	if uDel != nil {
		t.Error("usuario ainda existe após exclusão")
	}
}
