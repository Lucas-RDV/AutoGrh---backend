package testes

import (
	"context"
	"testing"
	
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
)

func newUsuarioService() *service.UsuarioService {
	return service.NewUsuarioService()
}

// --- CREATE NOVO -------------------------------------------------------------

func TestUsuario_Create_NewUser(t *testing.T) {
	defer func() { _ = truncateAll() }()
	svc := newUsuarioService()
	ctx := context.Background()

	reactivated, err := svc.Create(ctx, service.CreateUsuarioInput{
		Username: "alice",
		Password: "senha123",
		IsAdmin:  true,
	})
	if err != nil {
		t.Fatalf("Create (novo) erro: %v", err)
	}
	// contrato: true => reusou row inativa; novo deve ser false
	if reactivated {
		t.Fatalf("Create (novo) deveria retornar false (não reusou row)")
	}

	u, err := repository.GetUsuarioByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("GetUsuarioByUsername erro: %v", err)
	}
	if u == nil {
		t.Fatalf("usuário não encontrado após criar")
	}
	if !u.Ativo {
		t.Errorf("esperava ativo=true")
	}
	if !u.IsAdmin {
		t.Errorf("esperava isAdmin=true")
	}
	if u.Password == "" {
		t.Errorf("hash de senha não persistido")
	}
}

// --- LIST --------------------------------------------------------------------

func TestUsuario_List(t *testing.T) {
	defer func() { _ = truncateAll() }()
	svc := newUsuarioService()
	ctx := context.Background()

	// cria 3 usuários
	for _, name := range []string{"u1", "u2", "u3"} {
		_, err := svc.Create(ctx, service.CreateUsuarioInput{
			Username: name, Password: "x", IsAdmin: false,
		})
		if err != nil {
			t.Fatalf("Create %s erro: %v", name, err)
		}
	}
	// desativa um (soft) pelo service
	u2, _ := repository.GetUsuarioByUsername(ctx, "u2")
	if u2 == nil {
		t.Fatalf("pré-condição inválida: u2 nil")
	}
	if err := svc.Delete(ctx, u2.ID); err != nil {
		t.Fatalf("Delete(u2) erro: %v", err)
	}

	lista, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List erro: %v", err)
	}
	// contrato usual: List retorna apenas ativos
	if len(lista) != 2 {
		t.Fatalf("esperava 2 ativos (u1, u3), veio %d", len(lista))
	}
}

// --- UPDATE ------------------------------------------------------------------

func TestUsuario_Update(t *testing.T) {
	defer func() { _ = truncateAll() }()
	svc := newUsuarioService()
	ctx := context.Background()

	// cria base
	_, err := svc.Create(ctx, service.CreateUsuarioInput{
		Username: "carol", Password: "hash", IsAdmin: false,
	})
	if err != nil {
		t.Fatalf("Create(carol) erro: %v", err)
	}
	u, err := repository.GetUsuarioByUsername(ctx, "carol")
	if err != nil || u == nil {
		t.Fatalf("GetUsuarioByUsername(carol) erro: %v / %v", err, u)
	}

	newUser := "carol2"
	newAdmin := true
	newPass := "outraSenha"
	err = svc.Update(ctx, u.ID, service.UpdateUsuarioInput{
		Username: &newUser,
		Senha:    &newPass,
		IsAdmin:  &newAdmin,
	})
	if err != nil {
		t.Fatalf("Update erro: %v", err)
	}

	u2, err := repository.GetUsuarioByID(u.ID)
	if err != nil {
		t.Fatalf("GetUsuarioByID erro: %v", err)
	}
	if u2 == nil || u2.Username != "carol2" {
		t.Errorf("username não atualizado")
	}
	if !u2.IsAdmin {
		t.Errorf("isAdmin não atualizado")
	}
	if u2.Password == "hash" {
		t.Errorf("senha não foi re-hashada/atualizada")
	}
}

// --- DELETE (soft) + LIST ----------------------------------------------------

func TestUsuario_Delete_Soft_Then_List(t *testing.T) {
	defer func() { _ = truncateAll() }()
	svc := newUsuarioService()
	ctx := context.Background()

	// cria dois ativos
	for _, name := range []string{"d1", "d2"} {
		_, err := svc.Create(ctx, service.CreateUsuarioInput{
			Username: name, Password: "x", IsAdmin: false,
		})
		if err != nil {
			t.Fatalf("Create(%s) erro: %v", name, err)
		}
	}

	lista, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List erro: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("esperava 2 usuários, veio %d", len(lista))
	}

	// inativa (soft) um pelo service
	if err := svc.Delete(ctx, lista[0].ID); err != nil {
		t.Fatalf("Delete (soft) erro: %v", err)
	}

	lista2, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List(2) erro: %v", err)
	}
	if len(lista2) != 1 {
		t.Fatalf("esperava 1 usuário após delete, veio %d", len(lista2))
	}
}

// --- CREATE reaproveitando row por USERNAME ----------------------------------
// Fluxo real do seu backend: se username já existe e está inativo,
// Create reaproveita a mesma row (faz update) e sinaliza true.

func TestUsuario_Create_ReaproveitaRowExistenteInativa(t *testing.T) {
	defer func() { _ = truncateAll() }()
	svc := newUsuarioService()
	ctx := context.Background()

	// 1) cria
	_, err := svc.Create(ctx, service.CreateUsuarioInput{
		Username: "user_reuse", // <= <= 15 chars
		Password: "abc123",
		IsAdmin:  false,
	})
	if err != nil {
		t.Fatalf("Create inicial erro: %v", err)
	}

	u, err := repository.GetUsuarioByUsername(ctx, "user_reuse")
	if err != nil || u == nil {
		t.Fatalf("GetUsuarioByUsername falhou: %v / %v", err, u)
	}

	// 2) desativa (soft)
	if err := svc.Delete(ctx, u.ID); err != nil {
		t.Fatalf("Delete erro: %v", err)
	}

	// 3) Create de novo com mesmo username => deve reaproveitar row e sinalizar true
	reused, err := svc.Create(ctx, service.CreateUsuarioInput{
		Username: "user_reuse",
		Password: "novaSenha",
		IsAdmin:  true, // atualiza isAdmin e senha
	})
	if err != nil {
		t.Fatalf("Create (reaproveitar) erro: %v", err)
	}
	if !reused {
		t.Fatalf("esperava reused=true ao chamar Create com username inativo")
	}

	u2, err := repository.GetUsuarioByUsername(ctx, "user_reuse")
	if err != nil || u2 == nil {
		t.Fatalf("GetUsuarioByUsername pós-reaproveitar falhou: %v / %v", err, u2)
	}
	if !u2.Ativo || !u2.IsAdmin {
		t.Fatalf("usuário deveria estar ativo e com isAdmin=true após reaproveitar row via Create")
	}
}
