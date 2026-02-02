package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElfAstAhe/goph-keeper/internal/app/config"
	"github.com/ElfAstAhe/goph-keeper/internal/app/db"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
)

// App - приложение
type App struct {
	ctx    context.Context
	Cancel context.CancelFunc
	db     db.DB
	conf   *config.Config
	log    logger.Logger
}

// NewApp - конструктор структуры App
//
// app instance
//
//	app := bootstrap.NewApp()
func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		Cancel: cancel,
		log:    logger.NewStartupZapLogger(),
	}
}

// Init - обязательный метод инициализации
//
// app initialization
//
//	if err := app.Init(); err != nil {
//		logger.Errorf("app initialization failed [%v]", err)
//
//		os.Exit(1)
//	}
func (app *App) Init() error {
	// ToDo: implement

	return nil
}

// Run - метод запуска приложения
//
//	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//	    logger.Errorf("app run error [%v]", err)
//	}
func (app *App) Run() error {
	// ToDo: implement

	return nil
}

// Close - метод освобождения ресурсов приложения
//
//	if err := app.Close(); err != nil {
//		logger.Errorf("app close error [%v]", err)
//
//		os.Exit(1)
//	}
func (app *App) Close() error {
	log := app.log.GetLogger("bootstrap close")

	//log.Info("save in mem data")
	//if err := app.saveInMemData(); err != nil {
	//    return err
	//}

	log.Info("close db connection")
	if err := db.CloseDB(app.db); err != nil {
		return err
	}

	//log.Info("close audit event service")
	//if err := app.auditEventService.Close(); err != nil {
	//    return err
	//}

	return nil
}

// gracefulShutdown - внутренний метод "агрессивного" закрытия приложения (ctrl+c) + остальные сигналы OS на закрытие
func (app *App) gracefulShutdown() {
	//defer app.WG.Done()
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// awaiting signal
	select {
	case <-sig:
		{
			app.Cancel()
			break
		}
	case <-app.ctx.Done():
		{
			signal.Stop(sig)
			break
		}
	}

	// stop http
	//app.Log.Info("graceful shutdown http server")
	//if err := app.httpServer.Shutdown(context.Background()); err != nil {
	//    app.Log.Errorf("error graceful shutdown http server with error [%v]", err)
	//}
	// stop gRPC
	//app.Log.Info("graceful shutdown gRPC server")
	//app.grpcServer.GracefulStop()
}

func (app *App) GetLogger() logger.Logger {
	return app.log
}
