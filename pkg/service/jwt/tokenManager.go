package jwtm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"AutoGRH/pkg/service"
)

type HS256Manager struct {
	secret []byte
}

func NewHS256Manager(secret []byte) *HS256Manager {
	return &HS256Manager{secret: secret}
}

type authClaims struct {
	UserID int64  `json:"uid"`
	Nome   string `json:"nome"`
	Perfil string `json:"perfil"`
	jwt.RegisteredClaims
}

func (m *HS256Manager) SignAccess(c service.Claims) (string, error) {
	rc := jwt.RegisteredClaims{
		Issuer:    c.Issuer,
		Subject:   strconv.FormatInt(c.UserID, 10),
		IssuedAt:  jwt.NewNumericDate(c.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(c.ExpiresAt),
	}
	ac := authClaims{UserID: c.UserID, Nome: c.Nome, Perfil: c.Perfil, RegisteredClaims: rc}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	signed, err := t.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("assinar jwt: %w", err)
	}
	return signed, nil
}

func (m *HS256Manager) ParseAccess(token string) (service.Claims, error) {
	var ac authClaims
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	_, err := parser.ParseWithClaims(token, &ac, func(t *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return service.Claims{}, translateErr(err)
	}
	return service.Claims{
		UserID:    ac.UserID,
		Nome:      ac.Nome,
		Perfil:    ac.Perfil,
		Issuer:    ac.Issuer,
		IssuedAt:  ac.IssuedAt.Time,
		ExpiresAt: ac.ExpiresAt.Time,
	}, nil
}

func translateErr(err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token expirado")
	}
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return errors.New("assinatura inválida")
	}
	return fmt.Errorf("token inválido: %v", err)
}
