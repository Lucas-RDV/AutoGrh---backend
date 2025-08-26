package router

import (
	"net/http"

	"AutoGRH/pkg/controller"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"

	"github.com/go-chi/chi/v5"
)

func New(auth *service.AuthService, pessoaSvc *service.PessoaService, funcSvc *service.FuncionarioService) http.Handler {
	r := chi.NewRouter()

	authCtl := controller.NewAuthController(auth)
	users := service.NewUsuarioService()
	adminCtl := controller.NewAdminController(users)
	pessoaCtl := controller.NewPessoaController(pessoaSvc)
	funcCtl := controller.NewFuncionarioController(funcSvc)

	// Rota pública
	r.Post("/auth/login", authCtl.Login)

	// Rota autenticada básica
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(auth))
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			if claims, ok := middleware.GetClaims(r.Context()); ok {
				w.Write([]byte("ola " + claims.Nome))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		})
	})

	// Rotas com permissão para gerenciar usuários
	r.Route("/admin/usuarios", func(r chi.Router) {
		r.Use(middleware.RequirePerm(auth, "usuario:list"))
		r.Get("/", adminCtl.UsuariosList)

		r.With(middleware.RequirePerm(auth, "usuario:create")).
			Post("/", adminCtl.CreateUsuario)

		r.With(middleware.RequirePerm(auth, "usuario:update")).
			Put("/{id}", adminCtl.UpdateUsuario)

		r.With(middleware.RequirePerm(auth, "usuario:delete")).
			Delete("/{id}", adminCtl.DeleteUsuario)
	})

	// Rotas para pessoas
	r.Route("/pessoas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "pessoa:list")).
			Get("/", pessoaCtl.ListPessoas)

		r.With(middleware.RequirePerm(auth, "pessoa:read")).
			Get("/{id}", pessoaCtl.GetPessoaByID)

		r.With(middleware.RequirePerm(auth, "pessoa:create")).
			Post("/", pessoaCtl.CreatePessoa)

		r.With(middleware.RequirePerm(auth, "pessoa:update")).
			Put("/{id}", pessoaCtl.UpdatePessoa)

		r.With(middleware.RequirePerm(auth, "pessoa:delete")).
			Delete("/{id}", pessoaCtl.DeletePessoa)
	})

	// Rotas para funcionários
	r.Route("/funcionarios", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "funcionario:list")).
			Get("/ativos", funcCtl.ListFuncionariosAtivos)

		r.With(middleware.RequirePerm(auth, "funcionario:list")).
			Get("/inativos", funcCtl.ListFuncionariosInativos)

		r.With(middleware.RequirePerm(auth, "funcionario:list")).
			Get("/todos", funcCtl.ListTodosFuncionarios)

		r.With(middleware.RequirePerm(auth, "funcionario:read")).
			Get("/{id}", funcCtl.GetFuncionarioByID)

		r.With(middleware.RequirePerm(auth, "funcionario:create")).
			Post("/", funcCtl.CreateFuncionario)

		r.With(middleware.RequirePerm(auth, "funcionario:update")).
			Put("/{id}", funcCtl.UpdateFuncionario)

		r.With(middleware.RequirePerm(auth, "funcionario:delete")).
			Delete("/{id}", funcCtl.DeleteFuncionario)
	})

	return r
}
