package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "expvar"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	irepo "github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/service"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/facade"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/rest/handler"
	_ "github.com/ElfAstAhe/goph-keeper/migrations/server"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
	"golang.org/x/sync/errgroup"
)

// App - приложение
type App struct {
	// ctx - контекст приложения
	ctx context.Context
	// cancel - функция отмены контекста приложения
	cancel context.CancelFunc
	// db - БД приложения
	db utils.DB
	// dbHelper - набор вспомогательных утилит работы с БД
	dbHelper utils.DBHelper
	// conf - настройки приложения
	conf *config.Config
	// log - главный логер приложения
	log logger.Logger
	// keyCipher - генератор hash суммы
	keyCipher utils.Cipher
	// dataCipher - шифрование данных
	dataCipher utils.Cipher
	// dataCipherHelper - набор вспомогательных утилит шифрования
	dataCipherHelper *utils.CipherHelper
	// keysHelper - набор вспомогательных утилит Pub/Priv keys
	keysHelper *utils.RSAKeysHelper
	// jwtHelper - набор вспомогательных утилит jwt
	jwtHelper *utils.JWTHelper
	// jwtHTTPHelper - набор впомогательных утилит jwt/HTTP
	jwtHTTPHelper *utils.JWTHTTPHelper
	// jwtGRPCHelper - набор вспомогательных утилит jwt/gRPC
	jwtGRPCHelper *utils.JWTGRPCHelper
	// authHelper - набор вспомогательных утилит аутентификации
	authHelper *utils.AuthHelper
	// wg - рабочая группа
	wg sync.WaitGroup
	// userRepo - репозиторий пользователей
	userRepo irepo.UserRepository
	// userDataRepo - репозиторий данных пользователя
	userDataRepo irepo.UserDataRepository
	// userService - сервис работы с пользователем
	userService service.UserService
	// authService - сервис аутентификации
	authService service.AuthService
	// userDataService - сервис работы с данными пользователя, связка БЛ и данных
	userDataService service.UserDataService
	// authFacade - фасад аутентификации, связка между конечной точкой и сервисами
	authFacade facade.AuthFacade
	// userFacade - фасад работы с пользователямт, связка между конечной точкой (HTTP/gRPC) и сервисами (ядро)
	userFacade facade.UserFacade
	// userDataFacade - фасад работы с данными пользователя, связка между конечной точкой (HTTP/gRPC) и сервисами (ядро)
	userDataFacade facade.UserDataFacade
	// router - маршрутизация HTTP запросов
	router *handler.AppChiRouter
	// httpServer - HTTP сервер
	httpServer *http.Server
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

	var eg errgroup.Group
	log.Info("start servers...")
	// http
	eg.Go(func() error {
		if err := app.launchHTTPServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Error starting http server with error [%v]", err)

			return err
		}

		return nil
	})
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
	return eg.Wait()
}

func (app *App) launchHTTPServer() error {
	log := app.log.GetLogger("bootstrap http server launch")
	if app.conf.HTTPConfig.UseHTTPS {
		log.Info("enable https")
		return app.httpServer.ListenAndServeTLS(app.conf.HTTPConfig.CertPath, app.conf.HTTPConfig.PrivateKeyPath)
	}

	log.Info("enable http")

	return app.httpServer.ListenAndServe()
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
	app.log.Info("graceful shutdown http server")
	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		app.log.Errorf("error graceful shutdown http server with error [%v]", err)
	}
	// stop gRPC
	//app.Log.Info("graceful shutdown gRPC server")
	//app.grpcServer.GracefulStop()
}

func (app *App) GetLogger() logger.Logger {
	return app.log
}
