package Adapter

import (
	"AutoGRH/pkg/entity"
	"context"
)

type DocumentoRepositoryAdapter struct {
	create   func(ctx context.Context, d *entity.Documento) error
	getByFunc func(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error)
	getByID   func(ctx context.Context, id int64) (*entity.Documento, error)
	list     func(ctx context.Context) ([]*entity.Documento, error)
	delete   func(ctx context.Context, id int64) error
}

func NewDocumentoRepositoryAdapter(
	create func(ctx context.Context, d *entity.Documento) error,
	getByFunc func(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error),
	getByID func(ctx context.Context, id int64) (*entity.Documento, error),
	list func(ctx context.Context) ([]*entity.Documento, error),
	delete func(ctx context.Context, id int64) error,
) *DocumentoRepositoryAdapter {
	return &DocumentoRepositoryAdapter{
		create:    create,
		getByFunc: getByFunc,
		getByID:   getByID,
		list:      list,
		delete:    delete,
	}
}

// Implementação da interface DocumentoRepository

func (a *DocumentoRepositoryAdapter) Create(ctx context.Context, d *entity.Documento) error {
	return a.create(ctx, d)
}

func (a *DocumentoRepositoryAdapter) GetByFuncionarioID(ctx context.Context, funcionarioID int64) ([]*entity.Documento, error) {
	return a.getByFunc(ctx, funcionarioID)
}

func (a *DocumentoRepositoryAdapter) GetByID(ctx context.Context, id int64) (*entity.Documento, error) {
	return a.getByID(ctx, id)
}

func (a *DocumentoRepositoryAdapter) List(ctx context.Context) ([]*entity.Documento, error) {
	return a.list(ctx)
}

func (a *DocumentoRepositoryAdapter) Delete(ctx context.Context, id int64) error {
	return a.delete(ctx, id)
}
