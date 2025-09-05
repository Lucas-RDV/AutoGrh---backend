package service

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UsuarioService struct{}

func NewUsuarioService() *UsuarioService { return &UsuarioService{} }

// lista todos os usuários
func (s *UsuarioService) List(_ context.Context) ([]UserMinimal, error) {
	us, err := repository.GetAllUsuarios()
	if err != nil {
		return nil, err
	}
	out := make([]UserMinimal, 0, len(us))
	for _, u := range us {
		perfil := "usuario"
		if u.IsAdmin {
			perfil = "admin"
		}
		out = append(out, UserMinimal{
			ID:     u.ID,
			Nome:   u.Username,
			Login:  u.Username,
			Perfil: perfil,
		})
	}
	return out, nil
}

// cria um novo usuário ou reativa um existente
func (s *UsuarioService) Create(ctx context.Context, input CreateUsuarioInput) (bool, error) {
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" {
		return false, errors.New("usuário ou senha não podem estar vazios")
	}

	existing, err := repository.GetUsuarioByUsername(ctx, input.Username)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência do usuário: %w", err)
	}
	if existing != nil {
		if existing.Ativo {
			return false, errors.New("usuário já existe e está ativo")
		}
		// reativar usuário inativo
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			return false, fmt.Errorf("erro ao gerar hash da senha: %w", err)
		}
		existing.Password = string(hashed)
		existing.IsAdmin = input.IsAdmin
		existing.Ativo = true
		return true, repository.UpdateUsuario(existing)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return false, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	user := entity.NewUsuario(input.Username, string(hashed), input.IsAdmin)
	return false, repository.CreateUsuario(user)
}

// atualiza dados do usuário
func (s *UsuarioService) Update(ctx context.Context, id int64, input UpdateUsuarioInput) error {
	user, err := repository.GetUsuarioByID(id)
	if err != nil {
		return err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Senha != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Senha), 12)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %w", err)
		}
		user.Password = string(hashed)
	}

	if input.IsAdmin != nil {
		user.IsAdmin = *input.IsAdmin
	}

	return repository.UpdateUsuario(user)
}

// deleta usuário por ID
func (s *UsuarioService) Delete(ctx context.Context, id int64) error {
	return repository.DeleteUsuario(id)
}

// Reativa um usuário específico
func (s *UsuarioService) Reactivate(ctx context.Context, id int64) error {
	user, err := repository.GetUsuarioByID(id)
	if err != nil {
		return err
	}
	if user.Ativo {
		return errors.New("usuário já está ativo")
	}
	user.Ativo = true
	return repository.UpdateUsuario(user)
}

// Structs auxiliares para entrada de dados

// CreateUsuarioInput representa os dados para criar um novo usuário
type CreateUsuarioInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

// UpdateUsuarioInput representa os dados para atualizar um usuário existente
type UpdateUsuarioInput struct {
	Username *string `json:"username"`
	Senha    *string `json:"senha"`
	IsAdmin  *bool   `json:"isAdmin"`
}
