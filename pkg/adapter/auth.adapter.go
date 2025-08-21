package Adapter

import (
	"context"
	"strings"
	"time"

	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
)

type UserRepositoryAdapter struct {
	findByLogin     func(ctx context.Context, login string) (*entity.Usuario, error)
	updateLastLogin func(ctx context.Context, userID int64, when time.Time) error
}

func NewUserRepositoryAdapter(
	find func(ctx context.Context, login string) (*entity.Usuario, error),
	update func(ctx context.Context, userID int64, when time.Time) error,
) *UserRepositoryAdapter {
	return &UserRepositoryAdapter{findByLogin: find, updateLastLogin: update}
}

func (a *UserRepositoryAdapter) GetByLogin(ctx context.Context, login string) (*service.UserRecord, error) {
	u, err := a.findByLogin(ctx, login)
	if err != nil || u == nil {
		return nil, err
	}
	return mapUsuarioToUserRecord(u), nil
}

func (a *UserRepositoryAdapter) UpdateLastLogin(ctx context.Context, userID int64, when time.Time) error {
	if a.updateLastLogin == nil {
		return nil
	}
	return a.updateLastLogin(ctx, userID, when)
}

type LogRepositoryAdapter struct {
	create func(ctx context.Context, l *entity.Log) (int64, error)
}

func NewLogRepositoryAdapter(create func(ctx context.Context, l *entity.Log) (int64, error)) *LogRepositoryAdapter {
	return &LogRepositoryAdapter{create: create}
}

func (a *LogRepositoryAdapter) Create(ctx context.Context, entry service.LogEntry) (int64, error) {
	var uid int64
	if entry.UsuarioID != nil {
		uid = *entry.UsuarioID
	}
	l := &entity.Log{
		EventoID:  entry.EventoID,
		UsuarioID: uid,
		Data:      entry.Quando,
		Message:   entry.Detalhe,
	}
	return a.create(ctx, l)
}

func mapUsuarioToUserRecord(u *entity.Usuario) *service.UserRecord {
	perfil := "usuario"
	if u.IsAdmin {
		perfil = "admin"
	}
	return &service.UserRecord{
		ID:        u.ID,
		Nome:      u.Username,
		Login:     u.Username,
		Perfil:    perfil,
		Ativo:     true,
		SenhaHash: strings.TrimSpace(u.Password),
	}
}
