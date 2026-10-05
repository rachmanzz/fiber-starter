package cores

import (
	"sync"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type RouteFunc func(app *AppContracts) error
type PreStartHook func() error

type AppContracts struct {
	App              *fiber.App
	beforeStartHooks []PreStartHook
	errorMappers     []ErrorMapperFn
	once             sync.Once
}

func CreateContract() *AppContracts {
	return &AppContracts{}
}

func (app *AppContracts) RegisterErrorMapper(mapper ...ErrorMapperFn) {
	app.errorMappers = append(app.errorMappers, mapper...)
}

func (app *AppContracts) Initialize() *AppContracts {
	NewLogger()
	return app
}

func (app *AppContracts) CreateApp(config ...fiber.Config) *AppContracts {
	app.once.Do(func() {
		cfg := fiber.Config{
			AppName:         Config().App.Name,
			StructValidator: NewStructValidator(),
			ErrorHandler:    app.GlobalErrorHandler,
		}
		if len(config) > 0 {
			cfg = config[0]
			if cfg.StructValidator == nil {
				cfg.StructValidator = NewStructValidator()
			}
			if cfg.ErrorHandler == nil {
				cfg.ErrorHandler = app.GlobalErrorHandler
			}
		}
		app.App = fiber.New(cfg)
	})
	return app
}

func (app *AppContracts) RegisterBeforeStart(hook PreStartHook) {
	app.beforeStartHooks = append(app.beforeStartHooks, hook)
}

func (app *AppContracts) RegisterRoute(route RouteFunc) {
	route(app)
}

func (app *AppContracts) Start() error {
	for _, hook := range app.beforeStartHooks {
		if err := hook(); err != nil {
			zap.L().Error("pre-start hook failed, aborting server listen", zap.Error(err))
			return err
		}
	}
	return app.App.Listen(Config().App.Port)
}

func (app *AppContracts) SetupShutdownHook() {
	app.App.Hooks().OnPostShutdown(func(err error) error {
		if err != nil {
			zap.L().Error("shutdown error", zap.Error(err))
		} else {
			zap.L().Info("server shut down successfully")
		}
		zap.L().Sync()
		return nil
	})
}
