package cores

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	appMu            sync.Mutex
	shutdownOnce     sync.Once
}

func CreateContract() *AppContracts {
	return &AppContracts{}
}

func (app *AppContracts) RegisterErrorMapper(mapper ...ErrorMapperFn) *AppContracts {
	app.appMu.Lock()
	app.errorMappers = append(app.errorMappers, mapper...)
	app.appMu.Unlock()
	return app
}

func (app *AppContracts) Initialize() *AppContracts {
	NewLogger()
	return app
}

func (app *AppContracts) buildFiberApp(config ...fiber.Config) *fiber.App {
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
		custom := config[0]
		if custom.AppName == "" {
			custom.AppName = cfg.AppName
		}
		if custom.ReadTimeout == 0 {
			custom.ReadTimeout = cfg.ReadTimeout
		}
		if custom.WriteTimeout == 0 {
			custom.WriteTimeout = cfg.WriteTimeout
		}
		if custom.IdleTimeout == 0 {
			custom.IdleTimeout = cfg.IdleTimeout
		}
		if custom.StructValidator == nil {
			custom.StructValidator = cfg.StructValidator
		}
		if custom.ErrorHandler == nil {
			custom.ErrorHandler = cfg.ErrorHandler
		}
		if custom.MsgPackEncoder == nil {
			custom.MsgPackEncoder = cfg.MsgPackEncoder
		}
		if custom.MsgPackDecoder == nil {
			custom.MsgPackDecoder = cfg.MsgPackDecoder
		}
		cfg = custom
	}
	fApp := fiber.New(cfg)
	fApp.RegisterCustomBinder(NewMsgPackBinder())
	return fApp
}

func (app *AppContracts) CreateApp(config ...fiber.Config) *AppContracts {
	app.appMu.Lock()
	defer app.appMu.Unlock()

	if app.App == nil {
		app.App = app.buildFiberApp(config...)
	}
	return app
}

func (app *AppContracts) ensureApp() {
	app.CreateApp()
}

func (app *AppContracts) RegisterBeforeStart(hook PreStartHook) *AppContracts {
	app.appMu.Lock()
	app.beforeStartHooks = append(app.beforeStartHooks, hook)
	app.appMu.Unlock()
	return app
}

func (app *AppContracts) RegisterHook(hook HookFunc) *AppContracts {
	app.ensureApp()
	hook(app)
	return app
}

func (app *AppContracts) RegisterMiddleware(mw MiddlewareFunc) *AppContracts {
	app.ensureApp()
	mw(app.App)
	return app
}

func (app *AppContracts) RegisterRoute(routes ...RouteFunc) *AppContracts {
	app.ensureApp()
	for _, route := range routes {
		route(app.App)
	}
	return app
}

func (app *AppContracts) Start(config ...fiber.ListenConfig) error {
	app.ensureApp()
	app.appMu.Lock()
	hooks := append([]PreStartHook(nil), app.beforeStartHooks...)
	app.appMu.Unlock()

	for _, hook := range hooks {
		if err := hook(); err != nil {
			zap.L().Error("pre-start hook failed, aborting server listen", zap.Error(err))
			return err
		}
	}

	app.SetupShutdownHook()

	listenCfg := fiber.ListenConfig{
		ShutdownTimeout: 10 * time.Second,
	}
	if len(config) > 0 {
		custom := config[0]
		if custom.ShutdownTimeout == 0 {
			custom.ShutdownTimeout = listenCfg.ShutdownTimeout
		}
		listenCfg = custom
	}

	if listenCfg.GracefulContext == nil {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		listenCfg.GracefulContext = ctx
	}

	return app.App.Listen(Config().App.Port, listenCfg)
}

func (app *AppContracts) SetupShutdownHook() {
	app.ensureApp()
	app.shutdownOnce.Do(func() {
		app.App.Hooks().OnPostShutdown(func(err error) error {
			if err != nil {
				zap.L().Error("shutdown error", zap.Error(err))
			} else {
				zap.L().Info("server shut down successfully")
			}
			_ = zap.L().Sync()
			_ = CloseLogger()
			return nil
		})
	})
}
