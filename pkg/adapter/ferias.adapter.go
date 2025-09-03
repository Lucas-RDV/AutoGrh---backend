package Adapter

import (
	"AutoGRH/pkg/entity"
	"context"
)

// FeriasRepositoryAdapter adapta as funções do repository para a interface do service
type FeriasRepositoryAdapter struct {
	create                   func(f *entity.Ferias) error
	getFeriasByFuncionarioID func(funcionarioID int64) ([]*entity.Ferias, error)
	getByID                  func(id int64) (*entity.Ferias, error)
	update                   func(f *entity.Ferias) error
	delete                   func(id int64) error
	list                     func() ([]*entity.Ferias, error)
}

func NewFeriasRepositoryAdapter(
	create func(f *entity.Ferias) error,
	getFeriasByFuncionarioID func(funcionarioID int64) ([]*entity.Ferias, error),
	getByID func(id int64) (*entity.Ferias, error),
	update func(f *entity.Ferias) error,
	delete func(id int64) error,
	list func() ([]*entity.Ferias, error),
) *FeriasRepositoryAdapter {
	return &FeriasRepositoryAdapter{
		create:                   create,
		getFeriasByFuncionarioID: getFeriasByFuncionarioID,
		getByID:                  getByID,
		update:                   update,
		delete:                   delete,
		list:                     list,
	}
}

// Implementações da interface service.FeriasRepository

func (a *FeriasRepositoryAdapter) Create(_ context.Context, f *entity.Ferias) error {
	return a.create(f)
}

func (a *FeriasRepositoryAdapter) GetFeriasByFuncionarioID(_ context.Context, funcionarioID int64) ([]*entity.Ferias, error) {
	return a.getFeriasByFuncionarioID(funcionarioID)
}

func (a *FeriasRepositoryAdapter) GetByID(_ context.Context, id int64) (*entity.Ferias, error) {
	return a.getByID(id)
}

func (a *FeriasRepositoryAdapter) Update(_ context.Context, f *entity.Ferias) error {
	return a.update(f)
}

func (a *FeriasRepositoryAdapter) Delete(_ context.Context, id int64) error {
	return a.delete(id)
}

func (a *FeriasRepositoryAdapter) List(_ context.Context) ([]*entity.Ferias, error) {
	return a.list()
}
