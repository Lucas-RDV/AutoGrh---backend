package Adapter

import (
	"AutoGRH/pkg/entity"
)

type SalarioRepositoryAdapter struct {
	create     func(s *entity.Salario) error
	listByFunc func(funcionarioID int64) ([]*entity.Salario, error)
	update     func(s *entity.Salario) error
	delete     func(id int64) error
}

func NewSalarioRepositoryAdapter(
	create func(s *entity.Salario) error,
	listByFunc func(funcionarioID int64) ([]*entity.Salario, error),
	update func(s *entity.Salario) error,
	delete func(id int64) error,
) *SalarioRepositoryAdapter {
	return &SalarioRepositoryAdapter{
		create:     create,
		listByFunc: listByFunc,
		update:     update,
		delete:     delete,
	}
}

func (a *SalarioRepositoryAdapter) Create(s *entity.Salario) error {
	return a.create(s)
}
func (a *SalarioRepositoryAdapter) GetSalariosByFuncionarioID(funcionarioID int64) ([]*entity.Salario, error) {
	return a.listByFunc(funcionarioID)
}
func (a *SalarioRepositoryAdapter) Update(s *entity.Salario) error {
	return a.update(s)
}
func (a *SalarioRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}
