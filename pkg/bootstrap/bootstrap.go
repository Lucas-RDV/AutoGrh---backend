package Bootstrap

import (
	"AutoGRH/pkg/controller"
	"AutoGRH/pkg/worker"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"AutoGRH/pkg/Adapter"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"AutoGRH/pkg/service"
	jwtm "AutoGRH/pkg/service/jwt"
)

type AppConfig struct {
	JWTSecret string
	Auth      service.AuthConfig
	Perms     service.PermissionMap
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvIntDefault(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvInt64Default(k string, def int64) int64 {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

func parseCutoffHours(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return []int{12, 19}
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ',', ';', '|', ' ':
			return true
		}
		return false
	})

	seen := make(map[int]struct{})
	hours := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if n, err := strconv.Atoi(p); err == nil && n >= 0 && n <= 23 {
			if _, ok := seen[n]; !ok {
				seen[n] = struct{}{}
				hours = append(hours, n)
			}
		}
	}
	if len(hours) == 0 {
		return []int{12, 19}
	}
	sort.Ints(hours)
	return hours
}

func Load() AppConfig {
	issuer := getenvDefault("JWT_ISSUER", "AutoGRH")
	ttlH := getenvIntDefault("JWT_EXPIRES_IN_HOURS", 6)
	skewS := getenvIntDefault("JWT_CLOCK_SKEW_SECONDS", 60)
	successID := getenvInt64Default("EVENT_LOGIN_SUCCESS_ID", 1)
	failID := getenvInt64Default("EVENT_LOGIN_FAIL_ID", 2)
	tz := getenvDefault("AUTH_TIMEZONE", "America/Campo_Grande")
	cutoffs := parseCutoffHours(getenvDefault("AUTH_CUTOFF_HOURS", "12,19"))

	cfg := service.AuthConfig{
		Issuer:          issuer,
		AccessTTL:       time.Duration(ttlH) * time.Hour,
		ClockSkew:       time.Duration(skewS) * time.Second,
		LoginSuccessEvt: successID,
		LoginFailEvt:    failID,
		Timezone:        tz,
		CutoffHours:     cutoffs,
	}

	perms := service.PermissionMap{
		"admin": {
			"*": {}, // acesso total
		},
		"usuario": {
			// usuario
			"self:read":      {},
			"self:update":    {},
			"ferias:request": {},

			// pessoa
			"pessoa:create": {}, // pode cadastrar pessoa
			"pessoa:read":   {}, // pode visualizar pessoa
			"pessoa:list":   {}, // pode listar pessoas
			// pessoa:update/delete só admin efetiva

			// funcionario
			"funcionario:create": {}, // pode cadastrar funcionário
			"funcionario:read":   {}, // pode visualizar funcionário
			"funcionario:list":   {}, // pode listar funcionários
			// funcionario:update/delete só admin efetiva
		},
	}

	return AppConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
		Auth:      cfg,
		Perms:     perms,
	}
}

func ConnectDB() error {
	// adapte para a sua função concreta, se necessário
	if repository.DB == nil {
		repository.ConnectDB()
	}
	return repository.DB.PingContext(context.Background())
}

func InitWorkers(
	feriasSvc *service.FeriasService,
	descansoSvc *service.DescansoService,
	salarioRealSvc *service.SalarioRealService,
	funcionarioSvc *service.FuncionarioService,
	faltaSvc *service.FaltaService,
	folhaSvc *service.FolhaPagamentoService,
	avisoSvc *service.AvisoService,
) {
	feriasWorker := worker.NewFeriasWorker(
		feriasSvc,
		descansoSvc,
		salarioRealSvc,
		funcionarioSvc,
		faltaSvc,
	)

	folhaWorker := worker.NewFolhaWorker(folhaSvc)
	avisosWorker := worker.NewAvisosWorker(avisoSvc)

	folhaWorker.Start()
	feriasWorker.Start()
	avisosWorker.Start()
}

func BuildAuth(app AppConfig) *service.AuthService {
	tm := jwtm.NewHS256Manager([]byte(app.JWTSecret))
	userRepo := Adapter.NewUserRepositoryAdapter(repository.GetUsuarioByUsername, nil)
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) { return 0, repository.CreateLog(l) }
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)
	return service.NewAuthService(userRepo, logRepo, tm, app.Auth, app.Perms)
}

// BuildPessoaService inicializa o PessoaService com suas dependências
func BuildPessoaService(auth *service.AuthService) *service.PessoaService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) { return 0, repository.CreateLog(l) }
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	// repositório de Pessoa (direto do pacote repository)
	pessoaRepo := Adapter.NewPessoaRepositoryAdapter(
		repository.CreatePessoa,
		repository.GetPessoaByID,
		repository.GetPessoaByCPF,
		repository.UpdatePessoa,
		repository.DeletePessoa,
		repository.ExistsPessoaByCPF,
		repository.ExistsPessoaByRG,
		repository.SearchPessoaByNome,
		repository.ListPessoas,
	)

	return service.NewPessoaService(auth, logRepo, pessoaRepo)
}

// BuildFuncionarioService inicializa o FuncionarioService com suas dependências
func BuildFuncionarioService(auth *service.AuthService) *service.FuncionarioService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) { return 0, repository.CreateLog(l) }
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	// repositório de Funcionario (direto do pacote repository)
	funcRepo := Adapter.NewFuncionarioRepositoryAdapter(
		repository.CreateFuncionario,
		repository.GetFuncionarioByID,
		repository.UpdateFuncionario,
		repository.DeleteFuncionario,
		repository.ListFuncionariosAtivos,
		repository.ListFuncionariosInativos,
		repository.ListTodosFuncionarios,
	)

	return service.NewFuncionarioService(auth, logRepo, funcRepo)
}

// BuildDocumentoService constrói o DocumentoService com repositório e log
func BuildDocumentoService(auth *service.AuthService) *service.DocumentoService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	docRepo := Adapter.NewDocumentoRepositoryAdapter(
		repository.CreateDocumento,
		repository.GetDocumentosByFuncionarioID,
		repository.GetByID,
		repository.ListDocumentos,
		repository.DeleteDocumento,
	)

	return service.NewDocumentoService(auth, logRepo, docRepo)
}

func BuildFaltaService(auth *service.AuthService) *service.FaltaService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	faltaRepo := Adapter.NewFaltaRepositoryAdapter(
		repository.CreateFalta,
		repository.UpdateFalta,
		repository.DeleteFalta,
		repository.GetFaltaByID,
		repository.GetFaltasByFuncionarioID,
		repository.ListAllFaltas,
	)

	return service.NewFaltaService(auth, logRepo, faltaRepo)
}

func BuildFeriasService(auth *service.AuthService) *service.FeriasService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	repo := Adapter.NewFeriasRepositoryAdapter(
		repository.CreateFerias,
		repository.GetFeriasByFuncionarioID,
		repository.GetFeriasByID,
		repository.UpdateFerias,
		repository.DeleteFerias,
		repository.ListFerias,
	)

	return service.NewFeriasService(auth, logRepo, repo)
}

func BuildDescansoService(auth *service.AuthService) *service.DescansoService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	repo := Adapter.NewDescansoRepositoryAdapter(
		repository.CreateDescanso,
		repository.GetDescansoByID,
		repository.UpdateDescanso,
		repository.DeleteDescanso,
		repository.GetDescansosByFeriasID,
		repository.GetDescansosByFuncionarioID,
		repository.GetDescansosAprovados,
		repository.GetDescansosPendentes,
	)

	return service.NewDescansoService(auth, logRepo, repo)
}

func BuildSalarioService(auth *service.AuthService) *service.SalarioService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	repo := Adapter.NewSalarioRepositoryAdapter(
		repository.CreateSalario,
		repository.GetSalariosByFuncionarioID,
		repository.UpdateSalario,
		repository.DeleteSalario,
	)

	return service.NewSalarioService(auth, logRepo, repo)
}

func BuildSalarioRealService(auth *service.AuthService) *service.SalarioRealService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	repo := Adapter.NewSalarioRealRepositoryAdapter(
		repository.CreateSalarioReal,
		repository.GetSalariosReaisByFuncionarioID,
		repository.GetSalarioRealAtual,
		repository.UpdateSalarioReal,
		repository.DeleteSalarioReal,
	)

	return service.NewSalarioRealService(auth, logRepo, repo)
}

// BuildValeService constrói o ValeService com o adapter apropriado
func BuildValeService(auth *service.AuthService) *service.ValeService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	valeRepo := Adapter.NewValeRepositoryAdapter(
		repository.CreateVale,
		repository.GetValeByID,
		repository.GetValesByFuncionarioID,
		repository.UpdateVale,
		repository.SoftDeleteVale,
		repository.DeleteVale,
		repository.ListValesPendentes,
		repository.ListValesAprovadosNaoPagos,
	)
	return service.NewValeService(valeRepo, auth, logRepo)
}

// BuildValeController constrói o ValeController
func BuildValeController(auth *service.AuthService) *controller.ValeController {
	valeSvc := BuildValeService(auth)
	return controller.NewValeController(valeSvc)
}

// BuildFolhaPagamentoService constrói o FolhaPagamentoService
func BuildFolhaPagamentoService(auth *service.AuthService) *service.FolhaPagamentoService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	folhaRepo := Adapter.NewFolhaPagamentoRepositoryAdapter(
		repository.CreateFolhaPagamento,
		repository.GetFolhaPagamentoByID,
		repository.GetFolhaByMesAnoTipo,
		repository.UpdateFolhaPagamento,
		repository.DeleteFolhaPagamento,
		repository.ListFolhasPagamentos,
		repository.MarcarFolhaComoPaga,
	)
	return service.NewFolhaPagamentoService(folhaRepo, auth, logRepo)
}

// BuildPagamentoService constrói o PagamentoService com o adapter apropriado
func BuildPagamentoService(auth *service.AuthService) *service.PagamentoService {
	createLog := func(ctx context.Context, l *entity.Log) (int64, error) {
		return 0, repository.CreateLog(l)
	}
	logRepo := Adapter.NewLogRepositoryAdapter(createLog)

	pagamentoRepo := Adapter.NewPagamentoRepositoryAdapter(
		repository.CreatePagamento,
		repository.UpdatePagamento,
		repository.GetPagamentosByFolhaID,
		repository.DeletePagamentosByFolhaID,
		repository.GetPagamentoByID,
		repository.ListPagamentosByFuncionarioID,
	)
	return service.NewPagamentoService(pagamentoRepo, auth, logRepo)
}

// BuildPagamentoController constrói o PagamentoController
func BuildPagamentoController(auth *service.AuthService) *controller.PagamentoController {
	pagamentoSvc := BuildPagamentoService(auth)
	return controller.NewPagamentoController(pagamentoSvc)
}

func BuildAvisoService(auth *service.AuthService) *service.AvisoService {
	return service.NewAvisoService(auth)
}
