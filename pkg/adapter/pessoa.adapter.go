package Adapter

import (
	"context"

	"AutoGRH/pkg/entity"
)

type PessoaRepositoryAdapter struct {
	create     func(p *entity.Pessoa) error
	getByID    func(id int64) (*entity.Pessoa, error)
	getByCPF   func(cpf string) (*entity.Pessoa, error)
	update     func(p *entity.Pessoa) error
	delete     func(id int64) error
	existsCPF  func(cpf string) (bool, error)
	existsRG   func(rg string) (bool, error)
	searchNome func(nome string) ([]*entity.Pessoa, error)
	listAll    func() ([]*entity.Pessoa, error)
}

func NewPessoaRepositoryAdapter(
	create func(p *entity.Pessoa) error,
	getByID func(id int64) (*entity.Pessoa, error),
	getByCPF func(cpf string) (*entity.Pessoa, error),
	update func(p *entity.Pessoa) error,
	delete func(id int64) error,
	existsCPF func(cpf string) (bool, error),
	existsRG func(rg string) (bool, error),
	searchNome func(nome string) ([]*entity.Pessoa, error),
	listAll func() ([]*entity.Pessoa, error),
) *PessoaRepositoryAdapter {
	return &PessoaRepositoryAdapter{
		create:     create,
		getByID:    getByID,
		getByCPF:   getByCPF,
		update:     update,
		delete:     delete,
		existsCPF:  existsCPF,
		existsRG:   existsRG,
		searchNome: searchNome,
		listAll:    listAll,
	}
}

// Implementações das interfaces usadas no PessoaService
func (a *PessoaRepositoryAdapter) Create(ctx context.Context, p *entity.Pessoa) error {
	return a.create(p)
}

func (a *PessoaRepositoryAdapter) GetByID(ctx context.Context, id int64) (*entity.Pessoa, error) {
	return a.getByID(id)
}

func (a *PessoaRepositoryAdapter) GetByCPF(ctx context.Context, cpf string) (*entity.Pessoa, error) {
	return a.getByCPF(cpf)
}

func (a *PessoaRepositoryAdapter) Update(ctx context.Context, p *entity.Pessoa) error {
	return a.update(p)
}

func (a *PessoaRepositoryAdapter) Delete(ctx context.Context, id int64) error {
	return a.delete(id)
}

func (a *PessoaRepositoryAdapter) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	return a.existsCPF(cpf)
}

func (a *PessoaRepositoryAdapter) ExistsByRG(ctx context.Context, rg string) (bool, error) {
	return a.existsRG(rg)
}

func (a *PessoaRepositoryAdapter) SearchByNome(ctx context.Context, nome string) ([]*entity.Pessoa, error) {
	return a.searchNome(nome)
}

func (a *PessoaRepositoryAdapter) ListAll(ctx context.Context) ([]*entity.Pessoa, error) {
	return a.listAll()
}
