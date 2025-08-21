package entity

// Usuario representa um usuário que pode acessar o sistema, com permissões e registros de ações

type Usuario struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Logs     []Log  `json:"logs,omitempty"`
	Ativo    bool   `json:"ativo"`
}

// NewUsuario cria uma nova instância de Usuario com permissão definida
func NewUsuario(username, password string, isAdmin bool) *Usuario {
	return &Usuario{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
		Logs:     []Log{},
	}
}
