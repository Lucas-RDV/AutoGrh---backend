package router

import (
	"net/http"

	"AutoGRH/pkg/controller"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"

	"github.com/go-chi/chi/v5"
)

func New(
	auth *service.AuthService,
	pessoaSvc *service.PessoaService,
	funcionarioSvc *service.FuncionarioService,
	documentoSvc *service.DocumentoService,
	faltaSvc *service.FaltaService,
	feriasSvc *service.FeriasService,
	descansoSvc *service.DescansoService,
	salarioSvc *service.SalarioService,
	salarioRealSvc *service.SalarioRealService,
	valeSvc *service.ValeService,
	folhaSvc *service.FolhaPagamentoService,
	pagamentoSvc *service.PagamentoService,

) http.Handler {
	r := chi.NewRouter()

	authCtl := controller.NewAuthController(auth)
	users := service.NewUsuarioService()
	adminCtl := controller.NewAdminController(users)
	pessoaCtl := controller.NewPessoaController(pessoaSvc)
	funcionarioCtl := controller.NewFuncionarioController(funcionarioSvc)
	documentoCtl := controller.NewDocumentoController(documentoSvc)
	faltaCtl := controller.NewFaltaController(faltaSvc)
	feriasCtl := controller.NewFeriasController(feriasSvc)
	descansoCtl := controller.NewDescansoController(descansoSvc)
	salarioCtl := controller.NewSalarioController(salarioSvc)
	salarioRealCtl := controller.NewSalarioRealController(salarioRealSvc)
	valeCtl := controller.NewValeController(valeSvc)
	folhaCtl := controller.NewFolhaPagamentoController(folhaSvc)
	pagamentoCtl := controller.NewPagamentoController(pagamentoSvc)

	// Rota pública
	r.Post("/auth/login", authCtl.Login)
	r.Post("/auth/logout", authCtl.Logout)

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

	// Rotas de Pessoas
	r.Route("/pessoas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "pessoa:list")).Get("/", pessoaCtl.ListPessoas)
		r.With(middleware.RequirePerm(auth, "pessoa:create")).Post("/", pessoaCtl.CreatePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:update")).Put("/{id}", pessoaCtl.UpdatePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:delete")).Delete("/{id}", pessoaCtl.DeletePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:read")).Get("/{id}", pessoaCtl.GetPessoaByID)
	})

	// Rotas de Funcionários
	r.Route("/funcionarios", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/", funcionarioCtl.ListTodosFuncionarios)
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/ativos", funcionarioCtl.ListFuncionariosAtivos)
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/inativos", funcionarioCtl.ListFuncionariosInativos)

		r.With(middleware.RequirePerm(auth, "funcionario:create")).Post("/", funcionarioCtl.CreateFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:update")).Put("/{id}", funcionarioCtl.UpdateFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:delete")).Delete("/{id}", funcionarioCtl.DeleteFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:read")).Get("/{id}", funcionarioCtl.GetFuncionarioByID)

		// Documentos dentro de funcionário
		r.With(middleware.RequirePerm(auth, "documento:create")).Post("/{id}/documentos", documentoCtl.CreateDocumento)
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/{id}/documentos", documentoCtl.GetDocumentosByFuncionarioID)

		// Faltas dentro de funcionário
		r.With(middleware.RequirePerm(auth, "falta:list")).Get("/{id}/faltas", faltaCtl.GetFaltasByFuncionarioID)
		r.With(middleware.RequirePerm(auth, "falta:create")).Post("/{id}/faltas", faltaCtl.CreateFalta)

		// Ferias dentro de funcionário
		r.With(middleware.RequirePerm(auth, "ferias:list")).Get("/{id}/ferias", feriasCtl.GetFeriasByFuncionarioID)

		// Descansos dentro de funcionário
		r.With(middleware.RequirePerm(auth, "descanso:list")).Get("/funcionarios/{id}/descansos", descansoCtl.ListByFuncionario)

	})

	// Rotas diretas de Documentos
	r.Route("/documentos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/", documentoCtl.ListDocumentos)
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/{id}/download", documentoCtl.DownloadDocumento)
		r.With(middleware.RequirePerm(auth, "documento:delete")).Delete("/{id}", documentoCtl.DeleteDocumento)
	})

	// Rotas diretas de Faltas
	r.Route("/faltas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "falta:list")).Get("/", faltaCtl.ListAllFaltas)
		r.With(middleware.RequirePerm(auth, "falta:read")).Get("/{id}", faltaCtl.GetFaltaByID)
		r.With(middleware.RequirePerm(auth, "falta:update")).Put("/{id}", faltaCtl.UpdateFalta)
		r.With(middleware.RequirePerm(auth, "falta:delete")).Delete("/{id}", faltaCtl.DeleteFalta)
	})

	//rotas diretas de Ferias
	r.Route("/ferias", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "ferias:list")).Get("/", feriasCtl.ListFerias)
		r.With(middleware.RequirePerm(auth, "ferias:read")).Get("/{id}", feriasCtl.GetFeriasByID)
		r.With(middleware.RequirePerm(auth, "ferias:update")).Put("/{id}/vencida", feriasCtl.MarcarComoVencidas)
		r.With(middleware.RequirePerm(auth, "ferias:update")).Put("/{id}/terco-pago", feriasCtl.MarcarTercoComoPago)
		r.With(middleware.RequirePerm(auth, "ferias:read")).Get("/{id}/saldo", feriasCtl.GetSaldoFerias)

		// Descansos dentro de férias
		r.With(middleware.RequirePerm(auth, "descanso:list")).Get("/ferias/{id}/descansos", descansoCtl.ListByFerias)
	})

	// Rotas diretas de Descansos
	r.Route("/descansos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "descanso:create")).Post("/", descansoCtl.Create)
		r.With(middleware.RequirePerm(auth, "descanso:update")).Put("/{id}/aprovar", descansoCtl.Aprovar)
		r.With(middleware.RequirePerm(auth, "descanso:update")).Put("/{id}/pagar", descansoCtl.Pagar)
		r.With(middleware.RequirePerm(auth, "descanso:delete")).Delete("/{id}", descansoCtl.Delete)

		r.With(middleware.RequirePerm(auth, "descanso:list")).Get("/aprovados", descansoCtl.ListAprovados)
		r.With(middleware.RequirePerm(auth, "descanso:list")).Get("/pendentes", descansoCtl.ListPendentes)
	})

	// Salários (por funcionário e individuais)
	r.With(middleware.RequirePerm(auth, "salario:create")).Post("/funcionarios/{id}/salarios", salarioCtl.Create)
	r.With(middleware.RequirePerm(auth, "salario:list")).Get("/funcionarios/{id}/salarios", salarioCtl.ListByFuncionario)
	r.With(middleware.RequirePerm(auth, "salario:update")).Put("/salarios/{id}", salarioCtl.Update)
	r.With(middleware.RequirePerm(auth, "salario:delete")).Delete("/salarios/{id}", salarioCtl.Delete)

	// Salários reais (histórico e atual)
	r.With(middleware.RequirePerm(auth, "salarioReal:create")).Post("/funcionarios/{id}/salarios-reais", salarioRealCtl.Create)
	r.With(middleware.RequirePerm(auth, "salarioReal:list")).Get("/funcionarios/{id}/salarios-reais", salarioRealCtl.ListByFuncionario)
	r.With(middleware.RequirePerm(auth, "salarioReal:list")).Get("/funcionarios/{id}/salario-real-atual", salarioRealCtl.GetAtual)
	r.With(middleware.RequirePerm(auth, "salarioReal:delete")).Delete("/salarios-reais/{id}", salarioRealCtl.Delete)

	// Rotas diretas de Vales
	r.Route("/vales", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "vale:create")).Post("/", valeCtl.CriarVale)
		r.With(middleware.RequirePerm(auth, "vale:update")).Put("/{id}", valeCtl.AtualizarVale)
		r.With(middleware.RequirePerm(auth, "vale:delete")).Delete("/{id}", valeCtl.DeleteVale)
		r.With(middleware.RequirePerm(auth, "vale:list")).Get("/", valeCtl.ListarVales)
		r.With(middleware.RequirePerm(auth, "vale:read")).Get("/{id}", valeCtl.GetVale)
		r.With(middleware.RequirePerm(auth, "vale:update")).Put("/{id}/aprovar", valeCtl.AprovarVale)
		r.With(middleware.RequirePerm(auth, "vale:update")).Put("/{id}/pagar", valeCtl.MarcarValeComoPago)
		r.With(middleware.RequirePerm(auth, "vale:list")).Get("/pendentes", valeCtl.ListarValesPendentes)
		r.With(middleware.RequirePerm(auth, "vale:list")).Get("/aprovados-nao-pagos", valeCtl.ListarValesAprovadosNaoPagos)
	})

	// Folha de Pagamentos
	r.Route("/folhas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "folha:list")).Get("/", folhaCtl.ListarFolhas)
		r.With(middleware.RequirePerm(auth, "folha:read")).Get("/{id}", folhaCtl.BuscarFolha)
		r.With(middleware.RequirePerm(auth, "folha:read")).Get("/{mes}/{ano}/{tipo}", folhaCtl.BuscarFolhaPorMesAnoTipo)
		r.With(middleware.RequirePerm(auth, "folha:create")).Post("/vale", folhaCtl.CriarFolhaVale)
		r.With(middleware.RequirePerm(auth, "folha:update")).Put("/{id}/recalcular", folhaCtl.RecalcularFolha)
		r.With(middleware.RequirePerm(auth, "folha:update")).Put("/{id}/recalcular-vale", folhaCtl.RecalcularFolhaVale)
		r.With(middleware.RequirePerm(auth, "folha:update")).Put("/{id}/fechar", folhaCtl.FecharFolha)
		r.With(middleware.RequirePerm(auth, "folha:delete")).Delete("/{id}", folhaCtl.ExcluirFolha)
	})

	// Pagamentos
	r.Route("/pagamentos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "pagamento:read")).Get("/{id}", pagamentoCtl.GetPagamentoByID)
		r.With(middleware.RequirePerm(auth, "pagamento:update")).Put("/{id}", pagamentoCtl.UpdatePagamento)
		r.With(middleware.RequirePerm(auth, "pagamento:update")).Put("/{id}/pagar", pagamentoCtl.MarcarComoPago)
	})

	// Pagamentos por funcionário
	r.Route("/funcionarios/{id}/pagamentos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "pagamento:list")).Get("/", pagamentoCtl.ListarPagamentosFuncionario)
	})

	return r
}
