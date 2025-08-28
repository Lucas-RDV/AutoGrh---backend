package Adapter

import (
	"AutoGRH/pkg/entity"
)

// FaltaRepositoryAdapter adapta funções do repository para a interface service.FaltaRepository
type FaltaRepositoryAdapter struct {
	create    func(f *entity.Falta) error
	update    func(f *entity.Falta) error
	delete    func(id int64) error
	getByID   func(id int64) (*entity.Falta, error)
	getByFunc func(funcionarioID int64) ([]*entity.Falta, error)
	listAll   func() ([]*entity.Falta, error)
}

func NewFaltaRepositoryAdapter(
	create func(f *entity.Falta) error,
	update func(f *entity.Falta) error,
	delete func(id int64) error,
	getByID func(id int64) (*entity.Falta, error),
	getByFunc func(funcionarioID int64) ([]*entity.Falta, error),
	listAll func() ([]*entity.Falta, error),
) *FaltaRepositoryAdapter {
	return &FaltaRepositoryAdapter{
		create:    create,
		update:    update,
		delete:    delete,
		getByID:   getByID,
		getByFunc: getByFunc,
		listAll:   listAll,
	}
}

// Implementação da interface service.FaltaRepository

func (a *FaltaRepositoryAdapter) Create(f *entity.Falta) error {
	return a.create(f)
}

func (a *FaltaRepositoryAdapter) Update(f *entity.Falta) error {
	return a.update(f)
}

func (a *FaltaRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}

func (a *FaltaRepositoryAdapter) GetFaltaByID(id int64) (*entity.Falta, error) {
	return a.getByID(id)
}

func (a *FaltaRepositoryAdapter) GetFaltasByFuncionarioID(funcionarioID int64) ([]*entity.Falta, error) {
	return a.getByFunc(funcionarioID)
}

func (a *FaltaRepositoryAdapter) ListAll() ([]*entity.Falta, error) {
	return a.listAll()
}
