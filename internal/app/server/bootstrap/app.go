package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "expvar"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	irepo "github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

// App - приложение
type App struct {
	ctx              context.Context
	cancel           context.CancelFunc
	db               utils.DB
	conf             *config.Config
	log              logger.Logger
	keyCipher        utils.Cipher
	dataCipher       utils.Cipher
	dataCipherHelper *utils.CipherHelper
	jwtHelper        *utils.JWTHelper
	authHelper       *utils.AuthHelper
	wg               sync.WaitGroup
	userRepo         irepo.UserRepository
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
		cancel: cancel,
		log:    logger.NewStartupZapLogger(),
	}
}

// Init - обязательный метод инициализации
//
// app initialization
//
//	if err := app.Init(); err != nil {
//		log.Errorf("app initialization failed [%v]", err)
//		defer app.Close()
//
//		panic(errs.NewAppCommonError("app initialization failed", err))
//	}
func (app *App) Init() error {
	log := app.log.GetLogger("bootstrap init")
	//    defer _utl.CloseOnly(logger.(io.Closer))

	app.conf = config.NewConfig()
	log.Info("loading config")
	if err := app.loadConfig(); err != nil {
		return err
	}

	log.Info("init logger")
	if err := app.initLogger(); err != nil {
		return err
	}

	log.Info("init helpers")
	if err := app.initHelpers(); err != nil {
		return err
	}

	log.Info("init database")
	if err := app.initDatabase(); err != nil {
		return err
	}

	log.Info("migrate database")
	if err := app.migrateDatabase(); err != nil {
		return err
	}

	//log.Info("load im mem data")
	//if err := app.loadInMemData(); err != nil {
	//    return err
	//}

	log.Info("init dependencies")
	if err := app.initDependencies(); err != nil {
		return err
	}

	log.Info("init startup services")
	if err := app.initStartupServices(); err != nil {
		return err
	}

	log.Info("init http router")
	if err := app.initHTTPRouter(); err != nil {
		return err
	}

	log.Info("init http server")
	if err := app.initHTTPServer(); err != nil {
		return err
	}

	log.Info("init gRPC server")
	if err := app.initGRPCServer(); err != nil {
		return err
	}

	return nil
}

// Run - метод запуска приложения
//
//	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//		app.Stop()
//		log.Errorf("app run error [%v]", err)
//	}
func (app *App) Run() error {
	log := app.log.GetLogger("bootstrap run")

	log.Info("start graceful shutdown")
	app.wg.Add(1)
	go app.gracefulShutdown()

	//var eg errgroup.Group
	//log.Info("start servers...")
	//// http
	//eg.Go(func() error {
	//    if err := app.launchHTTPServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//        log.Errorf("Error starting http server with error [%v]", err)
	//
	//        return err
	//    }
	//
	//    return nil
	//})
	//// gRPC
	//eg.Go(func() error {
	//    if err := app.launchGRPCServer(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
	//        log.Errorf("Error starting gRPC server with error [%v]", err)
	//
	//        return err
	//    }
	//
	//    return nil
	//})
	//
	//return eg.Wait()

	return nil
}

// Stop - метод остановки приложения
func (app *App) Stop() {
	app.cancel()
}

func (app *App) WaitForStop() {
	app.wg.Wait()
}

// Close - метод освобождения ресурсов приложения
//
//	if err := app.Close(); err != nil {
//		log.Errorf("app close error [%v]", err)
//
//		panic(errs.NewAppCommonError("app close failed", err))
//	}
func (app *App) Close() error {
	log := app.log.GetLogger("bootstrap close")

	//log.Info("save in mem data")
	//if err := app.saveInMemData(); err != nil {
	//    return err
	//}

	log.Info("close db connection")
	if err := utils.DBClose(app.db); err != nil {
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
	defer app.wg.Done()
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// awaiting signal
	select {
	case <-sig:
		{
			app.cancel()
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
