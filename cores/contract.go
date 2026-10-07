package cores

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/shamaton/msgpack/v3"
	"go.uber.org/zap"
)

type HookFunc func(core *AppContracts)
type MiddlewareFunc func(app *fiber.App)
type RouteFunc func(app *fiber.App)
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

func (app *AppContracts) RegisterErrorMapper(mapper ...ErrorMapperFn) *AppContracts {
	app.errorMappers = append(app.errorMappers, mapper...)
	return app
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
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     120 * time.Second,
			MsgPackEncoder:  msgpack.Marshal,
			MsgPackDecoder:  SafeUnmarshal,
		}
		if len(config) > 0 {
			cfg = config[0]
			if cfg.StructValidator == nil {
				cfg.StructValidator = NewStructValidator()
			}
			if cfg.ErrorHandler == nil {
				cfg.ErrorHandler = app.GlobalErrorHandler
			}
			if cfg.MsgPackEncoder == nil {
				cfg.MsgPackEncoder = msgpack.Marshal
			}
			if cfg.MsgPackDecoder == nil {
				cfg.MsgPackDecoder = SafeUnmarshal
			}
		}
		app.App = fiber.New(cfg)
		app.App.RegisterCustomBinder(NewMsgPackBinder())
	})
	return app
}

func (app *AppContracts) RegisterBeforeStart(hook PreStartHook) *AppContracts {
	app.beforeStartHooks = append(app.beforeStartHooks, hook)
	return app
}

func (app *AppContracts) RegisterHook(hook HookFunc) *AppContracts {
	hook(app)
	return app
}

func (app *AppContracts) RegisterMiddleware(mw MiddlewareFunc) *AppContracts {
	mw(app.App)
	return app
}

func (app *AppContracts) RegisterRoute(routes ...RouteFunc) *AppContracts {
	for _, route := range routes {
		route(app.App)
	}
	return app
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
