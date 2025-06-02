package main

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"fmt"
	"log"
)

func main() {
	// 1. Conectar ao banco de dados
	repository.ConnectDB()

	// 2. Criar um novo usuário
	usuario := entity.NewUsuario("lucas_teste", "123456", false)
	err := repository.CreateUsuario(usuario)
	if err != nil {
		log.Fatal("Erro ao criar usuário:", err)
	}
	fmt.Println("Usuário criado com ID:", usuario.Id)

	// 3. Buscar usuário pelo ID
	buscado, err := repository.GetUsuarioByID(usuario.Id)
	if err != nil {
		log.Fatal("Erro ao buscar usuário:", err)
	}
	if buscado != nil {
		fmt.Println("Usuário encontrado:", *buscado)
	} else {
		fmt.Println("Usuário não encontrado.")
	}

	// 4. Atualizar usuário
	buscado.Username = "lucas_editado"
	buscado.IsAdmin = true
	err = repository.UpdateUsuario(buscado)
	if err != nil {
		log.Fatal("Erro ao atualizar usuário:", err)
	}
	fmt.Println("Usuário atualizado.")

	// 5. Listar todos os usuários
	usuarios, err := repository.GetAllUsuarios()
	if err != nil {
		log.Fatal("Erro ao listar usuários:", err)
	}
	fmt.Println("Lista de usuários:")
	for _, u := range usuarios {
		fmt.Printf("- ID: %d, Username: %s, Admin: %v\n", u.Id, u.Username, u.IsAdmin)
	}

	// 6. Deletar usuário
	err = repository.DeleteUsuario(usuario.Id)
	if err != nil {
		log.Fatal("Erro ao deletar usuário:", err)
	}
	fmt.Println("Usuário deletado com sucesso.")

}
