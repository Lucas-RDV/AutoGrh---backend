package entity

type Usuario struct {
	Id       int64
	Username string
	Password string
	IsAdmin  bool
	Logs     []Log
}

func NewUsuario(username string, password string, isAdmin bool) *Usuario {
	d := new(Usuario)
	d.Username = username
	d.Password = password
	d.IsAdmin = isAdmin
	return d
}
