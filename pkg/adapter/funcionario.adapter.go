package Adapter

import (
	"context"

	"AutoGRH/pkg/entity"
)

type FuncionarioRepositoryAdapter struct {
	create       func(f *entity.Funcionario) error
	getByID      func(id int64) (*entity.Funcionario, error)
	update       func(f *entity.Funcionario) error
	delete       func(id int64) error
	listAtivos   func() ([]*entity.Funcionario, error)
	listInativos func() ([]*entity.Funcionario, error)
	listTodos    func() ([]*entity.Funcionario, error)
}

func NewFuncionarioRepositoryAdapter(
	create func(f *entity.Funcionario) error,
	getByID func(id int64) (*entity.Funcionario, error),
	update func(f *entity.Funcionario) error,
	delete func(id int64) error,
	listAtivos func() ([]*entity.Funcionario, error),
	listInativos func() ([]*entity.Funcionario, error),
	listTodos func() ([]*entity.Funcionario, error),
) *FuncionarioRepositoryAdapter {
	return &FuncionarioRepositoryAdapter{
		create:       create,
		getByID:      getByID,
		update:       update,
		delete:       delete,
		listAtivos:   listAtivos,
		listInativos: listInativos,
		listTodos:    listTodos,
	}
}

// Implementações para FuncionarioService
func (a *FuncionarioRepositoryAdapter) Create(ctx context.Context, f *entity.Funcionario) error {
	return a.create(f)
}

func (a *FuncionarioRepositoryAdapter) GetByID(ctx context.Context, id int64) (*entity.Funcionario, error) {
	return a.getByID(id)
}

func (a *FuncionarioRepositoryAdapter) Update(ctx context.Context, f *entity.Funcionario) error {
	return a.update(f)
}

func (a *FuncionarioRepositoryAdapter) Delete(ctx context.Context, id int64) error {
	return a.delete(id)
}

func (a *FuncionarioRepositoryAdapter) ListAtivos(ctx context.Context) ([]*entity.Funcionario, error) {
	return a.listAtivos()
}

func (a *FuncionarioRepositoryAdapter) ListInativos(ctx context.Context) ([]*entity.Funcionario, error) {
	return a.listInativos()
}

func (a *FuncionarioRepositoryAdapter) ListTodos(ctx context.Context) ([]*entity.Funcionario, error) {
	return a.listTodos()
}
