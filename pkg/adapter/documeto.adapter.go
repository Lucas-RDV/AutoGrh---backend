package Adapter

import (
	"AutoGRH/pkg/entity"
	"context"
)

type DocumentoRepositoryAdapter struct {
	create             func(d *entity.Documento) error
	getByFuncionarioID func(funcionarioID int64) ([]entity.Documento, error)
	list               func() ([]entity.Documento, error)
	delete             func(id int64) error
}

func NewDocumentoRepositoryAdapter(
	create func(d *entity.Documento) error,
	getByFuncionarioID func(funcionarioID int64) ([]entity.Documento, error),
	list func() ([]entity.Documento, error),
	delete func(id int64) error,
) *DocumentoRepositoryAdapter {
	return &DocumentoRepositoryAdapter{create, getByFuncionarioID, list, delete}
}

func (a *DocumentoRepositoryAdapter) Create(ctx context.Context, d *entity.Documento) error {
	return a.create(d)
}

func (a *DocumentoRepositoryAdapter) GetByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error) {
	list, err := a.getByFuncionarioID(funcionarioID)
	if err != nil {
		return nil, err
	}
	result := make([]*entity.Documento, 0, len(list))
	for i := range list {
		result = append(result, &list[i])
	}
	return result, nil
}

func (a *DocumentoRepositoryAdapter) List(ctx context.Context) ([]*entity.Documento, error) {
	list, err := a.list()
	if err != nil {
		return nil, err
	}
	result := make([]*entity.Documento, 0, len(list))
	for i := range list {
		result = append(result, &list[i])
	}
	return result, nil
}

func (a *DocumentoRepositoryAdapter) Delete(ctx context.Context, id int64) error {
	return a.delete(id)
}
