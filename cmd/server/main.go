package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/bootstrap"
	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

// @title           GophKeeper API
// @version         1.0
// @description     Сервис безопасного хранения паролей, карт и файлов.
// @termsOfService  Free use

// @contact.name   API Support
// @contact.url    https://github.com/ElfAstAhe/goph-keeper
// @contact.email  elf.ast.ahe@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  BearerAuth
// @in                          cookie
// @name                        Authorization
// @description                 Введите токен в формате: Bearer <JWT_TOKEN>
func main() {
	fmt.Println(config.BuildVersionInfo())

	// app instance
	app := bootstrap.NewApp()

	// log
	log := app.GetLogger().GetLogger("main")

	// app initialization
	log.Info("app init")
	if err := app.Init(); err != nil {
		log.Errorf("app init failed [%v]", err)
		defer app.Close()

		panic(errs.NewAppCommonError("app initialization failed", err))
	}

	// app run
	log.Info("app run")
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Stop()
		log.Errorf("app run error [%v]", err)
	}

	app.WaitForStop()

	// app close
	log.Info("app close")
	if err := app.Close(); err != nil {
		log.Errorf("app close error [%v]", err)

		panic(errs.NewAppCommonError("app close failed", err))
	}

	log.Info("app shutdown")
}

func printOrNA(template, val string) {
	if strings.TrimSpace(val) == "" {
		fmt.Printf(template, "N/A")

		return
	}

	fmt.Printf(template, val)
}
